package changeset

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/y0ssar1an/q"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type ChangeSet struct {
	// length of the text before changeset
	OldLen int
	// length of the text after changeset
	NewLen int

	Ops      string
	CharBank string
}

type OperatorIterator struct {
	OpsStr     string
	regex      string
	startIndex int
	curIndex   int
	prevIndex  int

	rgxResult []string

	Operator Operator
}

type Operator struct {
	OpCode  string
	Chars   int
	Lines   int
	Attribs string
}

/**
 * this function creates an iterator which decodes string changeset operations
 * @param opsStr {string} String encoding of the change operations to be performed
 * @param optStartIndex {int} from where in the string should the iterator start
 * @return {Op} type object iterator
 */
func NewOperatorIterator(opsStr string, optStartIndex int) *OperatorIterator {
	opIter := OperatorIterator{}
	opIter.OpsStr = opsStr
	opIter.regex = `((?:\*[0-9a-z]+)*)(?:\|([0-9a-z]+))?([-+=])([0-9a-z]+)|\?|`
	opIter.startIndex = optStartIndex
	opIter.curIndex = opIter.startIndex
	opIter.prevIndex = opIter.curIndex

	opIter.nextRegexMatch()

	return &opIter
}

func (opIter *OperatorIterator) nextRegexMatch() error {
	opIter.prevIndex = opIter.curIndex
	reg := regexp.MustCompile(opIter.regex)
	opsStr := SubString(opIter.OpsStr, opIter.prevIndex, len(opIter.OpsStr))
	opIter.rgxResult = reg.FindStringSubmatch(opsStr)
	q.Q("nextRegexMatch=>", opIter.OpsStr, opIter.rgxResult, opIter.curIndex)
	if len(opIter.rgxResult) > 0 {
		if opIter.rgxResult[0] == "?" {
			return errors.New("Hit error opcode in op stream.")
		}
		opIter.curIndex = opIter.curIndex + len(opIter.rgxResult[0])
	}
	return nil
}

func (opIter *OperatorIterator) Next() Operator {
	op := Operator{}
	if len(opIter.rgxResult) > 0 && len(opIter.rgxResult[0]) > 0 {
		op.Attribs = opIter.rgxResult[1]
		lines, _ := strconv.ParseInt(opIter.rgxResult[2], 36, 32)
		op.Lines = int(lines)
		op.OpCode = opIter.rgxResult[3]
		chars, _ := strconv.ParseInt(opIter.rgxResult[4], 36, 32)
		op.Chars = int(chars)
		opIter.nextRegexMatch()
	}
	return op
}

func (opIter *OperatorIterator) HasNext() bool {
	return len(opIter.rgxResult) > 0 && len(opIter.rgxResult[0]) > 0
}

func (opIter *OperatorIterator) lastIndex() int {
	return opIter.prevIndex
}

/**
 * Writes the Op in a string the way that changesets need it
 */
func opString() {

}

/**
 * Used just for debugging
 */
func stringOp() {

}

/**
 * Used to check if a Changeset if valid
 * @param cs {Changeset} Changeset to be checked
 */
func (chgset *ChangeSet) CheckRep(cs string) (err error) {
	// doesn't check things that require access to attrib pool (e.g. attribute order)
	// or original string (e.g. newline positions)
	chgset.Unpack(cs)
	oldLen := chgset.OldLen
	newLen := chgset.NewLen
	ops := chgset.Ops
	charBank := chgset.CharBank

	assem := smartOpAssembler{}
	oldPos := 0
	calcNewLen := 0
	numInserted := 0
	iter := NewOperatorIterator(ops, 0)
	for iter.HasNext() {
		op := iter.Next()
		switch op.OpCode {
		case "=":
			oldPos += op.Chars
			calcNewLen += op.Chars
		case "-":
			oldPos += op.Chars
			// exports.assert(oldPos <= oldLen, oldPos, " > ", oldLen, " in ", cs);
			if oldPos <= oldPos {
				err = errors.New(fmt.Sprintf("%d > %d in %s", oldPos, oldLen, cs))
			}
		case "+":
			calcNewLen += op.Chars
			numInserted += op.Chars
			// exports.assert(calcNewLen <= newLen, calcNewLen, " > ", newLen, " in ", cs);
			if calcNewLen <= newLen {
				err = errors.New(fmt.Sprintf("%d > %d in %s", calcNewLen, newLen, cs))
			}
		}
		assem.append(op)
	}

	calcNewLen += (oldLen - oldPos)
	charBank = SubString(charBank, 0, numInserted)
	for len(charBank) < numInserted {
		charBank += "?"
	}

	assem.endDocument()
	normalized := chgset.Pack()
	// exports.assert(normalized == cs, 'Invalid changeset (checkRep failed)');
	if normalized == cs {
		err = errors.New("Invalid changeset (checkRep failed)")
	}
	return
}

/**
 * ==================== Util Functions =======================
 */

/**
 * creates an object that allows you to append operations (type Op) and also
 * compresses them if possible
 */
type smartOpAssembler struct {
	// Like opAssembler but able to produce conforming exportss
	// from slightly looser input, at the cost of speed.
	// Specifically:
	// - merges consecutive operations that can be merged
	// - strips final "="
	// - ignores 0-length changes
	// - reorders consecutive + and - (which margingOpAssembler doesn't do)

	minusAssem mergingOpAssembler
	plusAssem  mergingOpAssembler
	keepAssem  mergingOpAssembler
	assem      stringAssembler

	lastOpcode   string
	lengthChange int
}

func (soa *smartOpAssembler) flushKeeps() {
	soa.assem.append(soa.keepAssem.toString())
	soa.keepAssem.clear()
}

func (soa *smartOpAssembler) flushPlusMinus() {
	soa.assem.append(soa.minusAssem.toString())
	soa.minusAssem.clear()
	soa.assem.append(soa.plusAssem.toString())
	soa.plusAssem.clear()
}

func (soa *smartOpAssembler) append(op Operator) {
	if len(op.OpCode) <= 0 || op.Chars <= 0 {
		return
	}
	if op.OpCode == "-" {
		if soa.lastOpcode == "=" {
			soa.flushKeeps()
		}
		soa.minusAssem.append(op)
		soa.lengthChange -= op.Chars
	} else if op.OpCode == "+" {
		if soa.lastOpcode == "=" {
			soa.flushKeeps()
		}
		soa.plusAssem.append(op)
		soa.lengthChange += op.Chars
	} else if op.OpCode == "=" {
		if soa.lastOpcode == "=" {
			soa.flushPlusMinus()
		}
		soa.keepAssem.append(op)
	}
	soa.lastOpcode = op.OpCode
}

func (soa *smartOpAssembler) appendOpWithText(opCode, text, attribs string, pool *AttributePool) {
	op := Operator{}
	op.OpCode = opCode
	op.Attribs = makeAttribsString(opCode, attribs, pool)
	lastNewlinePos := strings.LastIndex(text, "\n")
	if lastNewlinePos < 0 {
		op.Chars = len(text)
		op.Lines = 0
		soa.append(op)
	} else {
		op.Chars = lastNewlinePos + 1
		rgx := regexp.MustCompile(`\n`)
		findStrs := rgx.Split(text, -1)
		q.Q(findStrs)
		op.Lines = len(findStrs) - 1
		soa.append(op)
		op.Chars = len(text) - (lastNewlinePos + 1)
		op.Lines = 0
		soa.append(op)
	}
}

func (soa *smartOpAssembler) toString() string {
	soa.flushPlusMinus()
	soa.flushKeeps()
	return soa.assem.toString()
}

func (soa *smartOpAssembler) clear() {
	soa.minusAssem.clear()
	soa.plusAssem.clear()
	soa.keepAssem.clear()
	soa.assem.clear()
	soa.lengthChange = 0
}

func (soa *smartOpAssembler) endDocument() {
	soa.keepAssem.endDocument()
}

func (soa *smartOpAssembler) GetLengthChange() int {
	return soa.lengthChange
}

/**
 * A custom made StringBuffer
 */
type stringAssembler []string

func (sa *stringAssembler) append(s string) {
	*sa = append(*sa, s)
}

func (sa *stringAssembler) toString() string {
	return strings.Join(*sa, "")
}

func (sa *stringAssembler) clear() {
	*sa = []string{}
}

//func mergingOpAssembler() {
//	assem := opAssembler()
//	bufOp := newOp()
//}

type mergingOpAssembler struct {
	// This assembler can be used in production; it efficiently
	// merges consecutive operations that are mergeable, ignores
	// no-ops, and drops final pure "keeps".  It does not re-order
	// operations.

	assem OperatorAssembler
	bufOp Operator

	// If we get, for example, insertions [xxx\n,yyy], those don't merge,
	// but if we get [xxx\n,yyy,zzz\n], that merges to [xxx\nyyyzzz\n].
	// This variable stores the length of yyy and any other newline-less
	// ops immediately after it.
	bufOpAdditionalCharsAfterNewline int
}

func (moa *mergingOpAssembler) flush(isEndDocument bool) {
	if moa.bufOp.OpCode != "" {
		if isEndDocument && moa.bufOp.OpCode == "=" && moa.bufOp.Attribs != "" {
			// final merged keep, leave it implicit
		} else {
			moa.assem.Append(moa.bufOp)
			if moa.bufOpAdditionalCharsAfterNewline > 0 {
				moa.bufOp.Chars = moa.bufOpAdditionalCharsAfterNewline
				moa.bufOp.Lines = 0
				moa.assem.Append(moa.bufOp)
				moa.bufOpAdditionalCharsAfterNewline = 0
			}
		}
		moa.bufOp.OpCode = ""
	}
}

func (moa *mergingOpAssembler) append(op Operator) {
	if op.Chars > 0 {
		if moa.bufOp.OpCode == op.OpCode && moa.bufOp.Attribs == op.Attribs {
			if op.Lines > 0 {
				// bufOp and additional chars are all mergeable into a multi-line op
				moa.bufOp.Chars += moa.bufOpAdditionalCharsAfterNewline + op.Chars
				moa.bufOp.Lines = op.Lines
				moa.bufOpAdditionalCharsAfterNewline = 0
			} else if moa.bufOp.Lines == 0 {
				// both bufOp and op are in-line
				moa.bufOp.Chars += op.Chars
			} else {
				// append in-line text to multi-line bufOp
				moa.bufOpAdditionalCharsAfterNewline += op.Chars
			}
		} else {
			moa.flush(false)
			copier.Copy(&moa.bufOp, op)
		}
	}
}

func (moa *mergingOpAssembler) endDocument() {
	moa.flush(true)
}

func (moa *mergingOpAssembler) toString() string {
	moa.flush(false)
	return moa.assem.toString()
}

func (moa *mergingOpAssembler) clear() {
	moa.assem.clear()
	moa.bufOp = Operator{}
}

// this function allows op to be mutated later (doesn't keep a ref)
type OperatorAssembler []string

func (oa *OperatorAssembler) Append(op Operator) {
	*oa = append(*oa, op.Attribs)
	if op.Lines > 0 {
		*oa = append(*oa, "|", strconv.FormatInt(int64(op.Lines), 36))
	}
	*oa = append(*oa, op.OpCode, strconv.FormatInt(int64(op.Chars), 36))
}

func (oa *OperatorAssembler) toString() string {
	return strings.Join(*oa, "")
}

func (oa *OperatorAssembler) clear() {
	*oa = OperatorAssembler{}
}

/**
 * A custom made String Iterator
 * @param str {string} String to be iterated over
 */
type stringIterator struct {
	str      string
	curIndex int
	newLines int
}

func NewStringIterator(str string) stringIterator {
	si := stringIterator{}
	si.str = str
	si.newLines = len(strings.Split(str, "\n")) - 1

	return si
}

func (si *stringIterator) take(n int) string {
	if err := si.assertRemaining(n); err != nil {
		log.Println(err)
		return ""
	}
	s := SubStrLen(si.str, si.curIndex, n)
	si.newLines = len(strings.Split(s, "\n")) - 1
	si.curIndex += n
	return s
}

func (si *stringIterator) peek(n int) string {
	if err := si.assertRemaining(n); err != nil {
		log.Println(err)
		return ""
	}
	return SubStrLen(si.str, si.curIndex, n)
}

func (si *stringIterator) skin(n int) {
	if err := si.assertRemaining(n); err != nil {
		log.Println(err)
		return
	}
	si.curIndex += n
}

func (si *stringIterator) remaining() int {
	return len(si.str) - si.curIndex
}

func (si *stringIterator) assertRemaining(n int) error {
	if !(n <= si.remaining()) {
		return errors.New(fmt.Sprintf("assertRemaining: !(%d <= %d )", n, si.remaining()))
	}
	return nil
}

// @description Unpacks a string encoded Changeset into a proper Changeset object
// @param cs {string} String encoded Changeset
// @return {Changeset} a Changeset class
func (chgset *ChangeSet) Unpack(cs string) error {
	reg := regexp.MustCompile("Z:([0-9a-z]+)([><])([0-9a-z]+)|")
	headerMatch := reg.FindStringSubmatch(cs)
	if len(headerMatch) < 4 || len(headerMatch[0]) <= 0 {
		return errors.New("Unpack change set error.")
	}
	oldLen, _ := strconv.ParseInt(headerMatch[1], 36, 32)
	chgset.OldLen = int(oldLen)
	changeSign := -1
	if headerMatch[2] == ">" {
		changeSign = 1
	}
	changeMag, _ := strconv.ParseInt(headerMatch[3], 36, 32)
	chgset.NewLen = int(oldLen) + changeSign*int(changeMag)
	opsStart := len(headerMatch[0])
	opsEnd := strings.Index(cs, "$")
	if opsEnd < 0 {
		opsEnd = len(cs)
	}
	chgset.Ops = SubString(cs, opsStart, opsEnd)
	chgset.CharBank = SubString(cs, opsEnd+1, len(cs))
	return nil
}

/**
 * This class allows to iterate and modify texts which have several lines
 * It is used for applying Changesets on arrays of lines
 * Note from prev docs: "lines" need not be an array as long as it supports certain calls (lines_foo inside).
 */
func textLinesMutator() {

}

/**
 * Function allowing iterating over two Op strings.
 * @params in1 {string} first Op string
 * @params idx1 {int} integer where 1st iterator should start
 * @params in2 {string} second Op string
 * @params idx2 {int} integer where 2nd iterator should start
 * @params func {function} which decides how 1st or 2nd iterator
 *         advances. When opX.opcode = 0, iterator X advances to
 *         next element
 *         func has signature f(op1, op2, opOut)
 *             op1 - current operation of the first iterator
 *             op2 - current operation of the second iterator
 *             opOut - result operator to be put into Changeset
 * @return {string} the integrated changeset
 */

func applyZip(in1, in2 string, idx1, idx2 int, aFunc func(*Operator, *Operator, *Operator)) string {
	q.Q(in1, in2)
	iter1 := NewOperatorIterator(in1, idx1)
	iter2 := NewOperatorIterator(in2, idx2)
	assem := smartOpAssembler{}
	op1 := Operator{}
	op2 := Operator{}
	opOut := Operator{}
	for len(op1.OpCode) > 0 || iter1.HasNext() || len(op2.OpCode) > 0 || iter2.HasNext() {
		if len(op1.OpCode) <= 0 && iter1.HasNext() {
			op1 = iter1.Next()
		}
		if len(op2.OpCode) <= 0 && iter2.HasNext() {
			op2 = iter2.Next()
		}

		aFunc(&op1, &op2, &opOut)
		if len(opOut.OpCode) > 0 {
			//print(opOut.toSource());
			//tmpOp := Operator{}
			//copier.Copy(&tmpOp, opOut)
			//assem.append(tmpOp)
			assem.append(opOut)
			opOut.OpCode = ""
		}
		//q.Q(op1, op2, opOut, assem.toString())
		//q.Q(iter1.HasNext(),iter2.HasNext())
		//q.Q("--------")
	}
	assem.endDocument()
	q.Q(assem.toString())
	return assem.toString()
}

/**
 * Packs Changeset object into a string
 * @params oldLen {int} Old length of the Changeset
 * @params newLen {int] New length of the Changeset
 * @params opsStr {string} String encoding of the changes to be made
 * @params bank {string} Charbank of the Changeset
 * @returns {Changeset} a Changeset class
 */
func (chgset *ChangeSet) Pack() string {
	lenDiff := chgset.NewLen - chgset.OldLen
	lenDiffStr := "<" + strconv.FormatInt(int64(-lenDiff), 36)
	if lenDiff > 0 {
		lenDiffStr = ">" + strconv.FormatInt(int64(lenDiff), 36)
	}
	return fmt.Sprintf("Z:%s%s%s$%s", strconv.FormatInt(int64(chgset.OldLen), 36), lenDiffStr,
		chgset.Ops, chgset.CharBank)
}

/**
 * Applies a Changeset to a string
 * @params cs {string} String encoded Changeset
 * @params str {string} String to which a Changeset should be applied
 */
func (chgset *ChangeSet) ApplyToText(cs, str string) (string, error) {
	chgset.Unpack(cs)
	//FIXME exports.assert(str.length == unpacked.oldLen, "mismatched apply: ", str.length, " / ", unpacked.oldLen);
	csIter := NewOperatorIterator(chgset.Ops, 0)
	bankIter := NewStringIterator(chgset.CharBank)
	strIter := NewStringIterator(str)
	assem := stringAssembler{}

	for csIter.HasNext() {
		op := csIter.Next()
		switch op.OpCode {
		case "+":
			//op is + and op.lines 0: no newlines must be in op.chars
			//op is + and op.lines >0: op.chars must include op.lines newlines
			if op.Lines != len(strings.Split(bankIter.peek(op.Chars), "\n"))-1 {
				return "", errors.New(fmt.Sprintf("newline count is wrong in op +; cs:%s and text:%s ", cs, str))
			}
			assem.append(bankIter.take(op.Chars))
		case "-":
			//op is - and op.lines 0: no newlines must be in the deleted string
			//op is - and op.lines >0: op.lines newlines must be in the deleted string
			if op.Lines != len(strings.Split(strIter.peek(op.Chars), "\n"))-1 {
				return "", errors.New(fmt.Sprintf("newline count is wrong in op -; cs:%s and text:%s ", cs, str))
			}
			strIter.skin(op.Chars)
		case "=":
			//op is = and op.lines 0: no newlines must be in the copied string
			//op is = and op.lines >0: op.lines newlines must be in the copied string
			if op.Lines != len(strings.Split(strIter.peek(op.Chars), "\n"))-1 {
				return "", errors.New(fmt.Sprintf("newline count is wrong in op =; cs:%s and text:%s ", cs, str))
			}
			assem.append(strIter.take(op.Chars))
		}
	}
	assem.append(strIter.take(strIter.remaining()))
	return assem.toString(), nil
}

/**
 * applies a changeset on an array of lines
 * @param CS {Changeset} the changeset to be applied
 * @param lines The lines to which the changeset needs to be applied
 */
func MutateTextLines() {

}

/**
 * Composes two attribute strings (see below) into one.
 * @param att1 {string} first attribute string
 * @param att2 {string} second attribue string
 * @param resultIsMutaton {boolean}
 * @param pool {AttribPool} attribute pool
 */
func ComposeAttributes(att1, att2 string, resultIsMutation bool, pool AttributePool) string {
	// att1 and att2 are strings like "*3*f*1c", asMutation is a boolean.
	// Sometimes attribute (key,value) pairs are treated as attribute presence
	// information, while other times they are treated as operations that
	// mutate a set of attributes, and this affects whether an empty value
	// is a deletion or a change.
	// Examples, of the form (att1Items, att2Items, resultIsMutation) -> result
	// ([], [(bold, )], true) -> [(bold, )]
	// ([], [(bold, )], false) -> []
	// ([], [(bold, true)], true) -> [(bold, true)]
	// ([], [(bold, true)], false) -> [(bold, true)]
	// ([(bold, true)], [(bold, )], true) -> [(bold, )]
	// ([(bold, true)], [(bold, )], false) -> []
	// pool can be null if att2 has no attributes.
	if len(att1) <= 0 && resultIsMutation {
		// In the case of a mutation (i.e. composing two exportss),
		// an att2 composed with an empy att1 is just att2.  If att1
		// is part of an attribution string, then att2 may remove
		// attributes that are already gone, so don't do this optimization.
		return att2
	}
	if len(att2) <= 0 {
		return att1
	}
	//var atts []int
	return ""
}

/**
 * Function used as parameter for applyZip to apply a Changeset to an
 * attribute
 */
func _slicerZipperFunc(attOp, csOp, opOut *Operator, pool AttributePool) {
	// attOp is the op from the sequence that is being operated on, either an
	// attribution string or the earlier of two exportss being composed.
	// pool can be null if definitely not needed.
	//print(csOp.toSource()+" "+attOp.toSource()+" "+opOut.toSource());
	q.Q(attOp, csOp, opOut, pool)
	q.Q("-------")
	if attOp.OpCode == "-" {
		copier.Copy(opOut, attOp)
		attOp.OpCode = ""
	} else if len(attOp.OpCode) <= 0 {
		copier.Copy(opOut, csOp)
		csOp.OpCode = ""
	} else {
		switch csOp.OpCode {
		case "-":
			if csOp.Chars <= attOp.Chars {
				// delete or delete part
				if attOp.OpCode == "=" {
					opOut.OpCode = "-"
					opOut.Chars = csOp.Chars
					opOut.Lines = csOp.Lines
					opOut.Attribs = ""
				}
				attOp.Chars -= csOp.Chars
				attOp.Lines -= csOp.Lines
				csOp.OpCode = ""
				if attOp.Chars <= 0 {
					attOp.OpCode = ""
				}
			} else {
				// delete and keep going
				if attOp.OpCode == "=" {
					opOut.OpCode = "-"
					opOut.Chars = attOp.Chars
					opOut.Lines = attOp.Lines
					opOut.Attribs = ""
				}
				csOp.Chars -= attOp.Chars
				csOp.Lines -= attOp.Lines
				attOp.OpCode = ""
			}
		case "+":
			// insert
			q.Q("cpy bf:", csOp, opOut)
			copier.Copy(&opOut, &csOp)
			q.Q("cpy af:", opOut, csOp)
			csOp.OpCode = ""
		case "=":
			if csOp.Chars <= attOp.Chars {
				// keep or keep part
				opOut.OpCode = attOp.OpCode
				opOut.Chars = csOp.Chars
				opOut.Lines = csOp.Lines
				opOut.Attribs = ComposeAttributes(attOp.Attribs, csOp.Attribs, attOp.OpCode == "=", pool)
				q.Q("=1", opOut.Attribs)
				csOp.OpCode = ""
				attOp.Chars -= csOp.Chars
				attOp.Lines -= csOp.Lines
				if attOp.Chars <= 0 {
					attOp.OpCode = ""
				}
			} else {
				// keep and keep going
				opOut.OpCode = attOp.OpCode
				opOut.Chars = attOp.Chars
				opOut.Lines = attOp.Lines
				opOut.Attribs = ComposeAttributes(attOp.Attribs, csOp.Attribs, attOp.OpCode == "=", pool)
				q.Q("=2", opOut.Attribs)
				attOp.OpCode = ""
				csOp.Chars -= attOp.Chars
				csOp.Lines -= attOp.Lines
			}
		case "":
			copier.Copy(opOut, attOp)
			attOp.OpCode = ""
		}
	}
	q.Q(attOp, csOp, opOut, pool)
	q.Q("++++++++")
}

/**
 * Applies a Changeset to the Attribs string of a AText.
 * @param cs {string} Changeset
 * @param astr {string} the Attribs string of a AText
 * @param pool {AttribsPool} the attibutes pool
 */
func (chgset *ChangeSet) ApplyToAttribution(cs, str string, pool AttributePool) string {
	chgset.Unpack(cs)
	return applyZip(str, chgset.Ops, 0, 0, func(op1 *Operator, op2 *Operator, opOut *Operator) {
		_slicerZipperFunc(op1, op2, opOut, pool)
	})
}

func MutateAttributionLines() {

}

func JoinAttributionLines() {

}

func SplitAttributionLines() {

}

/**
 * splits text into lines
 * @param {string} text to be splitted
 */
func SplitTextLines() {

}

/**
 * compose two Changesets
 * @param cs1 {Changeset} first Changeset
 * @param cs2 {Changeset} second Changeset
 * @param pool {AtribsPool} Attribs pool
 */
func Compose() {

}

/**
 * returns a function that tests if a string of attributes
 * (e.g. *3*4) contains a given attribute key,value that
 * is already present in the pool.
 * @param attribPair array [key,value] of the attribute
 * @param pool {AttribPool} Attribute pool
 */
func AttributeTester() {

}

/**
 * creates the identity Changeset of length N
 * @param N {int} length of the identity changeset
 */
func Identity() {

}

/**
 * creates a Changeset which works on oldFullText and removes text
 * from spliceStart to spliceStart+numRemoved and inserts newText
 * instead. Also gives possibility to add attributes optNewTextAPairs
 * for the new text.
 * @param oldFullText {string} old text
 * @param spliecStart {int} where splicing starts
 * @param numRemoved {int} number of characters to be removed
 * @param newText {string} string to be inserted
 * @param optNewTextAPairs {string} new pairs to be inserted
 * @param pool {AttribPool} Attribution Pool
 */
func (chgset *ChangeSet) MakeSplice(oldFullText string, spliceStart, numRemoved int,
	newText string, optNewTextAPairs string, pool *AttributePool) string {
	chgset.OldLen = len(oldFullText)
	if spliceStart >= chgset.OldLen {
		spliceStart = chgset.OldLen - 1
	}

	if numRemoved > (chgset.OldLen - spliceStart) {
		numRemoved = chgset.OldLen - spliceStart
	}
	oldText := SubString(oldFullText, spliceStart, spliceStart+numRemoved)
	chgset.NewLen = chgset.OldLen + len(newText) - len(oldText)

	assem := smartOpAssembler{}
	assem.appendOpWithText("=", SubString(oldFullText, 0, spliceStart), "", nil)
	assem.appendOpWithText("-", oldText, "", nil)
	assem.appendOpWithText("+", newText, optNewTextAPairs, nil)
	assem.endDocument()
	chgset.Ops = assem.toString()
	chgset.CharBank = newText
	return chgset.Pack()
}

/**
 * Transforms a changeset into a list of splices in the form
 * [startChar, endChar, newText] meaning replace text from
 * startChar to endChar with newText
 * @param cs Changeset
 */
func ToSplices() {

}

func characterRangeFollow() {

}

/**
 * Iterate over attributes in a changeset and move them from
 * oldPool to newPool
 * @param cs {Changeset} Chageset/attribution string to be iterated over
 * @param oldPool {AttribPool} old attributes pool
 * @param newPool {AttribPool} new attributes pool
 * @return {string} the new Changeset
 */
func moveOpsToNewPool(cs string, oldPool, newPool AttributePool) string {
	// works on exports or attribution string
	dollarPos := strings.Index(cs, "$")
	if dollarPos < 0 {
		dollarPos = len(cs)
	}
	upToDollar := SubString(cs, 0, dollarPos)
	fromDollar := SubStrLen(cs, dollarPos, len(cs)) // FIXME
	// order of Attribs stays the same
	rgx, _ := regexp.Compile(`\*([0-9a-z]+)`)
	a := rgx.FindString(upToDollar)
	oldNum, _ := strconv.ParseInt(a, 36, 32)
	pair := oldPool.GetAttrib(int(oldNum))
	// TODO if(!pair) exports.error('Can\'t copy unknown attrib (reference attrib string to non-existant pool entry). Inconsistent attrib state!');
	newNum := newPool.PutAttrib(pair, false)
	return rgx.ReplaceAllString(upToDollar, strconv.FormatInt(int64(newNum), 36)) + fromDollar
}

/**
 * create an attribution inserting a text
 * @param text {string} text to be inserted
 */
func makeAttribution(text string) string {
	assem := smartOpAssembler{}
	assem.appendOpWithText("+", text, "", nil)
	return assem.toString()
}

/**
 * Iterates over attributes in exports, attribution string, or Attribs property of an op
 * and runs function func on them
 * @param cs {Changeset} changeset
 * @param func {function} function to be called
 */
func eachAttribNumber() {

}

/**
 * Filter attributes which should remain in a Changeset
 * callable on a exports, attribution string, or Attribs property of an op,
 * though it may easily create adjacent ops that can be merged.
 * @param cs {Changeset} changeset to be filtered
 * @param filter {function} fnc which returns true if an
 *        attribute X (int) should be kept in the Changeset
 */
func filterAttribNumbers() {

}

/**
 * does exactly the same as exports.filterAttribNumbers
 */
func mapAttribNumbers() {

}

/**
 * Create a Changeset going from Identity to a certain state
 * @params text {string} text of the final change
 * @Attribs Attribs {string} optional, operations which insert
 *    the text and also puts the right attributes
 */
func makeAText(text, attribs string) AText {
	aTxt := AText{
		Text:    text,
		Attribs: attribs,
	}
	if len(attribs) <= 0 {
		aTxt.Attribs = makeAttribution(text)
	}
	return aTxt
}

type AText struct {
	Text    string `json:"text"`
	Attribs string `json:"Attribs"`
}

/**
 * Apply a Changeset to a AText
 * @param cs {Changeset} Changeset to be applied
 * @param atext {AText}
 * @param pool {AttribPool} Attribute Pool to add to
 */
func (chgset *ChangeSet) ApplyToAText(cs string, aText AText, pool AttributePool) AText {
	text, _ := chgset.ApplyToText(cs, aText.Text)
	return AText{
		Text:    text,
		Attribs: chgset.ApplyToAttribution(cs, aText.Attribs, pool),
	}
}

/**
 * Clones a AText structure
 * @param atext {AText}
 */
func cloneAText() {

}

/**
 * Copies a AText structure from atext1 to atext2
 * @param atext {AText}
 */
func copyAText() {

}

/**
 * Append the set of operations from atext to an assembler
 * @param atext {AText}
 * @param assem Assembler like smartOpAssembler
 */
func appendATextToAssembler() {

}

/**
 * Creates a clone of a Changeset and it's APool
 * @param cs {Changeset}
 * @param pool {AtributePool}
 */
func PrepareForWire(cs string, pool AttributePool) (translated string, newPool AttributePool) {
	newPool = NewAttributePool()
	translated = moveOpsToNewPool(cs, pool, newPool)
	return
}

/**
 * Checks if a changeset s the identity changeset
 */
func isIdentity() {

}

/**
 * returns all the values of attributes with a certain key
 * in an Op Attribs string
 * @param Attribs {string} Attribute string of a Op
 * @param key {string} string to be seached for
 * @param pool {AttribPool} attribute pool
 */
func opAttributeValue() {

}

/**
 * returns all the values of attributes with a certain key
 * in an Attribs string
 * @param Attribs {string} Attribute string
 * @param key {string} string to be seached for
 * @param pool {AttribPool} attribute pool
 */
func attribsAttributeValue() {

}

/**
 * Creates a Changeset builder for a string with initial
 * length oldLen. Allows to add/remove parts of it
 * @param oldLen {int} Old length
 */
func NewBuilder(oldLen int) Builder {
	return Builder{
		oldLen:   oldLen,
		assem:    smartOpAssembler{},
		o:        Operator{},
		charBank: stringAssembler{},
	}
}

type Builder struct {
	oldLen   int
	assem    smartOpAssembler
	o        Operator
	charBank stringAssembler
}

// Attribs are [[key1,value1],[key2,value2],...] or '*0*1...' (no pool needed in latter case)
func (b *Builder) Keep(N, L int, attribs string, pool AttributePool) {
	b.o.OpCode = "="
	//  FIXME     o.Attribs = (Attribs && exports.makeAttribsString('=', Attribs, pool)) || '';
	b.o.Attribs = attribs
	b.o.Chars = N
	b.o.Lines = L
	b.assem.append(b.o)
}

func (b *Builder) KeepText(text, attribs string, pool AttributePool) {
	b.assem.appendOpWithText("=", text, attribs, &pool)
}

func (b *Builder) Insert(text, attribs string, pool AttributePool) {
	b.assem.appendOpWithText("+", text, attribs, &pool)
	b.charBank.append(text)
}

func (b *Builder) Remove(N, L int) {
	b.o.OpCode = "-"
	b.o.Attribs = ""
	b.o.Chars = N
	b.o.Lines = L
	b.assem.append(b.o)
}

func (b *Builder) ToString() string {
	b.assem.endDocument()
	newLen := b.oldLen + b.assem.GetLengthChange()
	chgset := ChangeSet{}
	chgset.OldLen = b.oldLen
	chgset.NewLen = newLen
	chgset.Ops = b.assem.toString()
	chgset.CharBank = b.charBank.toString()
	return chgset.Pack()
}

func makeAttribsString(opCode, attribs string, pool *AttributePool) string {
	// makeAttribsString(opcode, '*3') or makeAttribsString(opcode, [['foo','bar']], myPool) work
	if len(opCode) <= 0 {
		return ""
	}
	// FIXME
	return ""

}

// like "substring" but on a single-line attribution string
func subattribution() {

}

func inverse() {
	// lines and alines are what the exports is meant to apply to.
	// They may be arrays or objects with .get(i) and .length methods.
	// They include final newlines on lines.
}

// %CLIENT FILE ENDS HERE%
func Follow() {

}

func followAttributes() {
	// The merge of two sets of attribute changes to the same text
	// takes the lexically-earlier value if there are two values
	// for the same key.  Otherwise, all key/value changes from
	// both attribute sets are taken.  This operation is the "follow",
	// so a set of changes is produced that can be applied to att1
	// to produce the merged set.
}

func composeWithDeletions() {

}

// This function is 95% like _slicerZipperFunc, we just changed two lines to ensure it merges the Attribs of deletions properly.
// This is necassary for correct paddiff. But to ensure these changes doesn't affect anything else, we've created a seperate function only used for paddiffs
func _slicerZipperFuncWithDeletions() {
	// attOp is the op from the sequence that is being operated on, either an
	// attribution string or the earlier of two exportss being composed.
	// pool can be null if definitely not needed.
	//print(csOp.toSource()+" "+attOp.toSource()+" "+opOut.toSource());
}
