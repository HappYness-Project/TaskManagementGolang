package usergroup

type UserGroup struct {
	GroupId   string `json:"id"`
	GroupName string `json:"name"`
	GroupDesc string `json:"description"`
	Type      string `json:"type"`
	Thumbnail string `json:"thumbnailurl"`
	IsActive  bool   `json:"is_active"`
}
