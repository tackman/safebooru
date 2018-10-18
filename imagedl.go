package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
)

type DatastoreItem struct {
	Directory    string   `json:"directory"`
	Hash         string   `json:"hash"`
	Height       int      `json:"height"`
	Id           int      `json:"Id"`
	Image        string   `json:"image"`
	Change       int      `json:"change"`
	Owner        string   `json:"owner"`
	ParentId     int      `json:"parent_id"`
	Rating       string   `json:"rating"`
	Sample       bool     `json:"sample"`
	SampleHeight int      `json:"sample_height"`
	SampleWidth  int      `json:"sample_width"`
	Score        int      `json:"score"`
	Tags         []string `json:"tags"`
	Width        int      `json:"width"`
}

func main() {
	const Kind = "SafebooruImage"
	const UrlFormat = "https://safebooru.org//images/%s/%s"
	const PathFormat = "/home/takuma/cvproj/white_1girl/%s-%s"

	ctx := context.Background()

	dsCli, _ := datastore.NewClient(ctx, "cvproject-166719")
	defer dsCli.Close()

	query := datastore.NewQuery(Kind)

	var item DatastoreItem
	for it := dsCli.Run(ctx, query); ; {
		_, err := it.Next(&item)

		if err != nil {
			log.Fatal(err)
			break
		}

		url := fmt.Sprintf(UrlFormat, item.Directory, item.Image)
		path := fmt.Sprintf(PathFormat, item.Directory, item.Image)
		dlImage(url, path)
	}

	return
}

func dlImage(url string, path string) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()

	file, fileErr := os.Create(path)

	if fileErr != nil {
		log.Fatal(err)
		panic("Raised file ERROR!!")
	}


	_, writeError := io.Copy(file, resp.Body)

	if writeError != nil {
		log.Fatal(writeError)
		return
	}

	file.Close()
}
