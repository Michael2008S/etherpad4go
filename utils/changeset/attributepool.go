package changeset

import (
	"github.com/y0ssar1an/q"
	"strings"
)

const (
	AttribKeyAuthor = "author"
)

/*
  An AttributePool maintains a mapping from [key,value] Pairs called
  Attributes to Numbers (unsigened integers) and vice versa. These numbers are
  used to reference Attributes in Changesets.
*/

type AttributePool struct {
	NumToAttrib map[int][]string `json:"numToAttrib"` // e.g. {0: ['foo','bar']}
	AttribToNum map[string]int   `json:"attribToNum"` // e.g. {'foo,bar': 0}
	NextNum     int              `json:"nextNum"`
}

func NewAttributePool() AttributePool {
	return AttributePool{
		NumToAttrib: map[int][]string{},
		AttribToNum: map[string]int{},
		NextNum:     0,
	}
}

func (a *AttributePool) PutAttrib(attrib []string, dontAddIfAbsent bool) int {
	attribToStr := strings.Join(attrib, ",")
	num, found := a.AttribToNum[attribToStr]
	if found {
		return num
	}
	if dontAddIfAbsent {
		return -1
	}
	num = a.NextNum
	a.NextNum++
	a.AttribToNum[attribToStr] = num
	// FIXME  this.numToAttrib[num] = [String(attrib[0] || ''), String(attrib[1] || '')];
	a.NumToAttrib[num] = attrib
	q.Q("PutAttrib:", a)
	return num
}

func (a *AttributePool) GetAttrib(num int) []string {
	return a.NumToAttrib[num]
}

func (a *AttributePool) GetAttribKey() {

}

func (a *AttributePool) GetAttribValue() {

}

func (a *AttributePool) EachAttrib() {

}

func (a *AttributePool) ToJsonAble() {

}

func (a *AttributePool) FromJsonAble() {

}
