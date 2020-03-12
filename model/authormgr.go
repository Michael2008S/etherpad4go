package model

import (
	"encoding/json"
	"github.com/Michael2008S/etherpad4go/store"
	"github.com/Michael2008S/etherpad4go/utils"
	"math/rand"
	"time"
)

const (
	AuthorKey = "globalAuthor:"
)

var ColorPalette = []string{"#ffc7c7", "#fff1c7", "#e3ffc7", "#c7ffd5", "#c7ffff", "#c7d5ff", "#e3c7ff",
	"#ffc7f1", "#ffa8a8", "#ffe699", "#cfff9e", "#99ffb3", "#a3ffff", "#99b3ff", "#cc99ff", "#ff99e5", "#e7b1b1",
	"#e9dcAf", "#cde9af", "#bfedcc", "#b1e7e7", "#c3cdee", "#d2b8ea", "#eec3e6", "#e9cece", "#e7e0ca", "#d3e5c7",
	"#bce1c5", "#c1e2e2", "#c1c9e2", "#cfc1e2", "#e0bdd9", "#baded3", "#a0f8eb", "#b1e7e0", "#c3c8e4", "#cec5e2",
	"#b1d5e7", "#cda8f0", "#f0f0a8", "#f2f2a6", "#f5a8eb", "#c5f9a9", "#ececbb", "#e7c4bc", "#daf0b2", "#b0a0fd",
	"#bce2e7", "#cce2bb", "#ec9afe", "#edabbd", "#aeaeea", "#c4e7b1", "#d722bb", "#f3a5e7", "#ffa8a8", "#d8c0c5",
	"#eaaedd", "#adc6eb", "#bedad1", "#dee9af", "#e9afc2", "#f8d2a0", "#b3b3e6",}

type AuthorMgr struct {
	store.Store
}

type Author struct {

}

func DoesAuthorExist(authorID string) {

}

func GetAuthor4Token(token string) {

}

// Returns the AuthorID for a mapper. We can map using a mapperkey, so far this is token2author and mapper2author
func CreateAuthorIfNotExistsFor(authorMapper, name string) {

}

func (a *AuthorMgr) CreateAuthor(name string) {
	author := utils.RandStringRunes(16)
	authObj := map[string]interface{}{}
	rand.Seed(time.Now().UnixNano())
	authObj["colorId"] = rand.Intn(len(ColorPalette))
	authObj["name"] = name
	authObj["timestamp"] = time.Now().Unix()
	authObjStr, _ := json.Marshal(authObj)
	a.Store.Set([]byte(AuthorKey+author), authObjStr, 0)
}

func GetAuthor(author string) {

}

func ListPadsOfAuthor(authorID string) {

}

func AddPad(authorID, padID string) {

}

//func RemovePad(authorID, padID string) {
//
//}
