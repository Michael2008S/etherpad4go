package changeset

import (
	"fmt"
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
	for opIter.hasNext() {
		aOp := opIter.Next()
		fmt.Println(aOp)
		fmt.Println(opIter.hasNext())
	}
}
