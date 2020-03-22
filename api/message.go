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
