package lineprinter

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

// Windows-874 encoder for Thai text
var tis620Encoder = charmap.Windows874.NewEncoder()

// Thai character levels
const (
	LevelUpper  = 1
	LevelBase   = 2
	LevelBottom = 3
)

// Pair of runes for mapping
type runePair struct {
	r1 rune
	r2 rune
}

// Combine Character map (vowel + tone mark) for Thai line printer encoding
var combinedCharacterMap = map[runePair][]byte{
	{0x0E34, 0x0E4C}: {135}, // ิ + ์  -> สระอิ + การันต์
	{0x0E34, 0x0E48}: {136}, // ิ + ่  -> สระอิ + ไม้เอก
	{0x0E34, 0x0E49}: {137}, // ิ + ้  -> สระอิ + ไม้โท
	{0x0E34, 0x0E4A}: {138}, // ิ + ๊  -> สระอิ + ไม้ตรี
	{0x0E34, 0x0E4B}: {139}, // ิ + ๋  -> สระอิ + ไม้จัตวา

	{0x0E35, 0x0E48}: {140}, // ี + ่  -> สระอี + ไม้เอก
	{0x0E35, 0x0E49}: {141}, // ี + ้  -> สระอี + ไม้โท
	{0x0E35, 0x0E4A}: {142}, // ี + ๊  -> สระอี + ไม้ตรี
	{0x0E35, 0x0E4B}: {143}, // ี + ๋  -> สระอี + ไม้จัตวา

	{0x0E36, 0x0E48}: {144}, // ึ + ่  -> สระอึ + ไม้เอก
	{0x0E36, 0x0E49}: {145}, // ึ + ้  -> สระอึ + ไม้โท
	{0x0E36, 0x0E4A}: {146}, // ึ + ๊  -> สระอึ + ไม้ตรี
	{0x0E36, 0x0E4B}: {147}, // ึ + ๋  -> สระอึ + ไม้จัตวา

	{0x0E37, 0x0E48}: {148}, // ื + ่  -> สระอื + ไม้เอก
	{0x0E37, 0x0E49}: {149}, // ื + ้  -> สระอื + ไม้โท
	{0x0E37, 0x0E4A}: {150}, // ื + ๊  -> สระอื + ไม้ตรี
	{0x0E37, 0x0E4B}: {151}, // ื + ๋  -> สระอื + ไม้จัตวา

	{0x0E31, 0x0E48}: {152}, // ั + ่  -> ไม้หันอากาศ + ไม้เอก
	{0x0E31, 0x0E49}: {153}, // ั + ้  -> ไม้หันอากาศ + ไม้โท
	{0x0E31, 0x0E4A}: {154}, // ั + ๊  -> ไม้หันอากาศ + ไม้ตรี
	{0x0E31, 0x0E4B}: {155}, // ั + ๋  -> ไม้หันอากาศ + ไม้จัตวา

	{0x0E4D, 0x0E48}: {156}, // ํ + ่  -> นิคหิต + ไม้เอก
	{0x0E4D, 0x0E49}: {157}, // ํ + ้  -> นิคหิต + ไม้โท
	{0x0E4D, 0x0E4A}: {158}, // ํ + ๊  -> นิคหิต + ไม้ตรี
	{0x0E4D, 0x0E4B}: {159}, // ํ + ๋  -> นิคหิต + ไม้จัตวา
}

// Classify Thai rune into levels
func classifyLevel(r rune) int {
	if r >= 0x0E01 && r <= 0x0E4E {
		switch r {
		case 0x0E38, 0x0E39, 0x0E3A:
			return LevelBottom

		case 0x0E48, 0x0E49, 0x0E4A, 0x0E4B,
			0x0E4C, 0x0E4D, 0x0E47,
			0x0E31, 0x0E34, 0x0E35,
			0x0E36, 0x0E37:
			return LevelUpper

		default:
			return LevelBase
		}
	}

	return LevelBase
}

// Convert rune to Windows-874 byte
func encodeRuneToTIS620Byte(r rune) []byte {
	src := make([]byte, utf8.RuneLen(r))
	utf8.EncodeRune(src, r)

	dst := make([]byte, 1)
	nDst, _, err := tis620Encoder.Transform(dst, src, true)
	if err != nil || nDst == 0 {
		return []byte{' '}
	}

	return dst[:nDst]
}

// Decode Windows-874 byte to rune
func decodeWindows874ByteToRune(b byte) rune {
	decoder := charmap.Windows874.NewDecoder()
	src := []byte{b}
	dst := make([]byte, 4)
	nDst, _, err := decoder.Transform(dst, src, true)
	if err != nil || nDst == 0 {

		return ' '
	}
	r, _ := utf8.DecodeRune(dst[:nDst])

	return r
}

// Cell holds bytes for each level
type cell struct {
	lvUpper  []byte
	lvBase   []byte
	lvBottom []byte
}

// Apply combined character for Thai vowels + tone marks
func applyCombinedCharacter(cells []cell) []cell {
	for i := range cells {
		c := &cells[i]
		if len(c.lvUpper) >= 2 {
			r1 := decodeWindows874ByteToRune(c.lvUpper[0])
			r2 := decodeWindows874ByteToRune(c.lvUpper[1])

			if proprietaryByte, found := combinedCharacterMap[runePair{r1, r2}]; found {
				newUpper := append([]byte{}, proprietaryByte...)
				newUpper = append(newUpper, c.lvUpper[2:]...)
				c.lvUpper = newUpper

				continue
			}

			if proprietaryByte, found := combinedCharacterMap[runePair{r2, r1}]; found {
				newUpper := append([]byte{}, proprietaryByte...)
				newUpper = append(newUpper, c.lvUpper[2:]...)
				c.lvUpper = newUpper
			}
		}
	}

	return cells
}

// Decomposes SaraAm
func decomposeSaraAm(text string) []rune {
	runesToProcess := make([]rune, 0, len(text)*2)
	for _, r := range text {
		if r == 0x0E33 { // ำ
			runesToProcess = append(runesToProcess, 0x0E4D)
			runesToProcess = append(runesToProcess, 0x0E32)
		} else {
			runesToProcess = append(runesToProcess, r)
		}
	}

	return runesToProcess
}

// Add a new base cell
func addBaseCell(cells *[]cell, base []byte) int {
	*cells = append(*cells, cell{lvBase: base})

	return len(*cells) - 1
}

// Handle upper and bottom level runes for the current cell
func handleUpperAndBottomLevels(cells *[]cell, cur *int, r rune, level int) {
	if *cur == -1 {
		*cur = addBaseCell(cells, encodeRuneToTIS620Byte(' '))
	}

	switch level {
	case LevelUpper:
		(*cells)[*cur].lvUpper = append((*cells)[*cur].lvUpper, encodeRuneToTIS620Byte(r)...)
	case LevelBottom:
		(*cells)[*cur].lvBottom = append((*cells)[*cur].lvBottom, encodeRuneToTIS620Byte(r)...)
	}
}

// Convert text into structured cells
func processTextToCells(text string) []cell {
	var cells []cell
	cur := -1

	processedRune := decomposeSaraAm(text)

	for _, r := range processedRune {
		level := classifyLevel(r)

		if level == LevelBase {
			cur = addBaseCell(&cells, encodeRuneToTIS620Byte(r))
			continue
		}

		handleUpperAndBottomLevels(&cells, &cur, r, level)
	}

	return cells
}

// Build byte lines for printing
func buildLines(cells []cell) ([]byte, []byte, []byte) {
	var b1, b2, b3 bytes.Buffer
	space := []byte{' '}

	for _, c := range cells {
		if len(c.lvUpper) > 0 {
			b1.Write(c.lvUpper)
		} else {
			b1.Write(space)
		}
		if len(c.lvBase) > 0 {
			b2.Write(c.lvBase)
		} else {
			b2.Write(space)
		}
		if len(c.lvBottom) > 0 {
			b3.Write(c.lvBottom)
		} else {
			b3.Write(space)
		}
	}

	return b1.Bytes(), b2.Bytes(), b3.Bytes()
}

// Split Thai text into upper, base, and bottom lines
func splitThaiText(text string) ([]byte, []byte, []byte) {
	cells := processTextToCells(text)
	cells = applyCombinedCharacter(cells)

	return buildLines(cells)
}

// Convert Thai text to 3-level Windows-874 byte format
func levelingThaiText(text string) []byte {
	lvUpper, lvBase, lvBottom := splitThaiText(text)

	var output bytes.Buffer

	output.Write(lvUpper)
	output.WriteByte('\n')
	output.Write(lvBase)
	output.WriteByte('\n')
	output.Write(lvBottom)
	output.WriteByte('\n')

	return output.Bytes()
}

func ToLP3(content string) []byte {
	lines := strings.Split(content, "\n")
	var formatted [][]byte

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			formatted = append(formatted, []byte{'\n'})
		} else {
			formatted = append(formatted, levelingThaiText(line))
		}
	}

	return bytes.Join(formatted, []byte{})
}
