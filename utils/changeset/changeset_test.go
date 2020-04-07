package changeset

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/y0ssar1an/q"
	"strings"
	"testing"
)

func TestChangeSet_Unpack(t *testing.T) {
	chgset := ChangeSet{}
	err := chgset.Unpack("Z:z>1|2=m=b*0|1+1$\n")
	if err != nil {

	}
	fmt.Println(chgset)
	cs := chgset.Pack()
	fmt.Println(cs)
}

func TestOperatorIterator_NewIterator(t *testing.T) {
	chgset := ChangeSet{}
	err := chgset.Unpack("Z:z>1|2=m=b*0|1+1$\n")
	if err != nil {

	}
	opIter := NewOperatorIterator(chgset.Ops, 0)
	fmt.Println(opIter)
	for opIter.HasNext() {
		aOp := opIter.Next()
		fmt.Println(aOp)
		fmt.Println(opIter.HasNext())
	}
}

func TestChangeSet_MakeSplice(t *testing.T) {
	text := `Welcome to Etherpad!

This pad text is synchronized~ https://github.com/ether/etherpad-lite

`
	//text = "Welcome to Etherpad!\\n\\nThis pad text is synchronized~ https:\\/\\/github.com\\/ether\\/etherpad-lite\\n"
	q.Q(text)
	cleanTxt := CleanText(text)
	q.Q("cleanTxt", cleanTxt)
	chgset := ChangeSet{}
	cs := chgset.MakeSplice("\n", 0, 0, cleanTxt, "", nil)

	q.Q(cs)

	atext := AText{
		Text:    "\n",
		Attribs: "|1+1",
	}
	atext = AText{
		Text:    "",
		Attribs: "",
	}
	pool := AttributePool{}
	newText := chgset.ApplyToAText(cs, atext, pool)

	q.Q(newText)

	q.Q("====================test new applyToAText==================")

	reqcs := "Z:2l>1|3=2k*0+1$a"
	reqAtext := chgset.ApplyToAText(reqcs, newText, pool)
	q.Q(reqAtext)
	expTxt := "Welcome to Etherpad!\n\nThis pad text is synchronized~ https://github.com/ether/etherpad-lite\na\n"
	assert.Equal(t, expTxt, reqAtext.Text)
	assert.Equal(t, "|3+2k*0+1|1+1", reqAtext.Attribs)

}

func CleanText(text string) string {
	strings.Replace(text, "\r\n", "\n", 0)
	strings.Replace(text, "\r", "\n", 0)
	strings.Replace(text, "\t", "        ", 0)
	strings.Replace(text, "\xa0", " ", 0)
	return text
}
