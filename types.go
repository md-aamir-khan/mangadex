package main

type MangaCache struct {
	DownloadedWebsiteIndices []int  `yaml:"downloadedWebsiteIndices"`
	Offset                   int    `yaml:"offset"`
	NextIdToCreatePDF        int    `yaml:"nextIdToCreatePDF"`
	Id                       string `yaml:"id"`
	Url                      string `yaml:"url"`
}
