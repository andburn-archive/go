package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

type Lexical []string

func (s Lexical) Len() int {
	return len(s)
}

func (s Lexical) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Lexical) Less(i, j int) bool {
	return s[i] < s[j]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("file argument required")
		os.Exit(1)
	}
	reader, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("failed to open file")
		os.Exit(1)
	}
	defer reader.Close()

	words := make([]string, 0)
	re := regexp.MustCompile("<p>(.+)</p>")
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		match := re.FindStringSubmatch(scanner.Text())
		words = append(words, match[1])
	}
	sort.Sort(Lexical(words))

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	for _, w := range words {
		fmt.Printf("<li>%s</li>\n", w)
	}
}
