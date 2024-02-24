package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var div float64 = 25.0

func main() {
	input := "./goht.svg"
	output := "./goht_16.svg"

	contents, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}

	// split by newlines
	lines := strings.Split(string(contents), "\n")

	newLines := make([]string, 0, len(lines))
	// first three lines are the header
	newLines = append(newLines, lines[0])
	for i := 1; i < len(lines); i++ {
		// read the value from d="<value>"
		// split the value by space
		// convert any numeric values to float and divide by 10
		// join the value back together
		// write the new value back to d="<new value>"
		// write the line back to the newLines
		lines[i] = regexp.MustCompile(`d="([^"]+)"`).ReplaceAllStringFunc(lines[i], func(s string) string {
			num := ""
			r := ""
			for _, c := range s {
				if c != '.' && c != '-' && !(c >= '0' && c <= '9') {
					if num != "" {
						f, err := strconv.ParseFloat(num, 64)
						if err != nil {
							panic(err)
						}
						f /= div
						r += fmt.Sprintf("%f", f)
						num = ""
					}
					r += string(c)
					continue
				}
				num += string(c)
			}
			if num != "" {
				f, err := strconv.ParseFloat(num, 64)
				if err != nil {
					panic(err)
				}
				f /= div
				r += fmt.Sprintf("%f", f)
			}
			return r
		})
		newLines = append(newLines, lines[i])
	}

	// join the newLines by newline
	newContents := strings.Join(newLines, "\n")
	// write the newContents to the output file
	err = os.WriteFile(output, []byte(newContents), 0644)
	if err != nil {
		panic(err)
	}
}
