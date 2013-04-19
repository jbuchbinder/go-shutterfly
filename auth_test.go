package shutterfly

import (
	"testing"
)

const (
	TEST_USER = "testsfly@mailinator.com"
	TEST_PASS = "test1ng"
)

func TestAuth(t *testing.T) {
	t.Log("Testing auth methods")
	s := Shutterfly{}
	authToken, err := s.Authorize(TEST_USER, TEST_PASS)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Got authToken = " + authToken)
	uid, err := s.GetUserId()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Got uid = " + uid)
}
