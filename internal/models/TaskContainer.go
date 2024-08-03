package models

type TaskContainer struct {
	ContainerId   string `json:id`
	ContainerName string `json:container_name`
	ContainerDesc string `json:container_desc`
}
