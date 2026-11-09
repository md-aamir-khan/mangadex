package main

type ChapterInfo struct {
	Id         int    `yaml:"id"`
	Url        string `yaml:"url"`
	Downloaded bool   `yaml:"downloaded"`
}

type MangaCache struct {
	DownloadedWebsiteIndices []int          `yaml:"downloadedWebsiteIndices"`
	Offset                   int            `yaml:"offset"`
	NextIdToCreatePDF        int            `yaml:"nextIdToCreatePDF"`
	Id                       string         `yaml:"id"`
	Url                      string         `yaml:"url"`
	ChaptersInfo             []*ChapterInfo `yaml:"chaptersInfo"`
}
