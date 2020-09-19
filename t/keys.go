package t

type Key string

const (
	NO         Key = "No"
	YES            = "Yes"
	CANCEL         = "Cancel"
	CANCELED       = "Canceled"
	COMPLETED      = "Completed"
	CONFIRM        = "Confirm"
	PAYAMOUNT      = "PayAmount"
	FAILURE        = "Failure"
	PROCESSING     = "Processing"
	WITHDRAW       = "Withdraw"
	ERROR          = "Error"
	CHECKING       = "Checking"
	TXPENDING      = "TxPending"
	TXCANCELED     = "TxCanceled"
	UNEXPECTED     = "Unexpected"

	CLAIMFAILED = "ClaimFailed"

	CALLBACKCOINFLIPWINNER = "CallbackCoinflipWinner"
	CALLBACKWINNER         = "CallbackWinner"
	CALLBACKERROR          = "CallbackError"
	CALLBACKEXPIRED        = "CallbackExpired"
	CALLBACKATTEMPT        = "CallbackAttempt"
	CALLBACKSENDING        = "CallbackSending"

	INLINEINVOICERESULT  = "InlineInvoiceResult"
	INLINEGIVEAWAYRESULT = "InlineGiveawayResult"
	INLINEGIVEFLIPRESULT = "InlineGiveflipResult"
	INLINECOINFLIPRESULT = "InlineCoinflipResult"
	INLINEHIDDENRESULT   = "InlineHiddenResult"

	LNURLUNSUPPORTED      = "LnurlUnsupported"
	LNURLERROR            = "LnurlError"
	LNURLAUTHSUCCESS      = "LnurlAuthSuccess"
	LNURLPAYPROMPT        = "LnurlPayPrompt"
	LNURLPAYPROMPTCOMMENT = "LnurlPayPromptComment"
	LNURLPAYSUCCESS       = "LnurlPaySuccess"
	LNURLPAYMETADATA      = "LnurlPayMetadata"
	LNURLPAYCOMMENT       = "LnurlPayComment"

	USERALLOWED       = "UserAllowed"
	SPAMFILTERMESSAGE = "SpamFilterMessage"

	RENAMABLEMSG      = "RenamableMsg"
	RENAMEPROMPT      = "RenamePrompt"
	GROUPNOTRENAMABLE = "GroupNotRenamable"

	PAYMENTFAILED       = "PaymentFailed"
	PAIDMESSAGE         = "PaidMessage"
	DBERROR             = "DBError"
	INSUFFICIENTBALANCE = "InsufficientBalance"
	RATELIMIT           = "RateLimit"
	OVERQUOTA           = "OverQuota"

	PAYMENTRECEIVED      = "PaymentReceived"
	FAILEDTOSAVERECEIVED = "FailedToSaveReceived"

	SPAMMYMSG           = "SpammyMsg"
	COINFLIPSENABLEDMSG = "CoinflipsEnabledMsg"
	LANGUAGEMSG         = "LanguageMsg"
	TICKETMSG           = "TicketMsg"
	FREEJOIN            = "FreeJoin"

	APPBALANCE = "AppBalance"

	HELPINTRO   = "HelpIntro"
	HELPSIMILAR = "HelpSimilar"
	HELPMETHOD  = "HelpMethod"

	RECEIVEHELP = "receiveHelp"

	PAYHELP = "payHelp"

	SENDHELP = "sendHelp"

	TRANSACTIONSHELP = "transactionsHelp"

	BALANCEHELP = "balanceHelp"

	GIVEAWAYHELP            = "giveawayHelp"
	GIVEAWAYMSG             = "GiveAwayMsg"
	GIVEAWAYCLAIM           = "GiveAwayClaim"
	GIVEAWAYSATSGIVENPUBLIC = "GiveawaySatsGivenPublic"

	COINFLIPHELP      = "coinflipHelp"
	COINFLIPWINNERMSG = "CoinflipWinnerMsg"
	COINFLIPGIVERMSG  = "CoinflipGiverMsg"
	COINFLIPAD        = "CoinflipAd"
	COINFLIPJOIN      = "CoinflipJoin"

	GIVEFLIPHELP      = "giveflipHelp"
	GIVEFLIPMSG       = "GiveFlipMsg"
	GIVEFLIPWINNERMSG = "GiveflipWinnerMsg"
	GIVEFLIPAD        = "GiveflipAd"
	GIVEFLIPJOIN      = "GiveflipJoin"

	FUNDRAISEHELP        = "fundraiseHelp"
	FUNDRAISEAD          = "FundraiseAd"
	FUNDRAISEJOIN        = "FundraiseJoin"
	FUNDRAISECOMPLETE    = "FundraiseComplete"
	FUNDRAISERECEIVERMSG = "FundraiseReceiverMsg"
	FUNDRAISEGIVERMSG    = "FundraiseGiverMsg"

	LIGHTNINGATMHELP       = "lightningatmHelp"
	BLUEWALLETHELP         = "bluewalletHelp"
	APIPASSWORDUPDATEERROR = "APIPasswordUpdateError"
	APICREDENTIALS         = "APICredentials"

	HIDEHELP             = "hideHelp"
	REVEALHELP           = "revealHelp"
	HIDDENREVEALBUTTON   = "HiddenRevealButton"
	HIDDENDEFAULTPREVIEW = "HiddenDefaultPreview"
	HIDDENWITHID         = "HiddenWithId"
	HIDDENSOURCEMSG      = "HiddenSourceMsg"
	HIDDENREVEALMSG      = "HiddenRevealMsg"
	HIDDENMSGNOTFOUND    = "HiddenMsgNotFound"
	HIDDENSHAREBTN       = "HiddenShareBtn"

	ETLENEUMHELP          = "etleneumHelp"
	ETLENEUMACCOUNT       = "EtleneumAccount"
	ETLENEUMHISTORY       = "EtleneumHistory"
	ETLENEUMCONTRACT      = "EtleneumContract"
	ETLENEUMCONTRACTSTATE = "EtleneumContractState"
	ETLENEUMCALL          = "EtleneumCall"
	ETLENEUMCONTRACTS     = "EtleneumContracts"
	ETLENEUMSUBSCRIBED    = "EtleneumSubscribed"
	ETLENEUMCONTRACTEVENT = "EtleneumContractEvent"

	MICROBETHELP                = "microbetHelp"
	MICROBETBETHEADER           = "MicrobetBetHeader"
	MICROBETPAIDBUTNOTCONFIRMED = "MicrobetPaidButNotConfirmed"
	MICROBETPLACING             = "MicrobetPlacing"
	MICROBETPLACED              = "MicrobetPlaced"
	MICROBETLIST                = "MicrobetList"

	BITREFILLHELP            = "bitrefillHelp"
	BITREFILLINVENTORYHEADER = "BitrefillInventoryHeader"
	BITREFILLPACKAGESHEADER  = "BitrefillPackagelHeader"
	BITREFILLNOPROVIDERS     = "BitrefillNoProviders"
	BITREFILLCONFIRMATION    = "BitrefillConfirmation"
	BITREFILLFAILEDSAVE      = "BitrefillFailedToSave"
	BITREFILLPURCHASEDONE    = "BitrefillPurchaseDone"
	BITREFILLPURCHASEFAILED  = "BitrefillPurchaseFailed"
	BITREFILLCOUNTRYSET      = "BitrefillCountrySet"
	BITREFILLINVALIDCOUNTRY  = "BitrefillInvalidCountry"

	SATELLITEHELP = "satelliteHelp"
	SATELLITEPAID = "SatellitePaid"

	FUNDBTCHELP   = "fundbtcHelp"
	FUNDBTCFAIL   = "fundbtcFail"
	FUNDBTCFINISH = "fundbtcFinish"

	BITCLOUDSHELP           = "bitcloudsHelp"
	BITCLOUDSCREATEHEADER   = "BitcloudsCreateHeader"
	BITCLOUDSCREATED        = "BitcloudsCreated"
	BITCLOUDSSTOPPEDWAITING = "BitcloudsStoppedWaiting"
	BITCLOUDSNOHOSTS        = "BitcloudsNoHosts"
	BITCLOUDSHOSTSHEADER    = "BitcloudsHostsHeader"
	BITCLOUDSSTATUS         = "BitcloudsStatus"
	BITCLOUDSREMINDER       = "BitcloudsReminder"

	SKYPEHELP = "skypeHelp"
	RUBHELP   = "rubHelp"

	GIFTSHELP       = "giftsHelp"
	GIFTSCREATED    = "GiftsCreated"
	GIFTSFAILEDSAVE = "GiftsFailedSave"
	GIFTSLIST       = "GiftsList"
	GIFTSSPENTEVENT = "GiftsSpentEvent"

	SATS4ADSHELP       = "sats4adsHelp"
	SATS4ADSTOGGLE     = "Sats4adsToggle"
	SATS4ADSBROADCAST  = "Sats4adsBroadcast"
	SATS4ADSSTART      = "Sats4adsStart"
	SATS4ADSPRICETABLE = "Sats4adsPriceTable"
	SATS4ADSADFOOTER   = "Sats4adsAdFooter"
	SATS4ADSVIEWED     = "Viewed"

	ETLENEUMFAILEDTOPAY = "EtleneumFailedToPay"

	TOGGLEHELP = "toggleHelp"

	HELPHELP = "helpHelp"

	STOPHELP = "stopHelp"

	PAYPROMPT          = "PayPrompt"
	FAILEDDECODE       = "FailedDecode"
	BALANCEMSG         = "BalanceMsg"
	TAGGEDBALANCEMSG   = "TaggedBalanceMsg"
	FAILEDUSER         = "FailedUser"
	LOTTERYMSG         = "LotteryMsg"
	INVALIDPARTNUMBER  = "InvalidPartNumber"
	INVALIDAMOUNT      = "InvalidAmount"
	USERSENTTOUSER     = "UserSentToUser"
	USERSENTYOUSATS    = "UserSentYouSats"
	RECEIVEDSATSANON   = "ReceivedSatsAnon"
	FAILEDSEND         = "FailedSend"
	QRCODEFAIL         = "QRCodeFail"
	SAVERECEIVERFAIL   = "SaveReceiverFail"
	CANTSENDNORECEIVER = "CantSendNoReceiver"
	GIVERCANTJOIN      = "GiverCantJoin"
	CANTJOINTWICE      = "CantJoinTwice"
	CANTCANCEL         = "CantCancel"
	FAILEDINVOICE      = "FailedInvoice"
	INVALIDAMT         = "InvalidAmt"
	STOPNOTIFY         = "StopNotify"
	WELCOME            = "Welcome"
	WRONGCOMMAND       = "WrongCommand"
	RETRACTQUESTION    = "RetractQuestion"
	RECHECKPENDING     = "RecheckPending"

	TXNOTFOUND = "TxNotFound"
	TXINFO     = "TxInfo"
	TXLIST     = "TxList"
	TXLOG      = "TxLog"

	TUTORIALWALLET  = "TutorialWallet"
	TUTORIALBLUE    = "TutorialBlue"
	TUTORIALAPPS    = "TutorialApps"
	TUTORIALTWITTER = "TutorialTwitter"
)
