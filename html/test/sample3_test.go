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
}
