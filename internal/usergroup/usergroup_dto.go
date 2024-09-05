package usergroup

type CreateUserGroupDto struct {
	GroupName string `json:"name"`
	GroupDesc string `json:"description"`
	GroupType string `json:"type"`
}

// type GetUserGroupDetailsDto struct {
// 	GroupName string `json:"name"`
// 	GroupDesc string `json:"description"`
// 	GroupType string `json:"type"`
// 	Users []User
// }
