package main

type MangaCache struct {
	MangaDir                 string `yaml:"directory"`
	DownloadedWebsiteIndices []int  `yaml:"downloadedWebsiteIndices"`
	Offset                   int    `yaml:"offset"`
	NextIdToCreatePDF        int    `yaml:"nextIdToCreatePDF"`
	Id                       string `yaml:"id"`
	Name                     string `yaml:"name"`
	Url                      string `yaml:"url"`
}
