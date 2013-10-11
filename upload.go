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

type UploadResponse struct {
	XMLName    xml.Name `xml:"feed"`
	ErrCode    int      `xml:"upload:errCode"`
	ErrMessage string   `xml:"upload:errMessage"`
	NumSuccess int      `xml:"upload:numSuccess"`
	NumFail    int      `xml:"upload:numFail"`
	AlbumPath  string   `xml:"upload:albumPath"`
}

func (self *Shutterfly) UploadPhotos(photos []string, folderName, albumName string) error {
	if self.AuthToken == "" {
		return errors.New("No AuthToken, please login first.")
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
		part, err := CreateJpegFormFile(writer /* "photo["+string(k)+"]" */, "Image.Data", filepath.Base(v))
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

	req, err := http.NewRequest("POST", "https://up1.shutterfly.com/images", body)
	if err != nil {
		fmt.Println("ERROR: ")
		fmt.Println(err)
		return err
	}

	req.Header.Set("X-OPENFLY-Authorization", "SFLY user-auth="+self.AuthToken)
	req.Header.Set("Content-type", writer.FormDataContentType())
	self.TestLog("X-OPENFLY-Authorization: SFLY user-auth=" + self.AuthToken)
	AuthHeaders(req, "")

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

	if res.StatusCode > 399 {
		return errors.New("Response: " + res.Status)
	}

	var xmlresp UploadResponse
	err = xml.Unmarshal(rbody, &xmlresp)
	if err != nil {
		return err
	}
	if xmlresp.NumFail > 0 {
		return errors.New(fmt.Sprintf("%d failed to upload", xmlresp.NumFail))
	}

	return nil
}

func CreateJpegFormFile(w *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", "image/jpeg")
	return w.CreatePart(h)
}

// Cribbed functions to support overriding MIME stuff

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
