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
