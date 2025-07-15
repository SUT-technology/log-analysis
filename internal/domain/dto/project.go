package dto


import "github.com/google/uuid"

type ProjectSummery struct {
	ProjectID      uuid.UUID         `json:"project_id"`
	ProjectName    string            `json:"project_name"`
	SearchableKeys []string 		 `json:"searchable_keys"`
}

type ProjectsListRespone struct {
	UserID 	 uuid.UUID			`json:"user_id"`
	Projects []ProjectSummery	`json:"projects"`
}