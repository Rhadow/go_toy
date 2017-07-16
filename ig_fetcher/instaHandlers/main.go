package instaHandlers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"

	"github.com/ahmdrz/goinsta"
	"github.com/rhadow/go_toy/ig_fetcher/fetcherResponse"
)

// LoginInstagram - Log in instagram
func LoginInstagram(username, password string) *goinsta.Instagram {
	insta := goinsta.New(username, password)
	fmt.Printf("Loading...")
	if err := insta.Login(); err != nil {
		panic(err)
	}
	return insta
}

// SearchUserHandler - Search user handler
func SearchUserHandler(query string, insta *goinsta.Instagram, simplifiedUsers *fetcherResponse.ByFollowersCount) {
	users, searchErr := insta.SearchUsername(query)
	if searchErr != nil {
		panic(searchErr)
	}
	*simplifiedUsers = fetcherResponse.ByFollowersCount{}
	for _, user := range users.Users {
		*simplifiedUsers = append(*simplifiedUsers, fetcherResponse.SimplifiedUser{
			Name:      user.FullName,
			Username:  user.Username,
			Followers: user.FollowerCount,
			ID:        user.Pk,
		})
	}
	for i, user := range *simplifiedUsers {
		fmt.Printf("[%d] %s - %s - %d followers\n", i+1, user.Username, user.Name, user.Followers)
	}
}

func collectPhotoURL(insta *goinsta.Instagram, targetUser fetcherResponse.SimplifiedUser) []string {
	const limit float64 = 300
	photoURLs := []string{}

	feeds, feedErr := insta.UserFeed(targetUser.ID, "", "")
	for len(photoURLs) <= int(limit) {
		fmt.Printf("Collecting...%.1f%%\n", (float64(len(photoURLs))/limit)*100)
		if feedErr != nil {
			panic(feedErr)
		}
		for _, item := range feeds.Items {
			if item.MediaType == 1 {
				photoURLs = append(photoURLs, item.ImageVersions2.Candidates[0].URL)
			}
		}
		if !feeds.MoreAvailable {
			break
		}
		feeds, _ = insta.UserFeed(targetUser.ID, feeds.NextMaxID, "")
	}
	return photoURLs
}

func downloadPictures(photoUrls []string, numberOfWorkers int, dirPath string) {
	fmt.Printf("100%% collected! Starting to download!\n")
	queue := make(chan string, len(photoUrls))
	resultQueue := make(chan int, len(photoUrls))
	for i := 0; i < numberOfWorkers; i++ {
		go startDownloadWorkers(dirPath, queue, resultQueue)
	}
	for _, url := range photoUrls {
		queue <- url
	}
	close(queue)
	for i := 0; i < len(photoUrls); i++ {
		progress := (float64(i) / float64(len(photoUrls))) * 100
		fmt.Printf("Downloaded %.1f%%\n", progress)
		<-resultQueue
	}
	fmt.Println("Download completed!")
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func startDownloadWorkers(dirPath string, urlChannel chan string, resultQueue chan int) {
	for url := range urlChannel {
		var imageType string
		if strings.Contains(url, ".png") {
			imageType = ".png"
		} else {
			imageType = ".jpg"
		}
		resp, httpErr := http.Get(url)
		if httpErr != nil {
			fmt.Println("HTTP Error: " + url)
			continue
		}
		defer resp.Body.Close()
		img, _, decodeErr := image.Decode(resp.Body)
		if decodeErr != nil {
			fmt.Println("Decode Error: " + url)
			continue
		}

		bounds := img.Bounds()
		if bounds.Size().X > 300 && bounds.Size().Y > 300 {
			imgHash := getMD5Hash(url)
			out, createFileErr := os.Create(dirPath + "/" + imgHash + imageType)
			if createFileErr != nil {
				fmt.Println("Create File Error: " + url)
				continue
			}
			defer out.Close()
			if imageType == ".png" {
				png.Encode(out, img)
			} else {
				jpeg.Encode(out, img, nil)
			}
		}
		resultQueue <- 1
	}
}

// DownloadPictureHandler - Download picture handler
func DownloadPictureHandler(insta *goinsta.Instagram, targetUser fetcherResponse.SimplifiedUser, numberOfWorkers int, baseDir string) {
	photoUrls := collectPhotoURL(insta, targetUser)
	dirPath := fmt.Sprintf("%v/%v", baseDir, targetUser.Username)
	os.MkdirAll(dirPath, 0755)
	downloadPictures(photoUrls, numberOfWorkers, dirPath)
}
