package main

import (
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"path"
	"strconv"
)

func downloadManga(mangaCache *MangaCache, limit int, mangaDir string) {
	// We configure limit and offset.
	// Create a chapterIds slice
	var chapterIds []string
	url, err := url.Parse(mangaCache.Url)
	if err != nil {
		panic(err)
	}
	query := url.Query()
	query.Set("limit", strconv.Itoa(limit))
	query.Set("offset", strconv.Itoa(mangaCache.Offset))
	url.RawQuery = query.Encode()

	log.Default().Printf("url: %s\n", url.String())
	ids, totGot := extractChapterIds(url.String()) // Add totGot to offset

	mangaCache.Offset += totGot
// https://api.mangadex.org/manga/bb71e1b9-8773-4b44-85b2-68ed44d124eb/feed?limit=96&includes[]=scanlation_group&includes[]=user&order[volume]=asc&order[chapter]=asc&offset=0&contentRating[]=safe&contentRating[]=suggestive&contentRating[]=erotica&contentRating[]=pornographic

	chapterIds = append(chapterIds, ids...)

	chapterUrls := createChapterUrlsFromChapterIds(chapterIds)

	tmpDir := path.Join(mangaDir, TmpDir)
	for _, chapterUrl := range chapterUrls {
		// Remove .tmp/*
		deleteFilesInDir(tmpDir)
		pdfFilePath := path.Join(mangaDir, fmt.Sprintf("%05d.pdf", mangaCache.NextIdToCreatePDF))

		log.Default().Printf("Downloading chapter: %s\n", pdfFilePath)

		downloadChapter(chapterUrl, tmpDir)

		// Add cmd to convert to pdf
		cmd := exec.Command("magick", "convert", fmt.Sprintf("%s/*", tmpDir), pdfFilePath)

		if err := cmd.Run(); err != nil {
			panic(err)
		}

		log.Default().Printf("PDF File Created %s\n", pdfFilePath)

		mangaCache.NextIdToCreatePDF++
	}
}
