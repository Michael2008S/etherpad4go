package model

import (
	"encoding/json"
	"fmt"
	"github.com/Michael2008S/etherpad4go/store"
	"github.com/Michael2008S/etherpad4go/utils/changeset"
	"github.com/jinzhu/copier"
	"strconv"
	"strings"
	"time"
)

const (
	PadKey         = "pad:"
	PadRevisionKey = ":revs:"
)

type Pad struct {
	dbStore        store.Store             `json:"-"`
	Id             string                  `json:"-"`
	Pool           changeset.AttributePool `json:"pool"`
	AText          changeset.AText         `json:"atext"`
	Head           int                     `json:"head"`
	ChatHead       int                     `json:"chatHead"`
	PublicStatus   bool                    `json:"publicStatus"`
	SavedRevisions []SavedRevision         `json:"savedRevisions"`
}

type SavedRevision struct {
	revNum    int
	savedById int
	label     string
	timestamp int
	id        string
}

type RevData struct {
	dbStore   store.Store `json:"-"`
	changeset string      `json:"changeset"`
	meta      meta        `json:"meta"`
}
type meta struct {
	author    string          `json:"author"`
	timestamp int             `json:"timestamp"`
	aText     changeset.AText `json:"atext"`
}

func (rd *RevData) SaveToDatabase(p Pad) {
	jsonStr, _ := json.Marshal(rd)
	rd.dbStore.Set([]byte(PadKey+p.Id+PadRevisionKey+strconv.Itoa(p.Head)), jsonStr, 0)
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

func (p *Pad) apool() changeset.AttributePool {
	return p.Pool
}

func (p *Pad) getHeadRevisionNumber() int {
	return p.Head
}

func (p *Pad) getSavedRevisionsNumber() int {
	return len(p.SavedRevisions)
}

func (p *Pad) getSavedRevisionsList() int {

}

func (p *Pad) appendRevision(aChangeset, author string) {
	cs := changeset.ChangeSet{}
	newAText := cs.ApplyToAText(aChangeset, p.AText, p.Pool)
	copier.Copy(p.AText, newAText)
	newRevData := RevData{
		changeset: aChangeset,
		meta: meta{
			author:    author,
			timestamp: int(time.Now().Unix()),
		},
	}
	// ex. getNumForAuthor
	if len(author) > 0 {
		p.Pool.PutAttrib("author", true)
	}
	p.Head++
	newRev := p.Head
	if newRev%100 == 0 {
		newRevData.meta.aText = p.AText
	}
	jsonStr, _ := json.Marshal(newRevData)
	p.dbStore.Set([]byte(PadKey+p.Id+PadRevisionKey+strconv.Itoa(newRev)), jsonStr, 0)
	p.saveToDatabase()

}

func (p *Pad) saveToDatabase() {
	jsonStr, _ := json.Marshal(p)
	p.dbStore.Set([]byte(PadKey+p.Id), jsonStr, 0)
}

func (p *Pad) getLastEdit() {
	revNum := p.getHeadRevisionNumber()
	p.dbStore.Get([]byte(PadKey + p.Id + PadRevisionKey + strconv.Itoa(revNum)))
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
