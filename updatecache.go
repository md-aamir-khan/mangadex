package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

func (m *MangaCache) notDownloadedCount() int {
	cnt := 0

	for _, chapterInfo := range m.ChaptersInfo {
		if !chapterInfo.Downloaded {
			cnt++
		}
	}
	return cnt
}

func (m *MangaCache) writeToCache(mangaDir string) {
	bytes, err := yaml.Marshal(m)
	if err != nil {
		panic(err)
	}

	cacheFile := getCacheFileLoc(mangaDir)

	file, err := os.OpenFile(cacheFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	file.Close()

	if err != nil {
		panic(err)
	}

}
