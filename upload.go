package shutterfly

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (self *Shutterfly) UploadPhotos(photos []string) error {
	if self.AuthToken == "" {
		return errors.New("No AuthToken, please login first.")
	}
	client := http.Client{}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range photos {
		file, err := os.Open(v)
		if err != nil {
			return err
		}
		defer file.Close()
		part, err := writer.CreateFormFile("photo["+string(k)+"]", filepath.Base(v))
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

	req, err := http.NewRequest("POST", "https://up3.shutterfly.com/images", body)
	if err != nil {
		fmt.Println("ERROR: ")
		fmt.Println(err)
		return err
	}
	req.Header.Set("X-OPENFLY-Authorization", "SFLY user-auth="+self.AuthToken)
	AuthHeaders(req, "")
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

	return nil
}
