package model

import (
	"errors"
	"time"
)

type ChatParticipant struct {
	Id       string    `json:"id"`
	ChatId   string    `json:"chat_id"`
	UserId   int       `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
	Role     string    `json:"role"`
	Status   string    `json:"status"`
}

const (
	RoleAdmin  = "admin"
	RoleMember = "member"
)
const (
	StatusActive  = "active"
	StatusLeft    = "left"
	StatusBanned  = "banned"
	StatusMuted   = "muted"
	StatusPending = "pending"
)

func NewChatParticipant(chatId string, userId int, role string) (*ChatParticipant, error) {
	if err := validateRole(role); err != nil {
		return nil, err
	}

	participant := &ChatParticipant{
		ChatId:   chatId,
		UserId:   userId,
		JoinedAt: time.Now(),
		Role:     role,
		Status:   StatusActive,
	}

	return participant, nil
}

func (cp *ChatParticipant) UpdateRole(role string) error {
	if err := validateRole(role); err != nil {
		return err
	}
	cp.Role = role
	return nil
}

func (cp *ChatParticipant) UpdateStatus(status string) error {
	if err := validateStatus(status); err != nil {
		return err
	}
	cp.Status = status
	return nil
}
func validateRole(role string) error {
	validRoles := []string{RoleAdmin, RoleMember}
	for _, validRole := range validRoles {
		if role == validRole {
			return nil
		}
	}
	return errors.New("invalid role")
}

func validateStatus(status string) error {
	validStatuses := []string{StatusActive, StatusLeft, StatusBanned, StatusMuted, StatusPending}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}
	return errors.New("invalid status")
}
