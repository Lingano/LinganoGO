package model

type Reading struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	User     *User  `json:"user"`
	Finished bool   `json:"finished"`
}