package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const API_URL = "https://api.notion.com"

type Parent struct {
	PageID string `json:"page_id"`
}

type TextProperties struct {
	Title Title `json:"title"`
}

type Title struct {
	TitleText []TitleText `json:"title"`
}

type TitleText struct {
	Text Text `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type CreateContent struct {
	Parent     Parent         `json:"parent"`
	Properties TextProperties `json:"properties"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	properties := TextProperties{}
	properties.Title.TitleText = make([]TitleText, 1)
	properties.Title.TitleText[0].Text.Content = "test content"

	requestContent := CreateContent{
		Parent: Parent{
			PageID: os.Getenv("NOTION_DB_ID"),
		},
		Properties: properties,
	}

	b, err := json.Marshal(requestContent)
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest(http.MethodPost, API_URL+"/v1/pages", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_TOKEN"))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Notion-Version", "2021-08-16")

	log.Println(request)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println(resp.StatusCode)

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(b))
}
