package shutterfly

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func (self *Shutterfly) Authorize(username, password string) (string, error) {
	client := http.Client{}
	payload := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
		"<entry xmlns=\"http://www.w3.org/2005/Atom\" xmlns:user=\"http://user.openfly.shutterfly.com/v1.0\">" +
		"<category term=\"user\" scheme=\"http://openfly.shutterfly.com/v1.0\" />" +
		"<user:password>" + XmlEncode(password) + "</user:password>" +
		"</entry>\n"
	url := "https://ws.shutterfly.com/user/" + username + "/auth"

	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		fmt.Println("ERROR: ")
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/atom+xml")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: ")
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("IO: ")
		fmt.Println(err)
		fmt.Println(body)
		return "", err
	}

	b := string(body)

	rx, err := regexp.Compile(`<user:newAuthToken>([^<]+)</user:newAuthToken>`)
	if err != nil {
		fmt.Println("REGEX: ")
		fmt.Println(err)
		fmt.Println(body)
		return "", err
	}

	sm := rx.FindStringSubmatch(b)
	self.AuthToken = sm[1]
	return sm[1], nil
}

func (self *Shutterfly) GetUserId() (string, error) {
	if self.AuthToken == "" {
		return "", errors.New("No AuthToken, please login first.")
	}
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://ws.shutterfly.com/auth", nil)
	if err != nil {
		fmt.Println("ERROR: ")
		fmt.Println(err)
		return "", err
	}
	req.Header.Set("X-OPENFLY-Authorization", "SFLY user-auth="+self.AuthToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: ")
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("IO: ")
		fmt.Println(err)
		fmt.Println(body)
		return "", err
	}

	b := string(body)

	rx, err := regexp.Compile(`<openfly:userid>([^<]+)</openfly:userid>`)
	if err != nil {
		fmt.Println("REGEX: ")
		fmt.Println(err)
		fmt.Println(body)
		return "", err
	}

	sm := rx.FindStringSubmatch(b)
	return sm[1], nil
}
