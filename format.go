package lineprinter

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// counts Thai combining characters that do not consume horizontal space
func countNonInlineThaiCharacters(str string) int {
	combiningChars := map[rune]bool{
		'\u0E31': true,
		'\u0E34': true,
		'\u0E35': true,
		'\u0E36': true,
		'\u0E37': true,
		'\u0E38': true,
		'\u0E39': true,
		'\u0E47': true,
		'\u0E48': true,
		'\u0E49': true,
		'\u0E4A': true,
		'\u0E4B': true,
		'\u0E4C': true,
		'\u0E4D': true,
		'\u0E3A': true,
		'\u0E4E': true,
	}

	count := 0
	for _, r := range str {
		if combiningChars[r] {
			count++
		}
	}

	return count
}

// pads string to the right based on rune length
func padRight(str string, length int) string {
	runeLen := utf8.RuneCountInString(str)
	if runeLen >= length {
		return str
	}

	return str + strings.Repeat(" ", length-runeLen)
}

// calculates required buffer length
func calculateMaxLength(positions []int, args []string) int {
	lastPos := positions[len(positions)-1]
	lastArg := args[len(args)-1]

	return lastPos - 1 + utf8.RuneCountInString(lastArg)
}

// initializes rune buffer with spaces
func createBuffer(length int) []rune {
	runes := make([]rune, length)
	for i := range runes {
		runes[i] = ' '
	}

	return runes
}

// inserts string into rune buffer at given position
func insertIntoBuffer(runes []rune, pos int, arg string) []rune {
	argRunes := []rune(arg)
	argLength := len(argRunes)
	remainingSpace := len(runes) - (pos + argLength)
	paddedStr := padRight(arg, remainingSpace)

	for j, r := range []rune(paddedStr) {
		if pos+j < len(runes) {
			runes[pos+j] = r
		}
	}

	return runes
}

// Format prints fixed-position text with Thai alignment handling
func Format(buffer io.Writer, positions []int, args ...string) error {
	if len(positions) != len(args) {
		return fmt.Errorf("positions and args length mismatch")
	}

	maxLength := calculateMaxLength(positions, args)
	totalPadding := 0

	padLengths := make([]int, len(args))
	for i, arg := range args {
		padLengths[i] = countNonInlineThaiCharacters(arg)
		totalPadding += padLengths[i]
	}

	runes := createBuffer(maxLength + totalPadding)

	totalExtraPadding := 0
	for i, arg := range args {
		pos := positions[i] - 1
		runes = insertIntoBuffer(runes, pos+totalExtraPadding, arg)
		totalExtraPadding += padLengths[i]
	}

	_, err := buffer.Write([]byte(string(runes) + "\n"))
	return err
}

// FormatLine writes separator
func FormatLine(buffer io.Writer) error {
	_, err := buffer.Write([]byte("-----------------------------------------------------------------------------------------------------\n"))

	return err
}

// FormatNewLine writes N line feeds
func FormatNewLine(buffer io.Writer, total int) error {
	_, err := buffer.Write([]byte(strings.Repeat("\n", total)))

	return err
}

// FormatNewPage writes form feed
func FormatNewPage(buffer io.Writer) error {
	_, err := buffer.Write([]byte{12})

	return err
}
