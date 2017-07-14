package instaHandlers

import (
	"fmt"

	"github.com/ahmdrz/goinsta"
	"github.com/rhadow/go_toy/ig_fetcher/fetcherResponse"
)

// LoginInstagram - Log in instagram
func LoginInstagram(username, password, loadingText string) *goinsta.Instagram {
	insta := goinsta.New(username, password)
	fmt.Printf(loadingText)
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

func collectPhotoURL(targetIndex int, insta *goinsta.Instagram, simplifiedUsers *fetcherResponse.ByFollowersCount) []string {
	const limit float64 = 300
	photoURLs := []string{}
	realIndex := targetIndex - 1
	if len(*simplifiedUsers) <= realIndex || realIndex < 0 {
		fmt.Printf("%d out of range\n", targetIndex)
		return []string{}
	}
	targetUser := (*simplifiedUsers)[realIndex]
	feeds, feedErr := insta.UserFeed(targetUser.ID, "", "")
	for feeds.MoreAvailable && len(photoURLs) <= int(limit) {
		fmt.Printf("Collecting...%.1f%%\n", (float64(len(photoURLs))/limit)*100)
		if feedErr != nil {
			panic(feedErr)
		}
		for _, item := range feeds.Items {
			if item.MediaType == 1 {
				photoURLs = append(photoURLs, item.ImageVersions2.Candidates[0].URL)
			}
		}
		feeds, _ = insta.UserFeed(targetUser.ID, feeds.NextMaxID, "")
	}
	return photoURLs
}

func downloadPictures(photoUrls []string, numberOfWorkers int) {
	fmt.Printf("100%% collected! Starting to download!\n")
	completeChannel := make(chan int)
	completeCount := 0
	partitionLength := len(photoUrls) / numberOfWorkers
	for i := 0; i < numberOfWorkers; i++ {
		start := i * partitionLength
		end := start + partitionLength
		if i < numberOfWorkers-1 {
			go startDownloadWorkers(photoUrls[start:end], completeChannel)
		} else {
			go startDownloadWorkers(photoUrls[start:], completeChannel)
		}
	}
	for range completeChannel {
		completeCount++
		fmt.Println(completeCount)
		if completeCount == numberOfWorkers {
			close(completeChannel)
		}
	}
}

// TODO: Complete download task
func startDownloadWorkers(photoUrls []string, completeChannel chan int) {
	for i, url := range photoUrls {
		fmt.Println(i, url)
	}
	completeChannel <- 1
}

// DownloadPictureHandler - Download picture handler
func DownloadPictureHandler(targetIndex int, insta *goinsta.Instagram, simplifiedUsers *fetcherResponse.ByFollowersCount, numberOfWorkers int) {
	photoUrls := collectPhotoURL(targetIndex, insta, simplifiedUsers)
	downloadPictures(photoUrls, numberOfWorkers)
	fmt.Println("Download completed!")
}
