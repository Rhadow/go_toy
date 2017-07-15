package instaHandlers

import (
	"fmt"

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

func downloadPictures(photoUrls []string, numberOfWorkers int, baseDir string) {
	fmt.Printf("100%% collected! Starting to download!\n")
	queue := make(chan string, len(photoUrls))
	resultQueue := make(chan int, len(photoUrls))
	for i := 0; i < numberOfWorkers; i++ {
		go startDownloadWorkers(baseDir, queue, resultQueue)
	}
	for _, url := range photoUrls {
		queue <- url
	}
	close(queue)
	for i := 0; i < len(photoUrls); i++ {
		<-resultQueue
	}
	fmt.Println("Download completed!")
}

func startDownloadWorkers(baseDir string, urlChannel chan string, resultQueue chan int) {
	for url := range urlChannel {
		fmt.Println(url)
		resultQueue <- 1
	}
}

// DownloadPictureHandler - Download picture handler
func DownloadPictureHandler(insta *goinsta.Instagram, targetUser fetcherResponse.SimplifiedUser, numberOfWorkers int, baseDir string) {
	photoUrls := collectPhotoURL(insta, targetUser)
	fmt.Println(len(photoUrls))
	downloadPictures(photoUrls, numberOfWorkers, baseDir)
}
