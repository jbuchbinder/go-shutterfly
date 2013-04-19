package shutterfly

import (
	"bytes"
	"encoding/xml"
)

// Wrap xml.Escape() so that it can be simply used with strings
func XmlEncode(s string) string {
	w := bytes.NewBuffer([]byte{})
	xml.Escape(w, []byte(s))
	return w.String()
}
