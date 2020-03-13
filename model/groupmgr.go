package model

import (
	"github.com/Michael2008S/etherpad4go/store"
	"github.com/Michael2008S/etherpad4go/utils"
)

const (
	GroupKey = "group:"
)

type GroupMgr struct {
	dbStore store.Store
}

func ListAllGroups() {

}

func DeleteGroup(groupID string) {

}

func DoesGroupExist(groupID string) {

}

func (g *GroupMgr) CreateGroup() {
	groupID := "g." + utils.RandStringRunes(16)
	g.dbStore.Set([]byte(GroupKey+groupID), []byte("TODO "), 0)
}

func CreateGroupIfNotExistsFor(groupMapper string) {

}

func CreateGroupPad(groupID, padName, text string) {

}

func ListPads(groupID string) {

}
