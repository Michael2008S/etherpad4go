package model

import (
	"log"
	"testing"
)

func TestPad_getKeyRevisionNumber(t *testing.T) {
	p := Pad{}
	result := p.getKeyRevisionNumber(91)
	log.Println(result)
}
