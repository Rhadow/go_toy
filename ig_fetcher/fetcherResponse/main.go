package fetcherResponse

// SimplifiedUser - Simplified twitter user data structure
type SimplifiedUser struct {
	ID        int64
	Name      string
	Username  string
	Followers int
}

// ByFollowersCount - Sort user by Followers count
type ByFollowersCount []SimplifiedUser

func (s ByFollowersCount) Len() int {
	return len(s)
}
func (s ByFollowersCount) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByFollowersCount) Less(i, j int) bool {
	return s[i].Followers > s[j].Followers
}
