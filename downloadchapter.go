package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "path"
)

func getPngsUrlFromChapterUrl(chapterUrl string) []string {
    type ChapterInfo struct {
        Hash        string   `json:"hash"`
        LinksToPngs []string `json:"data"`
    }

    type Chapter struct {
        Chapter ChapterInfo `json:"chapter"`
    }

    bytes := httpGetAndReturnBytes(chapterUrl)
    var chapter Chapter
    err := json.Unmarshal(bytes, &chapter)
    if err != nil {
        panic(err)
    }

    var pngUrls []string

    for _, png := range chapter.Chapter.LinksToPngs {
        pngUrls = append(pngUrls, fmt.Sprintf("https://cmdxd98sb0x3yprd.mangadex.network/data/%s/%s", chapter.Chapter.Hash, png))
    }

    return pngUrls
}

func downloadChapter(chapterUrl string, chapterDir string) {
    pngsUrls := getPngsUrlFromChapterUrl(chapterUrl)

    for idx, pngUrl := range pngsUrls {
        // fmt.Println(idx, pngUrl)
        downloadPng(pngUrl, path.Join(chapterDir, fmt.Sprintf("%04d.png", idx+1)))
    }
}

func downloadPng(pngUrl, fileName string) {
    resp, err := httpGetRespWithRetries(pngUrl, 10)
    if err != nil {
        fmt.Println("Error fetching image:", err)
        return
    }
    defer resp.Body.Close()

    // Check if request was successful
    if resp.StatusCode != http.StatusOK {
        fmt.Println("Failed to download image, status code:", resp.StatusCode)
        return
    }

    // Create file to save image
    file, err := os.Create(fileName)
    if err != nil {
        fmt.Println("Error creating file:", err)
        return
    }
    defer file.Close()

    // Copy response body to file
    _, err = io.Copy(file, resp.Body)
    if err != nil {
        fmt.Println("Error saving image:", err)
        return
    }
}
