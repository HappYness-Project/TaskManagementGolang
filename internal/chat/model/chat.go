package model

import (
	"errors"
	"time"
)

type Chat struct {
	Id          string    `json:"id"`
	Type        string    `json:"type"`
	UserGroupId *int      `json:"usergroup_id,omitempty"`
	ContainerId *string   `json:"container_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

const (
	ChatTypePrivate   = "private"
	ChatTypeGroup     = "group"
	ChatTypeContainer = "container"
)

func NewChat(chatType string, userGroupId *int, containerId *string) (*Chat, error) {
	if err := validateChatType(chatType); err != nil {
		return nil, err
	}

	if err := validateChatData(chatType, userGroupId, containerId); err != nil {
		return nil, err
	}

	chat := &Chat{
		Type:        chatType,
		UserGroupId: userGroupId,
		ContainerId: containerId,
		CreatedAt:   time.Now(),
	}

	return chat, nil
}

func NewGroupChat(userGroupId int) (*Chat, error) {
	return NewChat(ChatTypeGroup, &userGroupId, nil)
}

func NewContainerChat(userGroupId int, containerId string) (*Chat, error) {
	return NewChat(ChatTypeContainer, &userGroupId, &containerId)
}

func NewPrivateChat() (*Chat, error) {
	return NewChat(ChatTypePrivate, nil, nil)
}

func (c *Chat) IsGroupChat() bool {
	return c.Type == ChatTypeGroup
}

func (c *Chat) IsContainerChat() bool {
	return c.Type == ChatTypeContainer
}

func (c *Chat) IsPrivateChat() bool {
	return c.Type == ChatTypePrivate
}

func validateChatType(chatType string) error {
	validTypes := []string{ChatTypePrivate, ChatTypeGroup, ChatTypeContainer}
	for _, validType := range validTypes {
		if chatType == validType {
			return nil
		}
	}
	return errors.New("invalid chat type")
}

func validateChatData(chatType string, userGroupId *int, containerId *string) error {
	switch chatType {
	case ChatTypePrivate:
		if userGroupId != nil || containerId != nil {
			return errors.New("private chats should not have user group or container")
		}
	case ChatTypeGroup:
		if userGroupId == nil {
			return errors.New("group chats must have a user group")
		}
		if containerId != nil {
			return errors.New("group chats should not have container")
		}
	case ChatTypeContainer:
		if userGroupId == nil {
			return errors.New("container chats must have a user group")
		}
		if containerId == nil {
			return errors.New("container chats must have a container")
		}
	}
	return nil
}
