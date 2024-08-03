package models

type TaskContainer struct {
	ContainerId   string `json:id`
	ContainerName string `json:name`
	ContainerDesc string `json:description`
}
