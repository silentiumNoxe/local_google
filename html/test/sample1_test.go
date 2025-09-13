package test

import (
	"local_google/html"
	"strings"
	"testing"
)

func TestSample1(t *testing.T) {
	var s = "<div>hello world</div>"
	root, err := html.Parse(strings.NewReader(s))
	if err != nil {
		t.Errorf("Got error; %s", err.Error())
		return
	}

	if root == nil {
		t.Errorf("Got nil root")
		return
	}

	if root.Tag != "ROOT" {
		t.Errorf("Invalid tag: expected=ROOT; actual=%s", root.Tag)
		return
	}

	var div = root.FirstChild
	if div == nil {
		t.Errorf("No child")
		return
	}

	if div.Content != "hello world" {
		t.Errorf("Invalid content: expected=%s; actual=%s", "\"hello world\"", div.Content)
		return
	}
}
