/*
 *    Copyright 2021 chenquan
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package xstring

import (
	"errors"
	"fmt"
	"github.com/chenquan/go-pkg/internal/hack"
	"github.com/chenquan/go-pkg/xmath"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	empty         = "" // 空字符串
	indexNotFound = -1
)

var (
	ErrDecodeChar = errors.New("error occurred on char decoding")
)

// Char is rune Alias.
type Char = rune

// PadLeftChar left pad a string with a specified character in a larger string (specified size).
// if the size is less than the param string, the param string is returned.
// NOTE: size is unicode size.
func PadLeftChar(s string, size int, ch Char) string {
	return padCharLeftOrRight(s, size, ch, true)
}

// PadLeftSpace left pad a string with space character(' ') in a larger string(specified size).
// if the size is less than the param string, the param string is returned.
// NOTE: size is unicode size.
func PadLeftSpace(s string, size int) string {
	return PadLeftChar(s, size, ' ')
}

// PadRightChar right pad a string with a specified character in a larger string(specified size).
// if the size is less than the param string, the param string is returned.
// NOTE: size is unicode size.
func PadRightChar(s string, size int, ch Char) string {
	return padCharLeftOrRight(s, size, ch, false)
}

// PadRightSpace right pad a string with space character(' ') in a large string(specified size).
// if the size is less than the param string, the param string is returned.
// NOTE: size is unicode size.
func PadRightSpace(s string, size int) string {
	return PadRightChar(s, size, ' ')
}

// PadCenterChar center pad a string with a specified character in a larger string(specified size).
// if the size is less than the param string, the param string is returned.
// NOTE: size is unicode size.
func PadCenterChar(s string, size int, ch Char) string {
	if size <= 0 {
		return s
	}

	length := Len(s)
	pads := size - length
	if pads <= 0 {
		return s
	}

	// pad left
	leftPads := pads / 2
	if leftPads > 0 {
		s = padRawLeftChar(s, ch, leftPads)
	}
	// pad right
	rightPads := size - leftPads - length
	if rightPads > 0 {
		s = padRawRightChar(s, ch, rightPads)
	}

	return s
}

// PadCenterSpace center pad a string with space character(' ') in a larger string(specified size).
// if the size is less than the param string, the param string is returned.
// NOTE: size is unicode size.
func PadCenterSpace(s string, size int) string {
	return PadCenterChar(s, size, ' ')
}

func padCharLeftOrRight(s string, size int, ch Char, isLeft bool) string {
	if size <= 0 {
		return s
	}

	pads := size - Len(s)
	if pads <= 0 {
		return s
	}

	if isLeft {
		return padRawLeftChar(s, ch, pads)
	}

	return padRawRightChar(s, ch, pads)
}

func padRawLeftChar(s string, ch Char, padSize int) string {
	return RepeatChar(ch, padSize) + s
}

func padRawRightChar(s string, ch Char, padSize int) string {
	return s + RepeatChar(ch, padSize)
}

// RepeatChar returns padding using the specified delimiter repeated to a given length.
func RepeatChar(ch Char, repeat int) string {
	if repeat <= 0 {
		return ""
	}

	sb := strings.Builder{}
	sb.Grow(repeat)
	for i := 0; i < repeat; i++ {
		sb.WriteRune(ch)
	}

	return sb.String()
}

// RemoveChar removes all occurrences of a specified character from the string.
func RemoveChar(s string, rmVal Char) string {
	if s == "" {
		return s
	}
	sb := strings.Builder{}
	sb.Grow(len(s) / 2)

	for _, v := range s {
		if v != rmVal {
			sb.WriteRune(v)
		}
	}

	return sb.String()
}

// RemoveString removes all occurrences of a substring from the string.
func RemoveString(s, rmStr string) string {
	if s == "" || rmStr == "" {
		return s
	}

	return strings.ReplaceAll(s, rmStr, "")
}

// Rotate rotates(circular shift) a string of shift characters.
func Rotate(s string, shift int) string {
	if shift == 0 {
		return s
	}

	sLen := len(s)
	if sLen == 0 {
		return s
	}

	shiftMod := shift % sLen
	if shiftMod == 0 {
		return s
	}

	offset := -(shiftMod)
	sb := strings.Builder{}
	sb.Grow(sLen)
	_, _ = sb.WriteString(Left(s, offset))
	_, _ = sb.WriteString(Right(s, offset))

	return sb.String()
}

// Sub returns substring from specified string avoiding panics with index start and end.
// start, end are based on unicode(utf8) count.
func Sub(s string, start, end int) string {
	return sub(s, start, end)
}

// Left returns substring from specified string avoiding panics with start.
// start, end are based on unicode(utf8) count.
func Left(s string, end int) string {
	return sub(s, 0, end)
}

// Right returns substring from specified string avoiding panics with end.
// start, end are based on unicode(utf8) count.
func Right(s string, start int) string {
	return sub(s, start, math.MaxInt)
}

func sub(s string, start, end int) string {
	if s == "" {
		return ""
	}

	unicodeLen := Len(s)
	// end
	if end < 0 {
		end += unicodeLen
	}
	if end > unicodeLen {
		end = unicodeLen
	}
	// start
	if start < 0 {
		start += unicodeLen
	}
	if start > end {
		return ""
	}

	// start <= end
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start == 0 && end == unicodeLen {
		return s
	}

	sb := strings.Builder{}
	sb.Grow(end - start)
	runeIndex := 0
	for _, v := range s {
		if runeIndex >= end {
			break
		}
		if runeIndex >= start {
			sb.WriteRune(v)
		}
		runeIndex++
	}

	return sb.String()
}

// MustReverse reverses a string, panics when error happens.
func MustReverse(s string) string {
	result, err := Reverse(s)

	if err != nil {
		panic(err)
	}

	return result
}

// Reverse reverses a string with error status returned.
func Reverse(s string) (string, error) {
	if s == "" {
		return s, nil
	}

	src := hack.StringToBytes(s)
	dst := make([]byte, len(s))
	srcIndex := len(s)
	dstIndex := 0
	for srcIndex > 0 {
		r, n := utf8.DecodeLastRune(src[:srcIndex])

		if r == utf8.RuneError {
			return hack.BytesToString(dst), ErrDecodeChar
		}

		utf8.EncodeRune(dst[dstIndex:], r)
		srcIndex -= n
		dstIndex += n
	}

	return hack.BytesToString(dst), nil
}

// Shuffle shuffles runes in a string and returns.
func Shuffle(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	index := 0
	for i := len(runes) - 1; i > 0; i-- {
		index = rand.Intn(i + 1)
		if i != index {
			runes[i], runes[index] = runes[index], runes[i]
		}
	}

	return string(runes)
}

// ContainsAnySubstrings returns whether s contains any of substring in slice.
func ContainsAnySubstrings(s string, subs []string) bool {
	if len(subs) == 0 {
		return false
	}

	for _, v := range subs {
		if strings.Contains(s, v) {
			return true
		}
	}

	return false
}

// IsAlpha checks if the string contains only unicode letters.
func IsAlpha(s string) bool {
	if s == empty {
		return false
	}

	for _, v := range s {
		if !unicode.IsLetter(v) {
			return false
		}
	}

	return true
}

// IsAlphanumeric checks if the string contains only Unicode letters or digits.
func IsAlphanumeric(s string) bool {
	if s == empty {
		return false
	}

	for _, v := range s {
		if !isAlphanumeric(v) {
			return false
		}
	}

	return true
}

func isAlphanumeric(v Char) bool {
	return unicode.IsDigit(v) || unicode.IsLetter(v)
}

// IsNumeric checks if the string contains only digits. A decimal point is not a digit and returns false.
func IsNumeric(s string) bool {
	if s == empty {
		return false
	}

	for _, v := range s {
		if !unicode.IsDigit(v) {
			return false
		}
	}

	return true
}

// IsEmpty returns ture if s is empty.
func IsEmpty(s string) bool {
	return s == empty
}

// IsNotEmpty returns ture if s isn't empty.
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// IsAnyEmpty returns ture if strings exist empty of string.
func IsAnyEmpty(strings ...string) bool {
	if len(strings) == 0 {
		return true
	}
	for _, s := range strings {
		if IsEmpty(s) {
			return true
		}
	}

	return false
}

// IsNoneEmpty returns false if strings exist empty of string.
func IsNoneEmpty(strings ...string) bool {
	return !IsAnyEmpty(strings...)
}

// IsBlank returns ture if the string is empty or has a length of 0 or consists of whitespace.
func IsBlank(s string) bool {
	if s == empty {
		return true
	}

	for _, c := range s {
		if !unicode.IsSpace(c) {
			return false
		}
	}

	return true
}

//IsNotBlank returns false if the string is empty or has a length of 0 or consists of whitespace.
func IsNotBlank(s string) bool {
	return !IsBlank(s)
}

// IsAnyBlank returns true if there is a blank string in strings.
func IsAnyBlank(strings ...string) bool {
	if len(strings) == 0 {
		return true
	}

	for _, s := range strings {
		if IsBlank(s) {
			return true
		}
	}

	return false

}

// IsNoneBlank returns true if there is no blank string in strings.
func IsNoneBlank(strings ...string) bool {
	return !IsAnyBlank(strings...)
}

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// Trim returns a slice of the string s with all leading and
// trailing Unicode code points contained in cutSet removed.
func Trim(s string, cutSet string) string {
	if IsEmpty(s) {
		return s
	}

	return strings.Trim(s, cutSet)
}

// TrimLeft returns a slice of the string s with all leading
// Unicode code points contained in cutset removed.
//
// To remove a prefix, use TrimPrefix instead.
func TrimLeft(str string, stripStr string) string {
	return strings.TrimLeft(str, stripStr)
}

// TrimRight returns a slice of the string s, with all trailing
// Unicode code points contained in cutset removed.
//
// To remove a suffix, use TrimSuffix instead.
func TrimRight(str string, stripStr string) string {
	return strings.TrimRight(str, stripStr)
}

// Contains reports whether substr is within s.
func Contains(s string, searchChar string) bool {
	if s == empty {
		return false
	}

	return strings.Contains(s, searchChar)
}

// IsNumerical returns ture if a numerical.
func IsNumerical(s string) bool {
	reg, _ := regexp.Compile("^\\d+.?\\d*$")
	return reg.MatchString(s)
}

// IndexOfDifference compares all strings in an array and returns the index at which the
//   string begin to differ.
func IndexOfDifference(strings ...string) int {
	stringsLen := len(strings)
	if len(strings) > 1 {
		allStringsNull := true
		shortestStrLen := 2<<32 - 1
		longestStrLen := 0

		firstDiff := 0
		for ; firstDiff < stringsLen; firstDiff++ {
			allStringsNull = false
			runes := []rune(strings[firstDiff])
			shortestStrLen = xmath.MinInt(len(runes), shortestStrLen)
			longestStrLen = xmath.MaxInt(len(runes), longestStrLen)
		}
		if allStringsNull || longestStrLen == 0 {
			return indexNotFound
		} else if shortestStrLen == 0 {
			return 0
		} else {
			firstDiff = -1

			runes := []rune(strings[0])
			for stringPos := 0; stringPos < shortestStrLen; stringPos++ {
				comparisonChar := runes[stringPos]
				for arrayPos := 1; arrayPos < stringsLen; arrayPos++ {
					if []rune(strings[arrayPos])[stringPos] != comparisonChar {
						firstDiff = stringPos
						break
					}
				}
				if firstDiff != -1 {
					break
				}
			}
			if firstDiff == -1 && shortestStrLen != longestStrLen {
				return shortestStrLen
			} else {
				return firstDiff
			}
		}

	} else {
		return indexNotFound
	}
}

// IndexOfDifferenceWithTwoStr Compares two string, and returns the index at which the
// string begin to differ.
func IndexOfDifferenceWithTwoStr(a, b string) int {
	if a == b {
		return indexNotFound
	} else {
		aRunes := []rune(a)
		bRunes := []rune(b)
		aLen := len(aRunes)
		bLen := len(bRunes)
		i := 0
		for ; i < aLen && i < bLen && aRunes[i] == bRunes[i]; i++ {
		}
		if i >= aLen && i >= bLen {
			return indexNotFound
		} else {
			return i
		}
	}
}

// Difference Compares two Strings, and returns the portion where they differ.
// More precisely, return the remainder of the second String,
// starting from where it's different from the first. This means that
// the difference between "abc" and "ab" is the empty String and not "c".
func Difference(a, b string) string {
	i := IndexOfDifferenceWithTwoStr(a, b)
	if i == -1 {
		return empty
	} else {
		runes := []rune(b)
		return string(runes[i:])
	}
}

func CommonPrefix(strings ...string) string {
	if len(strings) != 0 {
		smallestIndexOfDiff := IndexOfDifference(strings...)
		if smallestIndexOfDiff == -1 {
			return strings[0]
		} else {
			if smallestIndexOfDiff == 0 {
				return empty
			} else {
				runes := []rune(strings[0])
				return string(runes[0:smallestIndexOfDiff])
			}
		}
	} else {
		return empty
	}
}

// Index returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func Index(s, substr string) int {
	return strings.Index(s, substr)
}

// IndexAny returns the index of the first instance of any Unicode code point
// from chars in s, or -1 if no Unicode code point from chars is present in s.
func IndexAny(s, chars string) int {
	return strings.IndexAny(s, chars)
}

// ContainsIgnoreCase checks if string contains a search string irrespective of case.
func ContainsIgnoreCase(str, searchStr string) bool {
	length := Len(searchStr)
	max := Len(str) - length

	for i := 0; i <= max; i++ {
		if RegionMatches(str, true, i, searchStr, 0, length) {
			return true
		}
	}

	return false
}

// RegionMatches Green implementation of regionMatches.
func RegionMatches(str string, ignoreCase bool, thisStart int, substr string, start int, length int) bool {
	if ignoreCase {
		str = strings.ToLower(str)
		substr = strings.ToLower(substr)
	}
	// Check the validity of the parameters
	if thisStart < 0 || start < 0 || length < 0 {
		return false
	}
	// The length of the remaining part of the string
	thisRetLen := Len(str) - thisStart
	subRetLen := Len(substr) - start
	if thisRetLen < length || subRetLen < length {
		return false
	}
	strToRunes := []rune(str)
	substrToRunes := []rune(substr)
	for ; length > 0; length-- {
		c1 := strToRunes[thisStart]
		c2 := substrToRunes[start]
		thisStart++
		start++
		if c1 == c2 {
			continue
		}

		if ignoreCase {
			c1 = unicode.ToUpper(c1)
			c2 = unicode.ToUpper(c2)
			if c1 == c2 {
				continue
			}
			if unicode.ToLower(c1) == unicode.ToLower(c2) {
				continue
			}
		}
		return false
	}

	return true
}

// Len returns str length.
func Len(str string) int {
	return utf8.RuneCountInString(str)
}

// DefaultIfBlank returns default String if str is blank.
func DefaultIfBlank(str, defaultStr string) string {
	if IsBlank(str) {
		return defaultStr
	}

	return str
}

// DefaultIfEmpty returns default String if str is empty.
func DefaultIfEmpty(str, defaultStr string) string {
	if IsEmpty(str) {
		return defaultStr
	}

	return str
}

// DeleteWhitespace deletes whitespace.
func DeleteWhitespace(str string) string {
	if str == empty {
		return empty
	}

	strLen := Len(str)
	builder := strings.Builder{}
	builder.Grow(strLen)

	for _, r := range str {
		if !unicode.IsSpace(r) {
			builder.WriteRune(r)
		}
	}

	if builder.Len() == strLen {
		return str
	}

	if builder.Len() == 0 {
		return empty
	}

	return builder.String()
}

// EndsWith returns true if the str
func EndsWith(str, suffix string, ignoreCase bool) bool {
	if str == suffix {
		return true
	}

	strLen := Len(str)
	suffixLen := Len(suffix)
	if suffixLen > strLen {
		return false
	}
	strOffset := strLen - suffixLen

	return RegionMatches(str, ignoreCase, strOffset, suffix, 0, suffixLen)
}

// EndsWithIgnoreCase case-insensitive check if a str ends with a specified suffix.
func EndsWithIgnoreCase(str, suffix string) bool {
	return EndsWith(str, suffix, true)
}

// EndsWithCase case-insensitive check if a str ends with a specified suffix.
func EndsWithCase(str, suffix string) bool {
	return EndsWith(str, suffix, false)
}

// EndsWithAny check if a sequence ends with any of an array of specified strings.
func EndsWithAny(sequence string, searchStrings ...string) bool {
	if IsEmpty(sequence) || len(searchStrings) == 0 {
		return false
	}

	for _, str := range searchStrings {
		if EndsWith(sequence, str, false) {
			return true
		}
	}

	return false
}

// EqualsIgnoreCase returns true if the two strings are equal ignoring case.
func EqualsIgnoreCase(str1, str2 string) bool {
	if str1 == str2 {
		return true
	}

	return RegionMatches(str1, true, 0, str2, 0, Len(str1))
}

// EqualsAny returns ture if str exists in strings.
func EqualsAny(str1 string, strings ...string) bool {
	if len(strings) != 0 {
		for _, str2 := range strings {
			if str1 == str2 {
				return true
			}
		}
	}

	return false
}

// Abbreviate abbreviates a string using ellipses or another given string.
func Abbreviate(str, abbrevMarker string, offset, maxWidth int) (string, error) {
	if IsNotEmpty(str) && abbrevMarker == empty && maxWidth > 0 {
		return Sub(str, 0, maxWidth), nil
	} else if IsAnyEmpty(str, abbrevMarker) {
		// 其中有一个字符串为,则直接返回原字符串
		return str, nil
	}
	abbrevMarkerLen := Len(abbrevMarker)
	// 最小缩减宽度
	minAbbrevWidth := abbrevMarkerLen + 1
	if maxWidth < minAbbrevWidth {
		return empty, fmt.Errorf("minimum abbreviation width is %d", minAbbrevWidth)
	}
	minAbbrevWidthOffset := 2*abbrevMarkerLen + 1
	strLen := Len(str)
	if strLen <= maxWidth {
		return str, nil
	}
	if offset > strLen {
		offset = strLen
	}
	if strLen-offset < maxWidth-abbrevMarkerLen {
		offset = strLen - (maxWidth - abbrevMarkerLen)
	}
	if offset <= abbrevMarkerLen+1 {
		return Sub(str, 0, maxWidth-abbrevMarkerLen) + abbrevMarker, nil
	}
	if maxWidth < minAbbrevWidthOffset {
		return empty, fmt.Errorf("minimum abbreviation width with offset is %d", minAbbrevWidthOffset)
	}
	if offset+maxWidth-abbrevMarkerLen < strLen {
		substr, err := Abbreviate(Sub(str, 0, offset), abbrevMarker, 0, maxWidth-abbrevMarkerLen)
		if err == nil {
			return abbrevMarker + substr, nil
		}
		return empty, err
	}

	return abbrevMarker + Sub(str, 0, strLen-(maxWidth-abbrevMarkerLen)), nil
}
