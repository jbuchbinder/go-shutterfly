package shutterfly

import (
	"errors"
	"testing"
)

func TestXmlEncode(t *testing.T) {
	t.Log("Testing XML encoding")
	orig := "This is & <> test"
	out := XmlEncode(orig)
	if out == "" {
		t.Error(errors.New("No string returned"))
		return
	}
	t.Log("Returned : " + out)
	if out != "This is &amp; &lt;&gt; test" {
		t.Error(errors.New("Fail"))
	}
}
