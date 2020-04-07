package changeset

/**
 * This method is called whenever there is an error in the sync process
 * @param msg {string} Just some message
 */

func Error() {

}

/**
 * This method is used for assertions with Messages
 * if assert fails, the error function is called.
 * @param b {boolean} assertion condition
 * @param msgParts {string} error to be passed if it fails
 */
func _assert() {

}

/**
 * Parses a number from string base 36
 * @param str {string} string of the number in base 36
 * @returns {int} number
 */
func parseNum() {

}

/**
 * Writes a number in base 36 and puts it in a string
 * @param num {int} number
 * @returns {string} string
 */
func numToString() {

}

/**
 * Converts stuff before $ to base 10
 * @obsolete not really used anywhere??
 * @param cs {string} the string
 * @return integer
 */
func toBaseTen() {

}

func SubString(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

func SubStrLen(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}
