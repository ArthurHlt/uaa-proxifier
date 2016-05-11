package main

type UserFromUaa struct {
	Email      string `json:"email"`
	FamilyName interface{} `json:"family_name"`
	GivenName  interface{} `json:"given_name"`
	Name       string `json:"name"`
	UserID     string `json:"user_id"`
	UserName   string `json:"user_name"`
}
type GitLabUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}