package shutterfly

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	//"net/http/httputil"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

type uploadResponse struct {
	XMLName    xml.Name `xml:"feed"`
	ErrCode    int      `xml:"upload:errCode"`
	ErrMessage string   `xml:"upload:errMessage"`
	NumSuccess int      `xml:"upload:numSuccess"`
	NumFail    int      `xml:"upload:numFail"`
	AlbumPath  string   `xml:"upload:albumPath"`
}

// UploadPhotos uploads photos to the Shutterfly service with the specified
// folder and album names.
func (self *Shutterfly) UploadPhotos(photos []string, folderName, albumName string) error {
	if self.AuthToken == "" {
		return errors.New("no AuthToken, please login first")
	}
	client := http.Client{}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add parameters
	_ = writer.WriteField("AuthenticationID", self.AuthToken)
	_ = writer.WriteField("Image.FolderName", folderName)
	_ = writer.WriteField("Image.AlbumName", albumName)

	for _, v := range photos {
		self.TestLog("Iterating through image " + v)
		file, err := os.Open(v)
		if err != nil {
			return err
		}
		defer file.Close()
		part, err := createJpegFormFile(writer /* "photo["+string(k)+"]" */, "Image.Data", filepath.Base(v))
		if err != nil {
			return err
		}
		_, err = io.Copy(part, file)
		//for key, val := range params {
		//	_ = writer.WriteField(key, val)
		//}
	}
	err := writer.Close()
	if err != nil {
		return err
	}

	// This was originally documented to use up3.shutterfly.com, but the
	// testing mechanism indicated this URL, so it has been adjusted to
	// use it.
	req, err := http.NewRequest("POST", "https://up1.shutterfly.com/images", body)
	if err != nil {
		fmt.Println("ERROR: ")
		fmt.Println(err)
		return err
	}

	// Push MIME boundary type header from multipart.Writer
	req.Header.Set("Content-type", writer.FormDataContentType())

	// Populate with authorization request headers
	authHeaders(req, self, "")

	//dump, _ := httputil.DumpRequest(req, true)
	//self.TestLog("Request: " + string(dump))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: ")
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()
	rbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("IO: ")
		fmt.Println(err)
		fmt.Println(rbody)
		return err
	}
	self.TestLog("Status: " + res.Status)
	self.TestLog(string(rbody))

	// Handle non-2xx/3xx series error codes by throwing an error
	if res.StatusCode > 399 {
		return errors.New("response: " + res.Status)
	}

	var xmlresp uploadResponse
	err = xml.Unmarshal(rbody, &xmlresp)
	if err != nil {
		return err
	}

	// Currently, the logic is that *any* failure in a batch upload
	// should indicate a failure.
	if xmlresp.NumFail > 0 {
		return fmt.Errorf("%d failed to upload", xmlresp.NumFail)
	}

	return nil
}

// Cribbed functions to support overriding MIME stuff

func createJpegFormFile(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	// Major change from the stock function is to declare this with
	// type == image/jpeg, otherwise Shutterfly's API will choke on
	// the request:
	h.Set("Content-Type", "image/jpeg")
	return w.CreatePart(h)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
