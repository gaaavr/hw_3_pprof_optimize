package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"hw3/models"
)

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	foundUsers := strings.Builder{}

	seenBrowsers := make(map[string]struct{}, 0)

	for i := 0; scanner.Scan(); i++ {
		user := models.User{}

		err = user.UnmarshalJSON(scanner.Bytes())
		if err != nil {
			panic(err)
		}

		var isAndroid, isMSIE bool

		for _, browser := range user.Browsers {
			hasAndroid := strings.Contains(browser, "Android")
			hasMSIE := strings.Contains(browser, "MSIE")

			_, ok := seenBrowsers[browser]

			isMSIE = hasMSIE || isMSIE
			isAndroid = hasAndroid || isAndroid

			if !ok && (hasAndroid || hasMSIE) {
				seenBrowsers[browser] = struct{}{}
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		foundUsers.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
