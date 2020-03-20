package api

import "github.com/Michael2008S/etherpad4go/utils/changeset"

const (
	MsgTypeClientVars = "CLIENT_VARS"
)

type WarpMsgResp struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type ClientVarsDataResp struct {
	SkinName     string `json:"skinName"`
	AccountPrivs struct {
		MaxRevisions int `json:"maxRevisions"`
	} `json:"accountPrivs"`
	AutomaticReconnectionTimeout int           `json:"automaticReconnectionTimeout"`
	InitialRevisionList          []interface{} `json:"initialRevisionList"`
	InitialOptions               struct {
		GuestPolicy string `json:"guestPolicy"`
	} `json:"initialOptions"`
	SavedRevisions     []interface{}      `json:"savedRevisions"`
	CollabClientVars   CollabClientVars   `json:"collab_client_vars"`
	ColorPalette       []string           `json:"colorPalette"`
	ClientIP           string             `json:"clientIp"`
	UserIsGuest        bool               `json:"userIsGuest"`
	UserColor          int                `json:"userColor"`
	PadID              string             `json:"padId"`
	PadOptions         PadOptions         `json:"padOptions"`
	PadShortcutEnabled PadShortcutEnabled `json:"padShortcutEnabled"`
	InitialTitle       string             `json:"initialTitle"`
	Opts               struct {
	} `json:"opts"`
	ChatHead                           int     `json:"chatHead"`
	NumConnectedUsers                  int     `json:"numConnectedUsers"`
	ReadOnlyID                         string  `json:"readOnlyId"`
	Readonly                           bool    `json:"readonly"`
	ServerTimestamp                    int64   `json:"serverTimestamp"`
	UserID                             string  `json:"userId"`
	AbiwordAvailable                   string  `json:"abiwordAvailable"`
	SofficeAvailable                   string  `json:"sofficeAvailable"`
	ExportAvailable                    string  `json:"exportAvailable"`
	Plugins                            Plugins `json:"plugins"`
	IndentationOnNewLine               bool    `json:"indentationOnNewLine"`
	ScrollWhenFocusLineIsOutOfViewport struct {
		Percentage struct {
			EditionAboveViewport int `json:"editionAboveViewport"`
			EditionBelowViewport int `json:"editionBelowViewport"`
		} `json:"percentage"`
		Duration                                 int  `json:"duration"`
		ScrollWhenCaretIsInTheLastLineOfViewport bool `json:"scrollWhenCaretIsInTheLastLineOfViewport"`
		PercentageToScrollWhenUserPressesArrowUp int  `json:"percentageToScrollWhenUserPressesArrowUp"`
	} `json:"scrollWhenFocusLineIsOutOfViewport"`
	InitialChangesets []interface{} `json:"initialChangesets"`
}

type CollabClientVars struct {
	InitialAttributedText struct {
		Text    string `json:"text"`
		Attribs string `json:"attribs"`
	} `json:"initialAttributedText"`
	ClientIP             string `json:"clientIp"`
	PadID                string `json:"padId"`
	HistoricalAuthorData struct {
		APy0WdSkbof4TM4DD struct {
			Name    interface{} `json:"name"`
			ColorID int         `json:"colorId"`
		} `json:"a.Py0WdSkbof4tM4DD"`
		AYjK4P2YxGHx8NNgf struct {
			Name    string `json:"name"`
			ColorID string `json:"colorId"`
		} `json:"a.YjK4P2yxGHx8NNgf"`
	} `json:"historicalAuthorData"`
	Apool changeset.AttributePool `json:"apool"`
	Rev   int                     `json:"rev"`
	Time  int64                   `json:"time"`
}

type PadOptions struct {
	NoColors         bool   `json:"noColors"`
	ShowControls     bool   `json:"showControls"`
	ShowChat         bool   `json:"showChat"`
	ShowLineNumbers  bool   `json:"showLineNumbers"`
	UseMonospaceFont bool   `json:"useMonospaceFont"`
	UserName         bool   `json:"userName"`
	UserColor        bool   `json:"userColor"`
	Rtl              bool   `json:"rtl"`
	AlwaysShowChat   bool   `json:"alwaysShowChat"`
	ChatAndUsers     bool   `json:"chatAndUsers"`
	Lang             string `json:"lang"`
}

type PadShortcutEnabled struct {
	AltF9     bool `json:"altF9"`
	AltC      bool `json:"altC"`
	CmdShift2 bool `json:"cmdShift2"`
	Delete    bool `json:"delete"`
	Return    bool `json:"return"`
	Esc       bool `json:"esc"`
	CmdS      bool `json:"cmdS"`
	Tab       bool `json:"tab"`
	CmdZ      bool `json:"cmdZ"`
	CmdY      bool `json:"cmdY"`
	CmdI      bool `json:"cmdI"`
	CmdB      bool `json:"cmdB"`
	CmdU      bool `json:"cmdU"`
	Cmd5      bool `json:"cmd5"`
	CmdShiftL bool `json:"cmdShiftL"`
	CmdShiftN bool `json:"cmdShiftN"`
	CmdShift1 bool `json:"cmdShift1"`
	CmdShiftC bool `json:"cmdShiftC"`
	CmdH      bool `json:"cmdH"`
	CtrlHome  bool `json:"ctrlHome"`
	PageUp    bool `json:"pageUp"`
	PageDown  bool `json:"pageDown"`
}

type Plugins struct {
	Plugins struct {
		EpEtherpadLite struct {
			Parts []struct {
				Name  string `json:"name"`
				Hooks struct {
					CreateServer  string `json:"createServer"`
					RestartServer string `json:"restartServer"`
				} `json:"hooks"`
				Plugin   string `json:"plugin"`
				FullName string `json:"full_name"`
			} `json:"parts"`
			Package struct {
				Name        string `json:"name"`
				Version     string `json:"version"`
				Description string `json:"description"`
				Main        string `json:"main"`
				Scripts     struct {
					Test string `json:"test"`
				} `json:"scripts"`
				Author   string `json:"author"`
				License  string `json:"license"`
				Invalid  bool   `json:"invalid"`
				RealName string `json:"realName"`
				Path     string `json:"path"`
				RealPath string `json:"realPath"`
				Link     string `json:"link"`
				Depth    int    `json:"depth"`
			} `json:"package"`
		} `json:"ep_etherpad-lite"`
	} `json:"plugins"`
	Parts []struct {
		Name  string `json:"name"`
		Hooks struct {
			ExpressCreateServer string `json:"expressCreateServer"`
		} `json:"hooks"`
		Plugin   string `json:"plugin"`
		FullName string `json:"full_name"`
	} `json:"parts"`
}

// 数据格式，参考：
//6384:42["message",{"type":"CLIENT_VARS","data":{"skinName":"no-skin","accountPrivs":{"maxRevisions":100},"automaticReconnectionTimeout":0,"initialRevisionList":[],"initialOptions":{"guestPolicy":"deny"},"savedRevisions":[],"collab_client_vars":{"initialAttributedText":{"text":"Welcome to Etherpad!\n\nThis pad text is synchronized as you type, so that everyone viewing this page sees the same text. This allows you to collaborate seamlessly on documents!\n\nGet involved with Etherpad at http://etherpad.org\n\nasdf\n\nadsf\n\nasdfasdfadsf阿斯顿发阿斯顿发阿斯顿发\nasdfasdfsdaf\nasdf\n给第三方\nad\nasfd\nsafaafdasdf asdf   tehsi. \n\n","attribs":"|5+6b*0|3+7*0+4*1|2+2*1+c*0+c*1|3+j*0+4*1|3+9*0|1+r|1+1"},"clientIp":"127.0.0.1","padId":"q","historicalAuthorData":{"a.Py0WdSkbof4tM4DD":{"name":null,"colorId":36},"a.YjK4P2yxGHx8NNgf":{"name":"","colorId":"#c4e7b1"}},"apool":{"numToAttrib":{"0":["author","a.Py0WdSkbof4tM4DD"],"1":["author","a.YjK4P2yxGHx8NNgf"]},"nextNum":2},"rev":33,"time":1583675281484},"colorPalette":["#ffc7c7","#fff1c7","#e3ffc7","#c7ffd5","#c7ffff","#c7d5ff","#e3c7ff","#ffc7f1","#ffa8a8","#ffe699","#cfff9e","#99ffb3","#a3ffff","#99b3ff","#cc99ff","#ff99e5","#e7b1b1","#e9dcAf","#cde9af","#bfedcc","#b1e7e7","#c3cdee","#d2b8ea","#eec3e6","#e9cece","#e7e0ca","#d3e5c7","#bce1c5","#c1e2e2","#c1c9e2","#cfc1e2","#e0bdd9","#baded3","#a0f8eb","#b1e7e0","#c3c8e4","#cec5e2","#b1d5e7","#cda8f0","#f0f0a8","#f2f2a6","#f5a8eb","#c5f9a9","#ececbb","#e7c4bc","#daf0b2","#b0a0fd","#bce2e7","#cce2bb","#ec9afe","#edabbd","#aeaeea","#c4e7b1","#d722bb","#f3a5e7","#ffa8a8","#d8c0c5","#eaaedd","#adc6eb","#bedad1","#dee9af","#e9afc2","#f8d2a0","#b3b3e6"],"clientIp":"127.0.0.1","userIsGuest":true,"userColor":36,"padId":"q","padOptions":{"noColors":false,"showControls":true,"showChat":true,"showLineNumbers":true,"useMonospaceFont":false,"userName":false,"userColor":false,"rtl":false,"alwaysShowChat":false,"chatAndUsers":false,"lang":"en-gb"},"padShortcutEnabled":{"altF9":true,"altC":true,"cmdShift2":true,"delete":true,"return":true,"esc":true,"cmdS":true,"tab":true,"cmdZ":true,"cmdY":true,"cmdI":true,"cmdB":true,"cmdU":true,"cmd5":true,"cmdShiftL":true,"cmdShiftN":true,"cmdShift1":true,"cmdShiftC":true,"cmdH":true,"ctrlHome":true,"pageUp":true,"pageDown":true},"initialTitle":"Pad: q","opts":{},"chatHead":3,"numConnectedUsers":0,"readOnlyId":"r.c97d8bdd5f1e326473442fc2f55793c3","readonly":false,"serverTimestamp":1583918908724,"userId":"a.Py0WdSkbof4tM4DD","abiwordAvailable":"no","sofficeAvailable":"no","exportAvailable":"no","plugins":{"plugins":{"ep_etherpad-lite":{"parts":[{"name":"express","hooks":{"createServer":"ep_etherpad-lite/node/hooks/express:createServer","restartServer":"ep_etherpad-lite/node/hooks/express:restartServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/express"},{"name":"static","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/static:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/static"},{"name":"i18n","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/i18n:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/i18n"},{"name":"specialpages","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/specialpages:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/specialpages"},{"name":"socketio","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/socketio:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/socketio"},{"name":"apicalls","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/apicalls:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/apicalls"},{"name":"webaccess","hooks":{"expressConfigure":"ep_etherpad-lite/node/hooks/express/webaccess:expressConfigure"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/webaccess"},{"name":"swagger","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/swagger:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/swagger"}],"package":{"name":"etherpad-lite","version":"1.0.0","description":"my customer etherpad","main":"index.js","scripts":{"test":"echo \"Error: no test specified\" && exit 1"},"author":"","license":"ISC","invalid":true,"realName":"ep_etherpad-lite","path":"/Volumes/RamDisk/goEtherpad/etherpad4go/etherpad-lite/node_modules/ep_etherpad-lite","realPath":"/Volumes/RamDisk/goEtherpad/etherpad4go/etherpad-lite/src","link":"/Volumes/RamDisk/goEtherpad/etherpad4go/etherpad-lite/src","depth":1}}},"parts":[{"name":"swagger","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/swagger:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/swagger"},{"name":"webaccess","hooks":{"expressConfigure":"ep_etherpad-lite/node/hooks/express/webaccess:expressConfigure"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/webaccess"},{"name":"apicalls","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/apicalls:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/apicalls"},{"name":"socketio","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/socketio:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/socketio"},{"name":"specialpages","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/specialpages:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/specialpages"},{"name":"i18n","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/i18n:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/i18n"},{"name":"static","hooks":{"expressCreateServer":"ep_etherpad-lite/node/hooks/express/static:expressCreateServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/static"},{"name":"express","hooks":{"createServer":"ep_etherpad-lite/node/hooks/express:createServer","restartServer":"ep_etherpad-lite/node/hooks/express:restartServer"},"plugin":"ep_etherpad-lite","full_name":"ep_etherpad-lite/express"}]},"indentationOnNewLine":true,"scrollWhenFocusLineIsOutOfViewport":{"percentage":{"editionAboveViewport":0,"editionBelowViewport":0},"duration":0,"scrollWhenCaretIsInTheLastLineOfViewport":false,"percentageToScrollWhenUserPressesArrowUp":0},"initialChangesets":[]}}]
