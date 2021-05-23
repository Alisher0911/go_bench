package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// вам надо написать более быструю оптимальную этой функции

type User struct {
	Browsers []string `json:"browsers"`
	Email string `json:"email"`
	Name string `json:"name"`
}

var dataPool = sync.Pool{
	New: func() interface{} {
		return &User{}
	},
}


func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Errorf("cannot open file")
	}

	in := bufio.NewScanner(file)
	seenBrowsers := make(map[string]struct{})
	count := 0

	fmt.Fprintln(out, "found users:")

	for in.Scan() {
		count++

		line := in.Bytes()
		user := dataPool.Get().(*User)

		err = json.Unmarshal(line, user)
		if err != nil {
			fmt.Errorf("cannot decode json")
		}

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			browserChecked := false

			if strings.Contains(browser, "Android") {
				isAndroid = true
				browserChecked = true
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				browserChecked = true
			}

			if browserChecked {
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
				}
			}
		}

		dataPool.Put(user)

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.Replace(user.Email,  "@", " [at] ", -1)
		fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", count - 1, user.Name, email))
	}

	fmt.Fprintln(out, "\nTotal unique browsers", len(seenBrowsers))
}