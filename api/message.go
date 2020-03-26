package api

import "github.com/Michael2008S/etherpad4go/utils/changeset"

type CollabRoomReqMessage struct {
	Type      string `json:"type"`
	Component string `json:"component"`
	Data      struct {
		Type      string                  `json:"type"`
		BaseRev   int                     `json:"baseRev"`
		Changeset string                  `json:"changeset"`
		Apool     changeset.AttributePool `json:"apool"`
	} `json:"data"`
}

type ClientReadyReq struct {
	Component       string `json:"component"`
	Type            string `json:"type"`
	PadID           string `json:"padId"`
	SessionID       string `json:"sessionID"`
	Password        string `json:"password"`
	Token           string `json:"token"`
	ProtocolVersion int    `json:"protocolVersion"`
}

type CollabRoomAcceptCommitResp struct {
	Type string `json:"type"`
	Data struct {
		Type   string `json:"type"`
		NewRev int    `json:"newRev"`
	} `json:"data"`
}

type CollabRoomNewChangesResp struct {
	Type string `json:"type"`
	Data struct {
		Type        string                  `json:"type"`
		NewRev      int                     `json:"newRev"`
		Changeset   string                  `json:"changeset"`
		Apool       changeset.AttributePool `json:"apool"`
		Author      string                  `json:"author"`
		CurrentTime int                     `json:"currentTime"`
		TimeDelta   int                     `json:"timeDelta"`
	} `json:"data"`
}

type UserNewInfoResp struct {
	Type string `json:"type"`
	Data struct {
		Type string `json:"type"`
		UserInfo UserInfo `json:"userInfo"`
	} `json:"data"`
}

type UserInfo struct {
	Ip string `json:"ip"`
	ColorId int `json:"colorId"`
	UserAgent string `json:"userAgent"`
	UserId   string `json:"userId"`
}
