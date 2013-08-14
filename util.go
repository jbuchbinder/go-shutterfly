package shutterfly

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"time"
)

// Wrap xml.Escape() so that it can be simply used with strings
func XmlEncode(s string) string {
	w := bytes.NewBuffer([]byte{})
	xml.Escape(w, []byte(s))
	return w.String()
}

func AuthHeaders(req *http.Request, urlParams string) {
	now := time.Now()

	req.Header.Set("oflyAppId", APP_ID)
	req.Header.Set("oflyHashMeth", "SHA1")
	req.Header.Set("oflyTimestamp", now.Format("2006-01-02T15:04:05.000-0700"))
}
