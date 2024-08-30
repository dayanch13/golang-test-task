package models

var Users = map[string]string{"user1": "password1", "user2": "password2"}
var Notes = make(map[string][]string)

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserNoteReq struct {
	Content string `json:"content"`
}
