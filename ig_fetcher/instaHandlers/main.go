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

// DownloadPictureHandler - Download picture handler
func DownloadPictureHandler(targetIndex int, insta *goinsta.Instagram, simplifiedUsers *fetcherResponse.ByFollowersCount) {
	realIndex := targetIndex - 1
	if len(*simplifiedUsers) <= realIndex || realIndex < 0 {
		fmt.Printf("%d out of range\n", targetIndex)
		return
	}
	targetUser := (*simplifiedUsers)[realIndex]
	feeds, _ := insta.UserFeed(targetUser.ID, "", "")
	// feeds.Items[0].ImageVersions2.Candidates[0].URL
	fmt.Print(feeds)
}
