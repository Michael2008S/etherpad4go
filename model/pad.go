package model

import (
	"fmt"
	"github.com/Michael2008S/etherpad4go/store"
	"github.com/Michael2008S/etherpad4go/utils/changeset"
	"strings"
)

type Pad struct {
	dbStore store.Store

	Id   string
	Text string
}

func (p *Pad) Init() {
	value, found := p.dbStore.Get([]byte(PadKey + p.Id))
	if found {
		// TODO
		fmt.Println(value)
	} else {
		// this pad doesn't exist, so create it
		cs := changeset.ChangeSet{}
		firstChangeset := cs.MakeSplice("\n", 0, 0, CleanText(p.Text), "", "")
		p.appendRevision(firstChangeset, "")
	}
}

func CleanText(text string) string {
	strings.Replace(text, "\r\n", "\n", 0)
	strings.Replace(text, "\r", "\n", 0)
	strings.Replace(text, "\t", "        ", 0)
	strings.Replace(text, "\xa0", " ", 0)
	return text
}

func (p *Pad) apool() {

}

func (p *Pad) getHeadRevisionNumber() {

}

func (p *Pad) getSavedRevisionsList() {

}

func (p *Pad) getPublicStatus() {

}

func (p *Pad) appendRevision(aChangeset, author string) {
	cs := changeset.ChangeSet{}
	cs.

}

func (p *Pad) saveToDatabase() {

}

func (p *Pad) getLastEdit() {

}

func (p *Pad) getRevisionChangeset() {

}

func (p *Pad) getRevisionAuthor() {

}

func (p *Pad) getRevisionDate() {

}

func (p *Pad) getAllAuthors() {

}

func (p *Pad) getInternalRevisionAText() {

}

func (p *Pad) getRevision() {

}

func (p *Pad) getKeyRevisionNumber() {

}

func (p *Pad) GetText() {

}
func (p *Pad) SetText() {

}
func (p *Pad) appendText() {

}
func (p *Pad) init() {

}
func (p *Pad) copy() {

}

func (p *Pad) remove() {

}

func (p *Pad) setPublicStatus() {

}

func (p *Pad) setPassword() {

}

func (p *Pad) isCorrectPassword() {

}

func (p *Pad) isPasswordProtected() {

}

func (p *Pad) addSavedRevision() {

}

func (p *Pad) getSavedRevisions() {

}

func generateSalt() {

}
