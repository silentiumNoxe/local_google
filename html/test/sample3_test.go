package test

import (
	"local_google/html"
	"strings"
	"testing"
)

func TestSkipTags(t *testing.T) {
	var s = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Simple HTML</title>
    <style>xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx</style>
</head>
<body>
<div id="main-content" class="container">
    <script>xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx</script>
    <h1>An interesting title</h1>
    <p>
        Here is a paragraph with a
        <a href="/about-us">link</a>
        and an empty tag.
    </p>
    <!--Comment-->
    <img src="/logo.png" alt="Company Logo">
</div>
<script>xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx</script>
</body>
</html>
`
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

	elem := root.FirstChild
	if elem == nil {
		t.Errorf("No child")
		return
	}

	if elem.Tag != "html" {
		t.Errorf("Invalid tag: expected=html; actual=%s", elem.Tag)
	}

	var head = elem.FirstChild
	if head == nil {
		t.Errorf("No child head")
		return
	}

	if head.Tag != "head" {
		t.Errorf("Invalid tag: expected=head; actual=%s", head.Tag)
		return
	}

	if head.FirstChild.Content != "Simple HTML" {
		t.Errorf("Invalid content of tag title: expected=\"Simple HTML\"; actual=%s", head.FirstChild.Content)
	}

	var body = head.Next
	if body == nil {
		t.Errorf("No child body")
		return
	}

	if body.Tag != "body" {
		t.Errorf("Invalid tag: expected=body; actual=%s", body.Tag)
		return
	}

	var aElement = body.FirstChild.FirstChild.Next.FirstChild
	if aElement == nil {
		t.Errorf("No child a")
		return
	}

	var href = aElement.Attr["href"]
	if href != "/about-us" {
		t.Errorf("Invalid href: expected=\"/about-us\"; actual=%s", href)
		return
	}

	var img = body.FirstChild.FirstChild.Next.Next
	if img == nil {
		t.Errorf("No child img")
		return
	}

	if img.Attr["alt"] != "Company Logo" {
		t.Errorf("Invalid alt: expected=\"Company Logo\"; actual=%s", img.Attr["alt"])
		return
	}
}

func TestSkipTagsLvl2(t *testing.T) {
	var s = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Simple HTML</title>
    <style>
		const test = 100 > 10
		console.log(test);
	</style>
</head>
<body>
<div id="main-content" class="container">
    <script>xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx</script>
    <h1>An interesting title</h1>
    <p>
        Here is a paragraph with a
        <a href="/about-us">link</a>
        and an empty tag.
    </p>
    <!--Comment-->
    <img src="/logo.png" alt="Company Logo">
</div>
<script>xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx</script>
</body>
</html>
`
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

	elem := root.FirstChild
	if elem == nil {
		t.Errorf("No child")
		return
	}

	if elem.Tag != "html" {
		t.Errorf("Invalid tag: expected=html; actual=%s", elem.Tag)
	}

	var head = elem.FirstChild
	if head == nil {
		t.Errorf("No child head")
		return
	}

	if head.Tag != "head" {
		t.Errorf("Invalid tag: expected=head; actual=%s", head.Tag)
		return
	}

	if head.FirstChild.Content != "Simple HTML" {
		t.Errorf("Invalid content of tag title: expected=\"Simple HTML\"; actual=%s", head.FirstChild.Content)
	}

	var body = head.Next
	if body == nil {
		t.Errorf("No child body")
		return
	}

	if body.Tag != "body" {
		t.Errorf("Invalid tag: expected=body; actual=%s", body.Tag)
		return
	}

	var aElement = body.FirstChild.FirstChild.Next.FirstChild
	if aElement == nil {
		t.Errorf("No child a")
		return
	}

	var href = aElement.Attr["href"]
	if href != "/about-us" {
		t.Errorf("Invalid href: expected=\"/about-us\"; actual=%s", href)
		return
	}

	var img = body.FirstChild.FirstChild.Next.Next
	if img == nil {
		t.Errorf("No child img")
		return
	}

	if img.Attr["alt"] != "Company Logo" {
		t.Errorf("Invalid alt: expected=\"Company Logo\"; actual=%s", img.Attr["alt"])
		return
	}
}
