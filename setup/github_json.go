package setup

import "time"

type GithubRelease struct {
	Name       string    `json:"name"`
	TagName    string    `json:"tag_name"`
	ID         int       `json:"id"`
	Prerelease bool      `json:"prerelease"`
	Created    time.Time `json:"created_at"`
	Published  time.Time `json:"published_at"`
	Zipball    string    `json:"zipball_url"`
}
