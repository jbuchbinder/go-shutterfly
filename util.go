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

// Add Shutterfly Open API authorization headers to an existing
// http.Request object.
func authHeaders(req *http.Request, sfly *Shutterfly, urlParams string) {
	now := time.Now()

	req.Header.Set("oflyAppId", APP_ID)
	req.Header.Set("oflyHashMeth", "SHA1")
	req.Header.Set("oflyTimestamp", now.Format("2006-01-02T15:04:05.000-0700"))
	if sfly.AuthToken != "" {
		req.Header.Set("X-OPENFLY-Authorization", "SFLY user-auth="+sfly.AuthToken)
	}
}
