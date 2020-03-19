package model

import "github.com/Michael2008S/etherpad4go/store"

type PadMgr struct {
	dbStore store.Store
}

func (p *PadMgr) GetPad(id, text string) Pad {
	return NewPad(id, text)
}

func (p *PadMgr) SavePad() {

}

func DoesPadExist(padId string) {

}

//sanitizedPadId

func RemovePad(padId string) {

}

func IsValidPadId(padId string) {

}
