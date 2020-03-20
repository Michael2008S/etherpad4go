package model

import (
	"encoding/json"
	"fmt"
	"github.com/Michael2008S/etherpad4go/store"
	"github.com/Michael2008S/etherpad4go/utils/changeset"
	"github.com/jinzhu/copier"
	"log"
	"math"
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

func NewPad(id, text string) Pad {
	p := Pad{Id: id,}
	p.Init(text)
	return p
}

func (rd *RevData) SaveToDatabase(p Pad) {
	jsonStr, _ := json.Marshal(rd)
	rd.dbStore.Set([]byte(PadKey+p.Id+PadRevisionKey+strconv.Itoa(p.Head)), jsonStr, 0)
}

func (p *Pad) Init(text string) {
	value, found := p.dbStore.Get([]byte(PadKey + p.Id))
	json.Unmarshal(value, p)
	// if this pad exists, load it
	if found {
		// copy all attr. To a transfrom via fromJsonable if necassary
		// TODO
		fmt.Println(value)
	} else {
		// this pad doesn't exist, so create it
		cs := changeset.ChangeSet{}
		firstChangeset := cs.MakeSplice("\n", 0, 0, CleanText(text), "", "")
		p.appendRevision(firstChangeset, "")
	}
	// TODO  hooks.callAll("padLoad", { 'pad':  this });
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

func (p *Pad) GetHeadRevisionNumber() int {
	return p.Head
}

func (p *Pad) getSavedRevisionsNumber() int {
	return len(p.SavedRevisions)
}

func (p *Pad) getSavedRevisionsList() int {
	return 0
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

	// TODO
	//if (this.head == 0) {
	//	hooks.callAll("padCreate", {'pad':this, 'author': author});
	//} else {
	//	hooks.callAll("padUpdate", {'pad':this, 'author': author});
	//}
}

func (p *Pad) saveToDatabase() {
	jsonStr, _ := json.Marshal(p)
	p.dbStore.Set([]byte(PadKey+p.Id), jsonStr, 0)
}

func (p *Pad) getLastEdit() RevData {
	revData := RevData{}
	revNum := p.GetHeadRevisionNumber()
	data, _ := p.dbStore.Get([]byte(PadKey + p.Id + PadRevisionKey + strconv.Itoa(revNum)))
	json.Unmarshal(data, &revData)
	return revData
}

func (p *Pad) getRevision(revNum int) (revData RevData) {
	data, _ := p.dbStore.Get([]byte(PadKey + p.Id + PadRevisionKey + strconv.Itoa(revNum)))
	json.Unmarshal(data, &revData)
	return
}

func (p *Pad) getRevisionChangeset(revNum int) string {
	revData := RevData{}
	data, _ := p.dbStore.Get([]byte(PadKey + p.Id + PadRevisionKey + strconv.Itoa(revNum)))
	json.Unmarshal(data, &revData)
	return revData.changeset
}

func (p *Pad) getRevisionAuthor(revNum int) string {
	revData := RevData{}
	data, _ := p.dbStore.Get([]byte(PadKey + p.Id + PadRevisionKey + strconv.Itoa(revNum)))
	json.Unmarshal(data, &revData)
	return revData.meta.author
}

func (p *Pad) getRevisionDate(revNum int) int {
	revData := RevData{}
	data, _ := p.dbStore.Get([]byte(PadKey + p.Id + PadRevisionKey + strconv.Itoa(revNum)))
	json.Unmarshal(data, &revData)
	return revData.meta.timestamp
}

func (p *Pad) getAllAuthors() (authors []string) {
	for _, val := range p.Pool.NumToAttrib {
		//	if val[0] == "author" && len(val[1]) > 0 {
		//		authors = append(authors, val[1])
		//	}
		// TODO
		log.Println(val)
	}
	return
}

func (p *Pad) getInternalRevisionAText(targetRev int) changeset.AText {
	keyRev := p.getKeyRevisionNumber(targetRev)
	// find out which changesets are needed
	var neededChangesets []int
	for curRev := keyRev; curRev < targetRev; curRev++ {
		neededChangesets = append(neededChangesets, curRev)
	}
	// get all needed data out of the database

	// start to get the atext of the key revision
	revData := p.getRevision(targetRev)
	var changesets map[int]string
	for _, v := range neededChangesets {
		data := p.getRevision(v)
		changesets[v] = data.changeset
	}

	// apply all changesets to the key changeset
	chgset := changeset.ChangeSet{}
	atext := revData.meta.aText
	for k, _ := range changesets {
		cs := changesets[k]
		atext = chgset.ApplyToAText(cs, atext, p.Pool)
	}
	return atext
}

func (p *Pad) getKeyRevisionNumber(revNum int) int {
	return int(math.Floor(float64(revNum)/100) * 100)
}

func (p *Pad) GetText() string {
	return p.AText.Text
}
func (p *Pad) SetText(newText string) {
	newText = CleanText(newText)
	oldText := p.GetText()
	newTxtLen := len(newText)
	chgset := changeset.ChangeSet{}
	// create the changeset
	// We want to ensure the pad still ends with a \n, but otherwise keep
	// getText() and setText() consistent.
	var cs string
	if string([]rune(newText)[newTxtLen-1:newTxtLen]) == "\n" {
		cs = chgset.MakeSplice(oldText, 0, len(oldText), newText, "", "")
	} else {
		cs = chgset.MakeSplice(oldText, 0, len(oldText)-1, newText, "", "")
	}
	// append the changeset
	p.appendRevision(cs, "")
}
func (p *Pad) appendText(newText string) {
	newText = CleanText(newText)
	oldText := p.GetText()
	chgset := changeset.ChangeSet{}
	cs := chgset.MakeSplice(oldText, len(oldText), 0, newText, "", "")
	p.appendRevision(cs, "")
}

//appendChatMessage
//getChatMessage
// getChatMessages

//func (p *Pad) copy() {
//}

//func (p *Pad) remove() {
//	padID := p.Id
//	// kick everyone from this pad
//
//}

//func (p *Pad) setPublicStatus() {
//
//}
//
//func (p *Pad) setPassword() {
//
//}

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
