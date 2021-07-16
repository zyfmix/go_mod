package tools

import "strings"

func UnicodeIndex(str, substr string) (startIndex, length int) {
	// 子串在字符串的字节位置
	startIndex = strings.Index(str, substr)
	if startIndex >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:startIndex]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		startIndex = len(rs)
		subStrUnicode := []rune(substr)
		length = len(subStrUnicode)
	}

	return startIndex, length
}

func SubUnicodeString(str string, begin, length int) (substr string) {
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
