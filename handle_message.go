package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/hoisie/mustache"
	"github.com/lucsky/cuid"
)

func handleMessage(message *tgbotapi.Message) {
	u, t, err := ensureUser(message.From.ID, message.From.UserName)
	if err != nil {
		log.Warn().Err(err).Int("case", t).
			Str("username", message.From.UserName).
			Int("id", message.From.ID).
			Msg("failed to ensure user")
		return
	}

	if message.Chat.Type == "private" {
		// after ensuring the user we should always enable him to
		// receive payment notifications and so on, as not all people will
		// remember to call /start
		u.setChat(message.Chat.ID)
	} else if message.Entities == nil || len(*message.Entities) == 0 ||
		// unless in the private chat, only messages starting with
		// bot commands will work
		(*message.Entities)[0].Type != "bot_command" ||
		(*message.Entities)[0].Offset != 0 {
		return
	}

	var (
		opts    = make(docopt.Opts)
		proceed = false
		text    = regexp.MustCompile("/([a-z]+)@"+s.ServiceId).ReplaceAllString(message.Text, "/$1")
	)

	log.Debug().Str("t", text).Str("user", u.Username).Msg("got message")

	// when receiving a forwarded invoice (from messages from other people?)
	// or just the full text of a an invoice (shared from a phone wallet?)
	if !strings.HasPrefix(text, "/") {
		if bolt11, ok := searchForInvoice(*message); ok {
			opts, _, _ = parse("/pay " + bolt11)
			goto parsed
		}
	}

	// individual transaction query
	if strings.HasPrefix(text, "/tx") {
		hashfirstchars := text[3:]
		txn, err := u.getTransaction(hashfirstchars)
		if err != nil {
			log.Warn().Err(err).Str("user", u.Username).Str("hash", hashfirstchars).
				Msg("failed to get transaction")
			u.notifyAsReply("Couldn't find transaction "+hashfirstchars+".", message.MessageID)
			return
		}

		txnreply := mustache.Render(`
<code>{{Status}}</code> {{#TelegramPeer.Valid}}{{PeerActionDescription}}{{/TelegramPeer.Valid}} on {{TimeFormat}} {{#IsUnclaimed}}(💤 unclaimed){{/IsUnclaimed}}
<i>{{Description}}</i>{{^TelegramPeer.Valid}} 
{{#Payee.Valid}}<b>Payee</b>: {{{PayeeLink}}} ({{PayeeAlias}}){{/Payee.Valid}}
<b>Hash</b>: {{Hash}}{{/TelegramPeer.Valid}}{{#HasPreimage}} 
<b>Preimage</b>: {{Preimage}}{{/HasPreimage}}
<b>Amount</b>: {{Satoshis}} sat
{{^IsReceive}}<b>Fee paid</b>: {{FeeSatoshis}}{{/IsReceive}}
        `, txn) + "\n" + renderLogInfo(hashfirstchars)
		id := u.notifyAsReply(txnreply, txn.TriggerMessage).MessageID

		if txn.Status == "PENDING" {
			// allow people to cancel pending if they're old enough
			editWithKeyboard(u.ChatId, id, text+"\n\nRecheck pending payment?",
				tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Yes", "check="+hashfirstchars),
					),
				),
			)
		}

		if txn.IsUnclaimed() {
			editWithKeyboard(u.ChatId, id, text+"\n\nRetract unclaimed tip?",
				tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Yes", "remunc="+hashfirstchars),
					),
				),
			)
		}

		return
	}

	// query failed transactions (only available in the first 24h after the failure)
	if strings.HasPrefix(text, "/log") {
		hashfirstchars := text[4:]
		u.notify(renderLogInfo(hashfirstchars))
		return
	}

	// otherwise parse the slash command
	opts, proceed, err = parse(text)
	if !proceed {
		return
	}
	if err != nil {
		if message.Chat.Type == "private" {
			// only tell we don't understand commands when in a private chat
			// because these commands we're not understanding
			// may be targeting other bots in a group, so we're spamming people.
			log.Debug().Err(err).Str("command", text).
				Msg("Failed to parse command")

			method := strings.Split(text, " ")[0][1:]
			handled := handleHelp(u, method)
			if !handled {
				u.notify("Could not understand the command. /help")
			}
		}
		return
	}

parsed:
	if opts["paynow"].(bool) {
		opts["pay"] = true
		opts["now"] = true
	}

	switch {
	case opts["start"].(bool):
		if message.Chat.Type == "private" {
			u.setChat(message.Chat.ID)
			u.notify("Your account is created.")
			handleHelp(u, "")
		}
		break
	case opts["stop"].(bool):
		if message.Chat.Type == "private" {
			u.unsetChat()
			u.notify("Notifications stopped.")
		}
		break
	case opts["receive"].(bool), opts["invoice"].(bool), opts["fund"].(bool):
		sats, err := opts.Int("<satoshis>")
		if err != nil {
			// couldn't get an integer, but maybe it's because nothing was specified, so
			// it's an invoice of undefined amount.

			if v, exists := opts["<satoshis>"]; exists && v != nil && v.(string) != "any" {
				// ok, it exists, so it's an invalid amount.
				u.notify("Invalid amount: " + v.(string))
				break
			}

			// will be this if "any"
			sats = INVOICE_UNDEFINED_AMOUNT
		}
		var desc string
		if idesc, ok := opts["<description>"]; ok {
			desc = strings.Join(idesc.([]string), " ")
		}

		var preimage string
		if param, ok := opts["--preimage"]; ok {
			preimage, _ = param.(string)
		}

		bolt11, qrpath, err := makeInvoice(u, sats, desc, message.MessageID, preimage)
		if err != nil {
			log.Warn().Err(err).Msg("failed to generate invoice")
			notify(message.Chat.ID, messageFromError(err, "Failed to generate invoice"))
			return
		}

		if qrpath == "" {
			u.notify(bolt11)
		} else {
			defer os.Remove(qrpath)
			photo := tgbotapi.NewPhotoUpload(message.Chat.ID, qrpath)
			photo.Caption = bolt11
			_, err := bot.Send(photo)
			if err != nil {
				log.Warn().Str("user", u.Username).Err(err).
					Msg("error sending photo")

					// send just the bolt11
				notify(message.Chat.ID, bolt11)
			}
		}

		break
	case opts["send"].(bool), opts["tip"].(bool):
		// default notify function to use depending on many things
		defaultNotify := func(m string) { u.notify(m) }
		if message.Chat.Type == "private" {
			defaultNotify = func(m string) { u.notifyAsReply(m, message.MessageID) }
		} else if isSpammy(message.Chat.ID) {
			defaultNotify = func(m string) { notifyAsReply(message.Chat.ID, m, message.MessageID) }
		}

		// sending money to others
		var (
			sats          int
			todisplayname string
			receiver      *User
			usernameval   interface{}
		)

		// get quantity
		sats, err := opts.Int("<satoshis>")

		if err != nil || sats <= 0 {
			// maybe the order of arguments is inverted
			if val, ok := opts["<satoshis>"].(string); ok && val[0] == '@' {
				// it seems to be
				usernameval = val
				if asats, ok := opts["<receiver>"].([]string); ok && len(asats) == 1 {
					sats, _ = strconv.Atoi(asats[0])
					goto gotusername
				}
			}

			defaultNotify("Invalid amount: " + opts["<satoshis>"].(string))
			break
		} else {
			usernameval = opts["<receiver>"]
		}

	gotusername:
		anonymous, _ := opts.Bool("--anonymous")

		receiver, todisplayname, err = parseUsername(message, usernameval)
		if err != nil {
			log.Warn().Interface("val", usernameval).Err(err).Msg("failed to parse username")
			break
		}
		if receiver != nil {
			goto ensured
		}

		// no username, this may be a reply-tip
		if message.ReplyToMessage != nil {
			log.Debug().Msg("it's a reply-tip")
			reply := message.ReplyToMessage

			var t int
			rec, t, err := ensureUser(reply.From.ID, reply.From.UserName)
			receiver = &rec
			if err != nil {
				log.Warn().Err(err).Int("case", t).
					Str("username", reply.From.UserName).
					Int("id", reply.From.ID).
					Msg("failed to ensure user on reply-tip")
				break
			}
			if reply.From.UserName != "" {
				todisplayname = "@" + reply.From.UserName
			} else {
				todisplayname = strings.TrimSpace(
					reply.From.FirstName + " " + reply.From.LastName,
				)
			}
			goto ensured
		}

		// if we ever reach this point then it's because the receiver is missing.
		defaultNotify("Can't send " + opts["<satoshis>"].(string) + ". Missing receiver!")
		break

	ensured:
		if err != nil {
			log.Warn().Err(err).
				Msg("failed to ensure target user on send/tip.")
			defaultNotify("Failed to save receiver. This is probably a bug.")
			break
		}

		errMsg, err := u.sendInternally(
			message.MessageID,
			*receiver,
			anonymous,
			sats*1000,
			nil,
			nil,
		)
		if err != nil {
			log.Warn().Err(err).
				Str("from", u.Username).
				Str("to", todisplayname).
				Msg("failed to send/tip")
			defaultNotify("Failed to send: " + errMsg)
			break
		}

		if receiver.ChatId != 0 {
			if anonymous {
				receiver.notify(fmt.Sprintf("Someone has sent you %d sat.", sats))
			} else {
				receiver.notify(fmt.Sprintf("%s has sent you %d sat.", u.AtName(), sats))
			}
		}

		if message.Chat.Type == "private" {
			warning := ""
			if receiver.ChatId == 0 {
				warning = fmt.Sprintf(
					" (couldn't notify %s as they haven't started a conversation with the bot)",
					todisplayname,
				)
			}
			u.notifyAsReply(
				fmt.Sprintf("%d sat sent to %s%s.", sats, todisplayname, warning),
				message.MessageID,
			)
			break
		}

		defaultNotify(fmt.Sprintf("%d sat sent to %s.", sats, todisplayname))
		break
	case opts["giveaway"].(bool):
		sats, err := opts.Int("<satoshis>")
		if err != nil {
			u.notify("Invalid amount: " + opts["<satoshis>"].(string))
			break
		}
		if !u.checkBalanceFor(sats, "giveaway") {
			break
		}

		chattable := tgbotapi.NewMessage(
			message.Chat.ID,
			fmt.Sprintf("%s is giving %d sat away!", u.AtName(), sats),
		)
		chattable.BaseChat.ReplyMarkup = giveAwayKeyboard(u, sats)
		bot.Send(chattable)
		break
	case opts["coinflip"].(bool), opts["lottery"].(bool):
		// open a lottery between a number of users in a group
		sats, err := opts.Int("<satoshis>")
		if err != nil {
			u.notify("Invalid amount: " + opts["<satoshis>"].(string))
			break
		}
		if !u.checkBalanceFor(sats, "coinflip") {
			break
		}

		nparticipants := 2
		if n, err := opts.Int("<num_participants>"); err == nil {
			if n < 2 || n > 100 {
				u.notify("Invalid number of participants: " + strconv.Itoa(n))
				break
			} else {
				nparticipants = n
			}
		}

		chattable := tgbotapi.NewMessage(
			message.Chat.ID,
			fmt.Sprintf(`A lottery round is starting!

Entry fee: %d sat
Total participants: %d
Prize: %d
Registered: %s`, sats, nparticipants, sats*nparticipants, u.AtName()),
		)

		coinflipid := cuid.Slug()
		rds.SAdd("coinflip:"+coinflipid, u.Id)
		rds.Expire("coinflip:"+coinflipid, s.GiveAwayTimeout)
		chattable.BaseChat.ReplyMarkup = coinflipKeyboard(coinflipid, nparticipants, sats)
		bot.Send(chattable)
	case opts["fundraise"].(bool), opts["crowdfund"].(bool):
		// many people join, we get all the money and transfer to the target
		sats, err := opts.Int("<satoshis>")
		if err != nil {
			u.notify("Invalid amount: " + opts["<satoshis>"].(string))
			break
		}
		if !u.checkBalanceFor(sats, "fundraise") {
			break
		}

		nparticipants, err := opts.Int("<num_participants>")
		if err != nil || nparticipants < 2 || nparticipants > 100 {
			u.notify("Invalid number of participants: " + strconv.Itoa(nparticipants))
			break
		}

		receiver, receiverdisplayname, err := parseUsername(message, opts["<receiver>"])
		if err != nil {
			log.Warn().Err(err).Msg("parsing fundraise receiver")
			u.notify("Failed to parse receiver name.")
			break
		}

		chattable := tgbotapi.NewMessage(
			message.Chat.ID,
			fmt.Sprintf(`A fundraising to %s was started!

Contributors needed for completion: %d
Each pays: %d sat
Final amount: %d
Have contributed: %s`, receiverdisplayname, nparticipants, sats, sats*nparticipants, u.AtName()),
		)

		fundraiseid := cuid.Slug()
		rds.SAdd("fundraise:"+fundraiseid, u.Id)
		rds.Expire("fundraise:"+fundraiseid, s.GiveAwayTimeout)
		chattable.BaseChat.ReplyMarkup = fundraiseKeyboard(fundraiseid, receiver.Id, nparticipants, sats)
		bot.Send(chattable)
	case opts["transactions"].(bool):
		// show list of transactions
		limit := 25
		offset := 0
		if page, err := opts.Int("--page"); err == nil {
			offset = limit * (page - 1)
		}

		txns, err := u.listTransactions(limit, offset)
		if err != nil {
			log.Warn().Err(err).Str("user", u.Username).
				Msg("failed to list transactions")
			break
		}

		title := fmt.Sprintf("Latest %d transactions", limit)
		if offset > 0 {
			title = fmt.Sprintf("Transactions from %d to %d", offset+1, offset+limit)
		}

		u.notify(mustache.Render(`<b>{{title}}</b>
{{#txns}}
<code>{{StatusSmall}}</code> <code>{{PaddedSatoshis}}</code> {{Icon}} {{PeerActionDescription}}{{^TelegramPeer.Valid}}<i>{{Description}}</i>{{/TelegramPeer.Valid}} <i>{{TimeFormatSmall}}</i> /tx{{HashReduced}}
{{/txns}}
        `, map[string]interface{}{"title": title, "txns": txns}))
		break
	case opts["balance"].(bool):
		// show balance
		info, err := u.getInfo()
		if err != nil {
			log.Warn().Err(err).Str("user", u.Username).Msg("failed to get info")
			break
		}

		u.notify(fmt.Sprintf(`
<b>Balance</b>: %.3f sat
<b>Total received</b>: %.3f sat
<b>Total sent</b>: %.3f sat
<b>Total fees paid</b>: %.3f sat
        `, info.Balance, info.TotalReceived, info.TotalSent, info.TotalFees))
		break
	case opts["pay"].(bool), opts["withdraw"].(bool), opts["decode"].(bool):
		// pay invoice
		askConfirmation := true
		if opts["now"].(bool) {
			askConfirmation = false
		}

		var bolt11 string
		// when paying, the invoice could be in the message this is replying to
		if ibolt11, ok := opts["<invoice>"]; !ok || ibolt11 == nil {
			if message.ReplyToMessage != nil {
				bolt11, ok = searchForInvoice(*message.ReplyToMessage)
				if !ok {
					u.notify("Invoice not provided.")
					break
				}
			}

			u.notify("Invoice not provided.")
			break
		} else {
			bolt11 = ibolt11.(string)
		}

		optsats, _ := opts.Int("<satoshis>")
		optmsats := optsats * 1000

		if askConfirmation {
			// decode invoice and show a button for confirmation
			inv, nodeAlias, usd, err := decodeInvoice(bolt11)
			if err != nil {
				errMsg := messageFromError(err, "Failed to decode invoice")
				notify(u.ChatId, errMsg)
				break
			}

			amount := int(inv.Get("msatoshi").Int())
			if amount == 0 {
				amount = optmsats
			}

			hash := inv.Get("payment_hash").String()
			text = fmt.Sprintf(`
%d sat (%s)
<i>%s</i>
<b>Hash</b>: %s
<b>Node</b>: %s (%s)
        `,
				amount/1000,
				usd,
				escapeHTML(inv.Get("description").String()),
				hash,
				nodeLink(inv.Get("payee").String()),
				nodeAlias,
			)

			msg := notify(u.ChatId, text)
			id := msg.MessageID

			hashfirstchars := hash[:5]
			rds.Set("payinvoice:"+hashfirstchars, bolt11, s.PayConfirmTimeout)
			rds.Set("payinvoice:"+hashfirstchars+":msats", optmsats, s.PayConfirmTimeout)

			editWithKeyboard(u.ChatId, id,
				text+"\n\nPay the invoice described above?",
				tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Cancel", fmt.Sprintf("cancel=%d", u.Id)),
						tgbotapi.NewInlineKeyboardButtonData("Yes", "pay="+hashfirstchars),
					),
				),
			)
		} else {
			u.payInvoice(message.MessageID, bolt11, optmsats)
		}
		break
	case opts["help"].(bool):
		command, _ := opts.String("<command>")
		handleHelp(u, command)
		break
	case opts["toggle"].(bool):
		if message.Chat.Type == "private" {
			break
		}

		switch {
		case opts["spammy"].(bool):
			if message.Chat.Type == "supergroup" {
				userchatconfig := tgbotapi.ChatConfigWithUser{
					ChatID:             message.Chat.ID,
					SuperGroupUsername: message.Chat.ChatConfig().SuperGroupUsername,
					UserID:             message.From.ID,
				}
				chatmember, err := bot.GetChatMember(userchatconfig)
				if err != nil ||
					(chatmember.Status != "administrator" && chatmember.Status != "creator") {
					log.Warn().Err(err).
						Int64("group", message.Chat.ID).
						Int("user", message.From.ID).
						Msg("toggle impossible. can't get user or not an admin.")
					break
				}
			} else if message.Chat.Type == "group" {
				// ok, everybody can toggle
			} else {
				break
			}

			log.Debug().Int64("group", message.Chat.ID).Msg("toggling spammy")
			spammy, err := toggleSpammy(message.Chat.ID)
			if err != nil {
				log.Warn().Err(err).Msg("failed to toggle spammy")
				break
			}

			if spammy {
				notify(message.Chat.ID, "This group is now spammy.")
			} else {
				notify(message.Chat.ID, "Not spamming anymore.")
			}
		}
	}
}
