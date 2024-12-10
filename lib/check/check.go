package check

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// upperThreshold: より大きい
// lowerThreshold: より小さい
// 場合にエラーを出す
func WarnLineLen(r io.Reader, w io.Writer, upperThreshold int, lowerThreshold int, fileName string) {
	fmt.Fprintf(w, "%s\n", fileName)

	scanner := bufio.NewScanner(r)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		if isKanjiNumeralsOnly(line) {
			continue
		}
		lineLen := len([]rune(line))
		if lineLen == 0 {
			continue
		}
		if lowerThreshold > lineLen || lineLen > upperThreshold {
			fmt.Fprintf(w, "Line: %d, Length: %d\n", lineNumber, lineLen)
			fmt.Fprintf(w, "  %s\n", line)
		}
		lineNumber++
	}
	fmt.Fprintln(w, strings.Repeat("-", 80))

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(w, "Error reading: %v\n", err)
	}
}

var kanjiNumerals = map[rune]bool{
	'一': true, '二': true, '三': true, '四': true,
	'五': true, '六': true, '七': true, '八': true,
	'九': true, '十': true, '百': true, '〇': true,
	'零': true, ' ': true, '編': true,
}

func isKanjiNumeralsOnly(input string) bool {
	for _, r := range input {
		if !kanjiNumerals[r] {
			return false
		}
	}

	return len(input) > 0
}

func WarnNotes(r io.Reader, w io.Writer, fileName string) {
	fmt.Fprintf(w, "%s\n", fileName)

	scanner := bufio.NewScanner(r)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		notes := []string{
			"＃",
			"※",
		}
		for _, noteMark := range notes {
			if strings.Contains(line, noteMark) {
				fmt.Fprintf(w, "Line: %d\n", lineNumber)
				fmt.Fprintf(w, "  %s\n", line)
			}
		}

		lineNumber++
	}
	fmt.Fprintln(w, strings.Repeat("-", 80))

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(w, "Error reading: %v\n", err)
	}
}
