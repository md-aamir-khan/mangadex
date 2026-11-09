package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"mangadex.com/download/auth"
)

/*
MangaName/ChapterNumber(in string)/`ChapterName - userName`
*/
const (
	MangaDexBase  = `https://api.mangadex.org/manga/%s/feed?limit=96&includes[]=scanlation_group&includes[]=user&order[volume]=asc&order[chapter]=asc&offset=0&contentRating[]=safe&contentRating[]=suggestive&contentRating[]=erotica&contentRating[]=pornographic`
	CacheFilename = `.cache.yaml`
	TmpDir        = ".tmp"
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

	loginCmd := loginMangaCmd()

	downloadCmd := downloadMangaCmd()

	rootCmd.AddCommand(downloadCmd)

	rootCmd.AddCommand(initDirCmd())

	rootCmd.AddCommand(loginCmd)

	rootCmd.Execute()
}

func initDirCmd() *cobra.Command {
	var id string
	var mangaDir string
	initDirCmd := &cobra.Command{
		Use:   "initdir",
		Short: "Inits manga directory with cache",
		Run: func(cmd *cobra.Command, args []string) {
			mangaUrl, err := url.Parse(fmt.Sprintf(MangaDexBase, id))
			if err != nil {
				panic(err)
			}
			if err := os.Mkdir(mangaDir, 0755); err != nil {
				panic(err)
			}
			if err := os.Mkdir(path.Join(mangaDir, TmpDir), 0755); err != nil {
				panic(err)
			}

			mangaCache := MangaCache{
				Id:                       id,
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
	initDirCmd.Flags().StringVarP(&mangaDir, "dir", "", "", "manga directory for downloading manga")

	initDirCmd.MarkFlagRequired("id")
	initDirCmd.MarkFlagRequired("dir")

	return initDirCmd
}

func loginMangaCmd() *cobra.Command {
	var username string
	var password string

	loginMangaCmd := &cobra.Command{
		Use:   "login",
		Short: "login manga with given id and to a directory",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(username, '\n', password)
			_, err := auth.Login(username, password)
			if err != nil {
				panic(err)
			}
			fmt.Println("login successful")
		},
	}
	loginMangaCmd.Flags().StringVarP(&username, "username", "u", "", "Describes the number of chapters to fetch in a batch")
	loginMangaCmd.Flags().StringVarP(&password, "password", "p", "", "manga directory for downloading manga")

	loginMangaCmd.MarkFlagRequired("username")
	loginMangaCmd.MarkFlagRequired("password")

	return loginMangaCmd
}

func downloadMangaCmd() *cobra.Command {
	var limit int
	var mangaDir string

	fmt.Println(mangaDir)

	downloadMangaCmd := &cobra.Command{
		Use:   "downloadmanga",
		Short: "Download manga with given id and to a directory",
		Run: func(cmd *cobra.Command, args []string) {
			cacheFile := path.Join(mangaDir, CacheFilename)
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
			tmpDir := path.Join(mangaDir, TmpDir)
			if !isDirectory(tmpDir) {
				panic(fmt.Sprintf("tmp dir does not exist: %s", tmpDir))
			}

			log.Default().Println("Starting to download manga")
			downloadManga(&mangaCache, limit, mangaDir)

		},
	}
	downloadMangaCmd.Flags().IntVarP(&limit, "limit", "", 50, "Describes the number of chapters to fetch in a batch")
	downloadMangaCmd.Flags().StringVarP(&mangaDir, "dir", "", "", "manga directory for downloading manga")

	downloadMangaCmd.MarkFlagRequired("dir")
	downloadMangaCmd.MarkFlagRequired("limit")

	return downloadMangaCmd
}
