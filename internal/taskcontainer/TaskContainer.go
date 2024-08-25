package taskcontainer

type TaskContainer struct {
	ContainerId   string `json:"id"`
	ContainerName string `json:"name"`
	ContainerDesc string `json:"description"`
	IsActive      bool   `json:"is_active"`
	UsergroupId   int    `json:"usergroup_id"`
}
