package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
)

type SafebooruItem struct {
	Directory    string `json:"directory"`
	Hash         string `json:"hash"`
	Height       int    `json:"height"`
	Id           int    `json:"Id"`
	Image        string `json:"image"`
	Change       int    `json:"change"`
	Owner        string `json:"owner"`
	ParentId     int    `json:"parent_id"`
	Rating       string `json:"rating"`
	Sample       bool   `json:"sample"`
	SampleHeight int    `json:"sample_height"`
	SampleWidth  int    `json:"sample_width"`
	Score        int    `json:"score"`
	Tags         string `json:"tags"`
	Width        int    `json:"width"`
}

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
	const PageCount int = 158759

	ctx := context.Background()

	dsCli, _ := datastore.NewClient(ctx, "cvproject-166719")
	defer dsCli.Close()

	for page := 0; page < PageCount; page++ {
		pageUrl := fmt.Sprintf("http://safebooru.org/index.php?page=dapi&s=post&q=index&limit=100&json=1&tags=white_background+1girl&pid=%d", page)

		fmt.Printf("begin download %s", pageUrl)

		b, httpErr := getPage(pageUrl)

		if httpErr != nil {
			log.Fatal(httpErr)
			continue
		}

		items, jsonErr := decodeJson(b)

		if jsonErr != nil {
			log.Fatal(jsonErr)
			continue
		}

		if len(items) <= 0 {
			log.Println("All jobs complete")
			break
		}

		storeItems(&ctx, dsCli, items)
	}
}

func getPage(page string) ([]byte, error) {
	resp, err := http.Get(page)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	b, ioErr := ioutil.ReadAll(resp.Body)

	if ioErr != nil {
		log.Fatal(ioErr)
		return nil, ioErr
	}

	return b, nil
}

func storeItems(ctx *context.Context, dsCli *datastore.Client, safeboorus []SafebooruItem) {
	for _, item := range safeboorus {
		slicedTags := strings.Split(item.Tags, " ")

		dbItem := &DatastoreItem{
			Directory:    item.Directory,
			Hash:         item.Hash,
			Id:           item.Id,
			Image:        item.Image,
			Change:       item.Change,
			Owner:        item.Owner,
			ParentId:     item.ParentId,
			Rating:       item.Rating,
			Sample:       item.Sample,
			SampleHeight: item.SampleHeight,
			SampleWidth:  item.SampleWidth,
			Score:        item.Score,
			Tags:         slicedTags,
			Width:        item.Width}

		putDb(ctx, dsCli, dbItem)
	}
}

func decodeJson(b []byte) ([]SafebooruItem, error) {
	var safeboorus []SafebooruItem
	jsonErr := json.Unmarshal(b, &safeboorus)

	if jsonErr != nil {
		log.Fatal(jsonErr)
		return nil, jsonErr
	}

	return safeboorus, nil
}

func putDb(ctx *context.Context, dsCli *datastore.Client, dbItem *DatastoreItem) {
	ancesttralKey, _ := datastore.DecodeKey("ahJifmN2cHJvamVjdC0xNjY3MTlyLAsSDVNhZmVib29ydVJvb3QiGVdoaXRlQmFja2dyb3VuZDFHaXJsSW1hZ2UMogEJU2FmZWJvb3J1")
	itemKey := datastore.NameKey("SafebooruImage", strconv.Itoa(dbItem.Id), ancesttralKey)

	_, err := dsCli.Put(*ctx, itemKey, dbItem)

	if err != nil {
		log.Fatal(err)
	}
}
