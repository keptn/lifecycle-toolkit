package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const split = "---"

var file string
var r *regexp.Regexp

var l = log.New(os.Stderr, "", 0)

func init() {
	args := os.Args[1:]
	if len(args) < 1 {
		l.Println("expected folder")
		os.Exit(1)
	}
	file = args[0]
	r = regexp.MustCompile("title:\\s(.+)")
}

func main() {
	b, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	content := string(b)
	out := content
	if strings.Contains(content, split) {
		s := strings.Split(content, split)
		pre := s[1]
		rest := s[2]
		m := r.FindStringSubmatch(pre)
		if len(m) > 0 {
			out = fmt.Sprintf("# %s%s", m[1], rest)
		} else {
			l.Printf("Warning! Missing title for %s\n", file)
			out = rest
		}
	}
	os.WriteFile(file, []byte(out), 0644)
	//fmt.Println(out)
}
