package model

// GithubUser is model for github user
type GithubUser struct {
	Login string `json:"login,omitempty"`
	Id    int64  `json:"id,omitempty"`
}
