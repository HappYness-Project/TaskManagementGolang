package taskcontainer

type CreateContainerDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	UserGroupId int    `json:"usergroup_id"`
}
