package model

type UserGroup struct {
	GroupId   int    `json:"id"`
	GroupName string `json:"name"`
	GroupDesc string `json:"description"`
	Type      string `json:"type"`
	Thumbnail string `json:"thumbnailurl"`
	IsActive  bool   `json:"is_active"`
}
