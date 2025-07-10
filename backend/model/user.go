package model

type User struct {
    ID       string   `json:"id"`       // UUID
    Username string   `json:"username"`
    Password string   `json:"password"` // hashed password
    ProjectIDs []string `json:"projectIds"` // slice of project IDs
}