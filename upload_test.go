package shutterfly

import (
	"testing"
)

func TestUpload(t *testing.T) {
	t.Log("Testing image upload")

	s := Shutterfly{
		Test: t,
	}
	_, err := s.Authorize(TEST_USER, TEST_PASS)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.UploadPhotos([]string{"stumped.jpg"}, "Test Folder", "Test Album")
	if err != nil {
		t.Error(err)
		return
	}
}
