package model

import (
	"fmt"
	"github.com/Michael2008S/etherpad4go/store"
)

func CleanText(text string) {

}

type Pad struct {
	dbStore store.Store


	Id   string
	Text string
}

func (p *Pad)Init(){
	value,found := p.dbStore.Get([]byte(PadKey+p.Id))
	if found {

	} else {
		// this pad doesn't exist, so create it
		fmt.Println(value)
	}
}

func (p *Pad) apool() {

}

func (p *Pad) getHeadRevisionNumber() {

}

func (p *Pad) getSavedRevisionsList() {

}

func (p *Pad) getPublicStatus() {

}

func (p *Pad) appendRevision() {

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
