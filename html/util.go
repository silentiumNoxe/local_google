package html

import (
	"errors"
	"io"
	"strings"
)

var voidElements = []string{
	"area",
	"base",
	"br",
	"col",
	"embed",
	"hr",
	"img",
	"input",
	"link",
	"meta",
	"param",
	"source",
	"track",
	"wbr",
}

func isVoidElement(tag string) bool {
	for _, v := range voidElements {
		if v == tag {
			return true
		}
	}
	return false
}

var ignoreContentTags = []string{
	"script",
	"style",
	"template",
}

func ignoreContent(tag string) bool {
	for _, v := range ignoreContentTags {
		if v == tag {
			return true
		}
	}

	return false
}

func replaceQuotes(s string) string {
	i := strings.Index(s, "\"")
	if i == 0 {
		return s[1 : len(s)-1]
	}

	i = strings.LastIndex(s, "'")
	if i == 0 {
		return s[1 : len(s)-1]
	}

	return s
}

func splitAttributes(s string) []string {
	var attrs []string

	var r = strings.NewReader(s)

	var buf []byte
	var rb = make([]byte, 1)
	var qOpen = false
	var qCloseRune byte

	for {
		_, err := r.Read(rb)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}

		buf = append(buf, rb[0])

		if !qOpen && rb[0] == '=' {
			qOpen = true
			_, err = r.Read(rb)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
			}

			buf = append(buf, rb[0])
			qCloseRune = rb[0]
			continue
		}

		if qOpen && rb[0] == qCloseRune {
			qCloseRune = 0
			qOpen = false
			attrs = append(attrs, string(buf))
			buf = nil
		}
	}

	return attrs
}
