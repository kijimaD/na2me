package check

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func WarnLongLine(r io.Reader, w io.Writer, threshold int, fileName string) {
	fmt.Fprintf(w, "%s\n", fileName)

	scanner := bufio.NewScanner(r)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		lineLen := len([]rune(line))
		if lineLen > threshold {
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
