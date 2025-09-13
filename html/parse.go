package html

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func Parse(r io.Reader) (*Node, error) {
	var reader = Wrap(r)

	var root = &Node{
		Tag: "ROOT",
	}

	var isIgnoreContent = false
	var ignoreUntilTag string
	for {
		token, err := reader.Next(ignoreUntilTag)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		if token == nil {
			panic("nil token")
		}

		if token.Type == OpenTagToken {
			var c = token.Content
			c = strings.Replace(c, "<", "", 1)
			c = strings.Replace(c, ">", "", 1)

			var i = strings.Index(c, " ")
			var tagName = ""
			if i == -1 {
				tagName = c
			} else {
				tagName = c[:i]
			}

			isIgnoreContent = ignoreContent(tagName)
			if isIgnoreContent {
				ignoreUntilTag = fmt.Sprintf("</%s>", tagName)
				continue
			}

			var isVoid = isVoidElement(tagName)
			var n = &Node{Tag: tagName, Attr: make(map[string]string), Parent: root}

			var attrs = splitAttributes(c[i+1:])

			for _, attr := range attrs {
				var kv = strings.Split(attr, "=")

				if len(kv) == 1 {
					n.Attr[kv[0]] = ""
					continue
				}

				n.Attr[kv[0]] = replaceQuotes(kv[1])
			}

			root.Append(n)
			if !isVoid {
				root = n
			}
			continue
		}

		if token.Type == CloseTagToken {
			var c = token.Content
			c = strings.Replace(c, "</", "", 1)
			c = strings.Replace(c, ">", "", 1)

			if isIgnoreContent && ignoreContent(c) {
				isIgnoreContent = false
				ignoreUntilTag = ""
				continue
			}

			var n = root
			if n.Tag == c {
				n = n.Parent
			}

			root = n
			continue
		}

		if token.Type == TextToken && !isIgnoreContent {
			var c = token.Content
			if len(c) == 0 {
				continue
			}

			root.Content = c
		}
	}

	return root, nil
}

type Reader struct {
	r   io.Reader
	n   int
	buf []byte
	rb  []byte
	prb []byte

	readOpen    bool
	readClose   bool
	readComment bool
}

func Wrap(r io.Reader) *Reader {
	return &Reader{r: r, rb: make([]byte, 1), prb: make([]byte, 1)}
}

func (r *Reader) Next(skipTag string) (*Token, error) {
	if skipTag != "" {
		if err := r.skipTo(skipTag); err != nil {
			return nil, err
		}
		return &Token{Type: CloseTagToken, Content: skipTag}, nil
	}

	var eof = false
	for {
		n, err := r.r.Read(r.rb)
		if err != nil {
			if errors.Is(err, io.EOF) {
				eof = true
			} else {
				return nil, err
			}
		}

		if n == 0 {
			break
		}

		if r.rb[0] == '<' {
			var b = r.buf
			r.buf = nil
			r.write(r.rb)

			r.readOpen = true

			if len(b) > 0 {
				return &Token{Type: TextToken, Content: string(b)}, nil
			}

			if eof {
				break
			}

			continue
		}

		if r.readOpen && r.rb[0] == '>' {
			r.write(r.rb)
			var b = r.buf
			r.buf = nil

			r.readOpen = false

			if len(b) > 0 {
				return &Token{Type: OpenTagToken, Content: string(b)}, nil
			}

			if eof {
				break
			}

			continue
		}

		if r.readClose && r.rb[0] == '>' {
			r.write(r.rb)
			var b = r.buf
			r.buf = nil

			r.readClose = false

			if len(b) > 0 {
				return &Token{Type: CloseTagToken, Content: string(b)}, nil
			}

			if eof {
				break
			}

			continue
		}

		if r.readComment && r.rb[0] == '>' {
			r.write(r.rb)
			r.buf = nil

			r.readComment = false

			if eof {
				break
			}

			continue
		}

		if r.rb[0] == '/' && r.prb[0] == '<' {
			r.write(r.rb)
			r.readClose = true
			r.readOpen = false

			if eof {
				break
			}

			continue
		}

		if r.rb[0] == '!' && r.prb[0] == '<' {
			r.write(r.rb)
			r.readComment = true
			r.readOpen = false

			if eof {
				break
			}

			continue
		}

		r.write(r.rb)
	}

	if eof {
		return nil, io.EOF
	}

	return nil, nil
}

func (r *Reader) skipTo(s string) error {
	var l = len(s)
	var buf []byte
	rb := make([]byte, 1)
	for {
		_, err := r.r.Read(rb)
		if err != nil {
			return err
		}

		buf = append(buf, rb[0])

		if len(buf) >= l {
			if string(buf[len(buf)-l:]) == s {
				return nil
			} else {
				buf = nil
			}
		}
	}
}

func (r *Reader) write(b []byte) {
	r.buf = append(r.buf, b...)
	r.prb[0] = r.rb[0]
}
