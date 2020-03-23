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
