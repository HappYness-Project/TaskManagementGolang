package taskcontainer

type TaskContainer struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Type           string `json:"type"`
	IsActive       bool   `json:"is_active"`
	activity_level int    `json:"activity_level"`
	UsergroupId    int    `json:"usergroup_id"`
}
