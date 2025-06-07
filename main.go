package main

import (
    "fmt"
    "log"
    "net/url"
    "os"
    "path"

    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"
)

/*
MangaName/ChapterNumber(in string)/`ChapterName - userName`
*/
const (
    MangaDexBase = `https://api.mangadex.org/manga/%s/feed?limit=96&includes[]=scanlation_group&includes[]=user&order[volume]=asc&order[chapter]=asc&offset=0&contentRating[]=safe&contentRating[]=suggestive&contentRating[]=erotica&contentRating[]=pornographic`
    CacheFilename = `.cache.yaml`
    TmpDir = ".tmp"
)

type MangaInfo struct {
    Url   *url.URL
    Limit int
}

func main() {
    // mangaId := "bb71e1b9-8773-4b44-85b2-68ed44d124eb"

    rootCmd := &cobra.Command{
        Use:   "mangadex plugin",
        Short: "Mangadex Download, List",
    }

    downloadCmd := downloadMangaCmd()

    rootCmd.AddCommand(downloadCmd)

    rootCmd.AddCommand(initDirCmd())

    rootCmd.Execute()
}

func initDirCmd() *cobra.Command {
    var id string
    var parentDir string
    var mangaName string
    initDirCmd := &cobra.Command{
        Use:   "initdir",
        Short: "Inits manga directory with cache",
        Run: func(cmd *cobra.Command, args []string) {
            mangaUrl, err := url.Parse(fmt.Sprintf(MangaDexBase, id))
            if err != nil {
                panic(err)
            }
            if !isDirectory(parentDir) {
                panic("directory is not valid")
            }
            mangaDir := path.Join(parentDir, mangaName)
            if err := os.Mkdir(mangaDir, 0755); err != nil {
                panic(err)
            }
            if err := os.Mkdir(path.Join(mangaDir, TmpDir), 0755); err != nil {
                panic(err)
            }

            mangaCache := MangaCache{
                MangaDir:                 mangaDir,
                Id:                       id,
                Name:                     mangaName,
                Offset:                   0,
                NextIdToCreatePDF:        1,
                DownloadedWebsiteIndices: []int{},
                Url:                      mangaUrl.String(),
            }

            fileName := path.Join(mangaDir, CacheFilename)
            // Create file to save image
            file, err := os.Create(fileName)
            if err != nil {
                fmt.Println("Error creating file:", err)
                return
            }
            defer file.Close()

            writeYamlToFile(mangaCache, file)
        },
    }
    initDirCmd.Flags().StringVarP(&id, "id", "", "", "id of the manga in mangadex")
    initDirCmd.Flags().StringVarP(&parentDir, "dir", "", "", "parent directory for downloading manga")
    initDirCmd.Flags().StringVarP(&mangaName, "name", "", "", "manga Name")

    initDirCmd.MarkFlagRequired("id")
    initDirCmd.MarkFlagRequired("dir")
    initDirCmd.MarkFlagRequired("name")

    return initDirCmd
}

func downloadMangaCmd() *cobra.Command {
    var limit int
    var parentDir string
    var mangaName string

    fmt.Println(parentDir);

    downloadMangaCmd := &cobra.Command{
        Use:   "downloadmanga",
        Short: "Download manga with given id and to a directory",
        Run: func(cmd *cobra.Command, args []string) {
            cacheFile := path.Join(parentDir, mangaName, CacheFilename)
            if !fileExists(cacheFile) {
                panic("cache file does not exist")
            }

            // Read cache file
            data, err := os.ReadFile(cacheFile)
            if err != nil {
                panic(err)
            }

            var mangaCache MangaCache
            if err := yaml.Unmarshal(data, &mangaCache); err != nil {
                panic(err)
            }
            if !isDirectory(path.Join(mangaCache.MangaDir, TmpDir)) {
                panic("tmp dir does not exist")
            }

            log.Default().Println("Starting to download manga")
            downloadManga(&mangaCache, limit)

            file, err := os.OpenFile(cacheFile, os.O_WRONLY|os.O_CREATE, 0644)
            if err != nil {
                panic(err)
            }
            writeYamlToFile(mangaCache, file)
        },
    }
    downloadMangaCmd.Flags().IntVarP(&limit, "limit", "", 50, "Describes the number of chapters to fetch in a batch")
    downloadMangaCmd.Flags().StringVarP(&parentDir, "dir", "", "", "parent directory for downloading manga")
    downloadMangaCmd.Flags().StringVarP(&mangaName, "name", "", "", "manga Name")

    downloadMangaCmd.MarkFlagRequired("dir")
    downloadMangaCmd.MarkFlagRequired("name")

    return downloadMangaCmd
}
