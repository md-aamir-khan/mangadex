package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func extractChapterIds(mangaUrl string) ([]string, int) {
	// fmt.Println(mangaUrl)
	type ChapterAttributes struct {
		Language string `json:"translatedLanguage"`
		ChapterName string `json:"chapter"`
		Title string `json:"title"`
	}
	type ChapterInfo struct {
		ChapterId  string            `json:"id"`
		Type       string            `json:"type"`
		Attributes ChapterAttributes `json:"attributes"`
	}

	type Manga struct {
		Chapters []ChapterInfo `json:"data"`
	}

	bytes := httpGetAndReturnBytes(mangaUrl)
	var manga Manga
	err := json.Unmarshal(bytes, &manga)
	if err != nil {
		log.Fatalln(mangaUrl+"\n", string(bytes), "\n", err)
	}

	var chapterIds []string

	for _, chap := range manga.Chapters {
		if chap.Type != "chapter" || chap.Attributes.Language != "en" {
			continue
		}
		chapterIds = append(chapterIds, chap.ChapterId)
	}

	return chapterIds, len(manga.Chapters)
}

func createChapterUrlsFromChapterIds(chapterIds []string) []string {
	var chapterUrls []string

	for _, id := range chapterIds {
		url := fmt.Sprintf("https://api.mangadex.org/at-home/server/%s?forcePort443=false", id)
		chapterUrls = append(chapterUrls, url)
	}
	return chapterUrls
}
