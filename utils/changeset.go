package utils

// @description Unpacks a string encoded Changeset into a proper Changeset object
// @param cs {string} String encoded Changeset
// @return {Changeset} a Changeset class
func Unpack() {

}

/**
 * Packs Changeset object into a string
 * @params oldLen {int} Old length of the Changeset
 * @params newLen {int] New length of the Changeset
 * @params opsStr {string} String encoding of the changes to be made
 * @params bank {string} Charbank of the Changeset
 * @returns {Changeset} a Changeset class
 */
func Pack() {

}

/**
 * Applies a Changeset to a string
 * @params cs {string} String encoded Changeset
 * @params str {string} String to which a Changeset should be applied
 */
func ApplyToText() {

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
func ComposeAttributes() {
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
}

/**
 * Function used as parameter for applyZip to apply a Changeset to an
 * attribute
 */

func _slicerZipperFunc() {
	// attOp is the op from the sequence that is being operated on, either an
	// attribution string or the earlier of two exportss being composed.
	// pool can be null if definitely not needed.
	//print(csOp.toSource()+" "+attOp.toSource()+" "+opOut.toSource());

}

/**
 * Applies a Changeset to the attribs string of a AText.
 * @param cs {string} Changeset
 * @param astr {string} the attribs string of a AText
 * @param pool {AttribsPool} the attibutes pool
 */
func ApplyToAttribution() {

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
func MakeSplice() {

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
func moveOpsToNewPool() {

}

/**
 * create an attribution inserting a text
 * @param text {string} text to be inserted
 */
func makeAttribution() {

}

/**
 * Iterates over attributes in exports, attribution string, or attribs property of an op
 * and runs function func on them
 * @param cs {Changeset} changeset
 * @param func {function} function to be called
 */
func eachAttribNumber() {

}

/**
 * Filter attributes which should remain in a Changeset
 * callable on a exports, attribution string, or attribs property of an op,
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
 * @attribs attribs {string} optional, operations which insert
 *    the text and also puts the right attributes
 */
func makeAText() {

}

/**
 * Apply a Changeset to a AText
 * @param cs {Changeset} Changeset to be applied
 * @param atext {AText}
 * @param pool {AttribPool} Attribute Pool to add to
 */
func applyToAText() {

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
func prepareForWire() {

}

/**
 * Checks if a changeset s the identity changeset
 */
func isIdentity() {

}

/**
 * returns all the values of attributes with a certain key
 * in an Op attribs string
 * @param attribs {string} Attribute string of a Op
 * @param key {string} string to be seached for
 * @param pool {AttribPool} attribute pool
 */
func opAttributeValue(){

}

/**
 * returns all the values of attributes with a certain key
 * in an attribs string
 * @param attribs {string} Attribute string
 * @param key {string} string to be seached for
 * @param pool {AttribPool} attribute pool
 */
func attribsAttributeValue (){

}


/**
 * Creates a Changeset builder for a string with initial
 * length oldLen. Allows to add/remove parts of it
 * @param oldLen {int} Old length
 */
func Builder(){

}

func makeAttribsString(){

}


// like "substring" but on a single-line attribution string
func subattribution (){

}

func inverse(){
	// lines and alines are what the exports is meant to apply to.
	// They may be arrays or objects with .get(i) and .length methods.
	// They include final newlines on lines.
}

// %CLIENT FILE ENDS HERE%
func Follow(){

}

func followAttributes(){
	// The merge of two sets of attribute changes to the same text
	// takes the lexically-earlier value if there are two values
	// for the same key.  Otherwise, all key/value changes from
	// both attribute sets are taken.  This operation is the "follow",
	// so a set of changes is produced that can be applied to att1
	// to produce the merged set.
}

func composeWithDeletions(){

}

// This function is 95% like _slicerZipperFunc, we just changed two lines to ensure it merges the attribs of deletions properly.
// This is necassary for correct paddiff. But to ensure these changes doesn't affect anything else, we've created a seperate function only used for paddiffs
func _slicerZipperFuncWithDeletions(){
	// attOp is the op from the sequence that is being operated on, either an
	// attribution string or the earlier of two exportss being composed.
	// pool can be null if definitely not needed.
	//print(csOp.toSource()+" "+attOp.toSource()+" "+opOut.toSource());
}