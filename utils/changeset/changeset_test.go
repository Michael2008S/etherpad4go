package changeset

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/y0ssar1an/q"
	"log"
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

func Test_moveOpsToNewPool(t *testing.T) {

	cs := "|1+l*1+2|2+1z*0+2|1+1"

	oldPool := AttributePool{
		NumToAttrib: map[int][]string{
			0: []string{"author", "a.glqITynU8VYvF40s"},
			1: []string{"author", "a.UaSfrktmubohgvYq"},
		},
		AttribToNum: map[string]int{
			"author,a.glqITynU8VYvF40s": 0,
			"author,a.UaSfrktmubohgvYq": 1,
		},
		NextNum: 2,
	}
	newPool := AttributePool{
		NumToAttrib: map[int][]string{},
		AttribToNum: map[string]int{},
		NextNum:     0,
	}
	newCs := moveOpsToNewPool(cs, &oldPool, &newPool)

	assert.Equal(t, "|1+l*0+2|2+1z*1+2|1+1", newCs)
	log.Println(fmt.Sprintf("%+v", oldPool))
	log.Println(fmt.Sprintf("%+v", newPool))
}

func Test_moveOpsToNewPool2(t *testing.T) {

	cs := "|1+l*0+1*1+1|1+1+5*2+3|1+1q*0+1*1+3|1+1"

	oldPool := AttributePool{
		NumToAttrib: map[int][]string{
			0: []string{"author", "a.glqITynU8VYvF40s"},
			1: []string{"author", "a.UaSfrktmubohgvYq"},
			2: []string{"strikethrough", "true"},
		},
		AttribToNum: map[string]int{
			"author,a.glqITynU8VYvF40s": 0,
			"author,a.UaSfrktmubohgvYq": 1,
			"strikethrough,true":        2,
		},
		NextNum: 3,
	}
	newPool := AttributePool{
		NumToAttrib: map[int][]string{},
		AttribToNum: map[string]int{},
		NextNum:     0,
	}
	newCs := moveOpsToNewPool(cs, &oldPool, &newPool)

	assert.Equal(t, "|1+l*0+1*1+1|1+1+5*2+3|1+1q*0+1*1+3|1+1", newCs)
	log.Println(fmt.Sprintf("%+v", oldPool))
	log.Println(fmt.Sprintf("%+v", newPool))
}
