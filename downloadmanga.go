package main

import (
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"path"
	"strconv"
)

const (
	QUERY_LIMIT = 50
)

func queryAndGetChapterUrls(mangaCache *MangaCache) []string {
	var chapterIds []string
	url, err := url.Parse(mangaCache.Url)
	if err != nil {
		panic(err)
	}
	query := url.Query()
	query.Set("limit", strconv.Itoa(QUERY_LIMIT))
	query.Set("offset", strconv.Itoa(mangaCache.Offset))
	url.RawQuery = query.Encode()

	log.Default().Printf("url: %s\n", url.String())
	ids, totGot := extractChapterIds(url.String()) // Add totGot to offset

	mangaCache.Offset += totGot
	// https://api.mangadex.org/manga/bb71e1b9-8773-4b44-85b2-68ed44d124eb/feed?limit=96&includes[]=scanlation_group&includes[]=user&order[volume]=asc&order[chapter]=asc&offset=0&contentRating[]=safe&contentRating[]=suggestive&contentRating[]=erotica&contentRating[]=pornographic

	chapterIds = append(chapterIds, ids...)

	chapterUrls := createChapterUrlsFromChapterIds(chapterIds)

	return chapterUrls
}

func downloadChapterListAndUpdateCache(mangaCache *MangaCache, limit int, mangaDir string) {

	for limit > mangaCache.notDownloadedCount() {

		chapterUrls := queryAndGetChapterUrls(mangaCache)
		if len(chapterUrls) == 0 {
			fmt.Println("Reached Last Manga Chapter Info")
			break
		}
		nextChapterId := 1
		chaptersInfo := mangaCache.ChaptersInfo
		if len(chaptersInfo) > 0 {
			nextChapterId = chaptersInfo[len(chaptersInfo)-1].Id + 1
		}

		for _, chapterUrl := range chapterUrls {
			mangaCache.ChaptersInfo = append(mangaCache.ChaptersInfo, &ChapterInfo{
				Url:        chapterUrl,
				Id:         nextChapterId,
				Downloaded: false,
			})
			nextChapterId++
		}

		mangaCache.writeToCache(mangaDir)
	}
}

func downloadManga(mangaCache *MangaCache, limit int, mangaDir string) {
	// We configure limit and offset.
	// Create a chapterIds slice

	downloadChapterListAndUpdateCache(mangaCache, limit, mangaDir)

	tmpDir := path.Join(mangaDir, TmpDir)

	curDownloadedCnt := 0
	for _, chapterInfo := range mangaCache.ChaptersInfo {
		// Remove .tmp/*
		if chapterInfo.Downloaded {
			continue
		}
		deleteFilesInDir(tmpDir)
		pdfFilePath := path.Join(mangaDir, fmt.Sprintf("%05d.pdf", chapterInfo.Id))

		log.Default().Printf("Downloading chapter: %s\n", pdfFilePath)

		downloadChapter(chapterInfo.Url, tmpDir)

		// Add cmd to convert to pdf
		cmd := exec.Command("magick", "convert", fmt.Sprintf("%s/*", tmpDir), pdfFilePath)

		if err := cmd.Run(); err != nil {
			panic(err)
		}

		log.Default().Printf("PDF File Created %s\n", pdfFilePath)

		chapterInfo.Downloaded = true

		mangaCache.writeToCache(mangaDir)

		curDownloadedCnt++

		if curDownloadedCnt >= limit {
			break
		}
	}
}
