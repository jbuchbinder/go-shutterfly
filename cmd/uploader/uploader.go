package main

import (
	"flag"
	"fmt"

	shutterfly "github.com/jbuchbinder/go-shutterfly"
)

var (
	Username = flag.String("username", "", "Shutterfly username")
	Password = flag.String("password", "", "Shutterfly password")
	Folder   = flag.String("folder", "My Folder", "Shutterfly folder name")
	Album    = flag.String("album", "My Album", "Shutterfly album name")
)

func main() {
	flag.Parse()
	if *Username == "" || *Password == "" {
		flag.PrintDefaults()
		return
	}

	images := flag.Args()

	if len(images) < 1 {
		fmt.Println("You must specify one or more images to upload.")
		return
	}

	s := shutterfly.Shutterfly{}
	_, err := s.Authorize(*Username, *Password)
	if err != nil {
		fmt.Println("Failed to authenticate with Shutterfly!")
		panic(err)
	}

	err = s.UploadPhotos(images, *Folder, *Album)
	if err != nil {
		panic(err)
	}

	fmt.Printf("SUCCESS: %d photos uploaded.\n", len(images))
}
