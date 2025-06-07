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
    }
    type ChapterInfo struct {
        ChapterId  string            `json:"id"`
        Type       string            `json:"type"`
        Attributes ChapterAttributes `json:"attributes"`
    }

    type Manga struct {
        Manga []ChapterInfo `json:"data"`
    }

    bytes := httpGetAndReturnBytes(mangaUrl)
    var manga Manga
    err := json.Unmarshal(bytes, &manga)
    if err != nil {
        log.Fatalln(mangaUrl+"\n", string(bytes), "\n", err)
    }

    var chapterIds []string

    for _, chap := range manga.Manga {
        if chap.Type != "chapter" || chap.Attributes.Language != "en" {
            continue
        }
        chapterIds = append(chapterIds, chap.ChapterId)
    }

    return chapterIds, len(manga.Manga)
}

func createChapterUrlsFromChapterIds(chapterIds []string) []string {
    var chapterUrls []string

    for _, id := range chapterIds {
        url := fmt.Sprintf("https://api.mangadex.org/at-home/server/%s?forcePort443=false", id)
        chapterUrls = append(chapterUrls, url)
    }
    return chapterUrls
}

func downloadChapterList(chapterListUrl string) {
    chapterIds, _ := extractChapterIds(chapterListUrl)
    chapterUrls := createChapterUrlsFromChapterIds(chapterIds)
    fmt.Println(chapterUrls[0])
    fmt.Println(len(chapterUrls))
}
