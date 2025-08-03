package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChat(t *testing.T) {
	t.Run("when creating new group chat with valid data, Then return chat with correct fields", func(t *testing.T) {
		// Given
		chatType := ChatTypeGroup
		userGroupId := 1

		// When
		chat, err := NewChat(chatType, &userGroupId, nil)

		// Then
		require.NoError(t, err)
		require.NotNil(t, chat)
		assert.Equal(t, chatType, chat.Type)
		assert.Equal(t, &userGroupId, chat.UserGroupId)
		assert.Nil(t, chat.ContainerId)
		assert.NotZero(t, chat.CreatedAt)
	})

	t.Run("when creating new container chat with valid data, Then return chat with correct fields", func(t *testing.T) {
		// Given
		chatType := ChatTypeContainer
		userGroupId := 1
		containerId := "container-123"

		// When
		chat, err := NewChat(chatType, &userGroupId, &containerId)

		// Then
		require.NoError(t, err)
		require.NotNil(t, chat)
		assert.Equal(t, chatType, chat.Type)
		assert.Equal(t, &userGroupId, chat.UserGroupId)
		assert.Equal(t, &containerId, chat.ContainerId)
		assert.NotZero(t, chat.CreatedAt)
	})

	t.Run("when creating new private chat with valid data, Then return chat with correct fields", func(t *testing.T) {
		// Given
		chatType := ChatTypePrivate

		// When
		chat, err := NewChat(chatType, nil, nil)

		// Then
		require.NoError(t, err)
		require.NotNil(t, chat)
		assert.Equal(t, chatType, chat.Type)
		assert.Nil(t, chat.UserGroupId)
		assert.Nil(t, chat.ContainerId)
		assert.NotZero(t, chat.CreatedAt)
	})

	t.Run("when creating chat with invalid type, Then return error", func(t *testing.T) {
		// Given
		invalidType := "invalid"

		// When
		chat, err := NewChat(invalidType, nil, nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid chat type")
		assert.Nil(t, chat)
	})

	t.Run("when creating group chat without user group, Then return error", func(t *testing.T) {
		// Given
		chatType := ChatTypeGroup

		// When
		chat, err := NewChat(chatType, nil, nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group chats must have a user group")
		assert.Nil(t, chat)
	})

	t.Run("when creating group chat with container, Then return error", func(t *testing.T) {
		// Given
		chatType := ChatTypeGroup
		userGroupId := 1
		containerId := "container-123"

		// When
		chat, err := NewChat(chatType, &userGroupId, &containerId)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group chats should not have container")
		assert.Nil(t, chat)
	})

	t.Run("when creating container chat without user group, Then return error", func(t *testing.T) {
		// Given
		chatType := ChatTypeContainer
		containerId := "container-123"

		// When
		chat, err := NewChat(chatType, nil, &containerId)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "container chats must have a user group")
		assert.Nil(t, chat)
	})

	t.Run("when creating container chat without container, Then return error", func(t *testing.T) {
		// Given
		chatType := ChatTypeContainer
		userGroupId := 1

		// When
		chat, err := NewChat(chatType, &userGroupId, nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "container chats must have a container")
		assert.Nil(t, chat)
	})

	t.Run("when creating private chat with user group, Then return error", func(t *testing.T) {
		// Given
		chatType := ChatTypePrivate
		userGroupId := 1

		// When
		chat, err := NewChat(chatType, &userGroupId, nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "private chats should not have user group or container")
		assert.Nil(t, chat)
	})

	t.Run("when creating private chat with container, Then return error", func(t *testing.T) {
		// Given
		chatType := ChatTypePrivate
		containerId := "container-123"

		// When
		chat, err := NewChat(chatType, nil, &containerId)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "private chats should not have user group or container")
		assert.Nil(t, chat)
	})

	t.Run("when creating multiple chats, Then each chat has unique timestamps", func(t *testing.T) {
		// Given
		userGroupId1 := 1
		userGroupId2 := 2

		// When
		chat1, err1 := NewChat(ChatTypeGroup, &userGroupId1, nil)
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
		chat2, err2 := NewChat(ChatTypeGroup, &userGroupId2, nil)

		// Then
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.True(t, chat2.CreatedAt.After(chat1.CreatedAt))
	})
}

func TestNewGroupChat(t *testing.T) {
	t.Run("when creating new group chat, Then return valid group chat", func(t *testing.T) {
		// Given
		userGroupId := 1

		// When
		chat, err := NewGroupChat(userGroupId)

		// Then
		require.NoError(t, err)
		require.NotNil(t, chat)
		assert.Equal(t, ChatTypeGroup, chat.Type)
		assert.Equal(t, &userGroupId, chat.UserGroupId)
		assert.Nil(t, chat.ContainerId)
		assert.NotZero(t, chat.CreatedAt)
	})

	t.Run("when creating multiple group chats, Then each chat has unique timestamps", func(t *testing.T) {
		// Given
		userGroupId1 := 1
		userGroupId2 := 2

		// When
		chat1, err1 := NewGroupChat(userGroupId1)
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
		chat2, err2 := NewGroupChat(userGroupId2)

		// Then
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.True(t, chat2.CreatedAt.After(chat1.CreatedAt))
	})
}

func TestNewContainerChat(t *testing.T) {
	t.Run("when creating new container chat, Then return valid container chat", func(t *testing.T) {
		// Given
		userGroupId := 1
		containerId := "container-123"

		// When
		chat, err := NewContainerChat(userGroupId, containerId)

		// Then
		require.NoError(t, err)
		require.NotNil(t, chat)
		assert.Equal(t, ChatTypeContainer, chat.Type)
		assert.Equal(t, &userGroupId, chat.UserGroupId)
		assert.Equal(t, &containerId, chat.ContainerId)
		assert.NotZero(t, chat.CreatedAt)
	})

	t.Run("when creating multiple container chats, Then each chat has unique timestamps", func(t *testing.T) {
		// Given
		userGroupId := 1
		containerId1 := "container-1"
		containerId2 := "container-2"

		// When
		chat1, err1 := NewContainerChat(userGroupId, containerId1)
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
		chat2, err2 := NewContainerChat(userGroupId, containerId2)

		// Then
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.True(t, chat2.CreatedAt.After(chat1.CreatedAt))
	})
}

func TestNewPrivateChat(t *testing.T) {
	t.Run("when creating new private chat, Then return valid private chat", func(t *testing.T) {
		// When
		chat, err := NewPrivateChat()

		// Then
		require.NoError(t, err)
		require.NotNil(t, chat)
		assert.Equal(t, ChatTypePrivate, chat.Type)
		assert.Nil(t, chat.UserGroupId)
		assert.Nil(t, chat.ContainerId)
		assert.NotZero(t, chat.CreatedAt)
	})

	t.Run("when creating multiple private chats, Then each chat has unique timestamps", func(t *testing.T) {
		// When
		chat1, err1 := NewPrivateChat()
		time.Sleep(1 * time.Millisecond) // Ensure different timestamps
		chat2, err2 := NewPrivateChat()

		// Then
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.True(t, chat2.CreatedAt.After(chat1.CreatedAt))
	})
}

func TestChatTypeChecking(t *testing.T) {
	t.Run("when checking group chat type, Then return correct boolean", func(t *testing.T) {
		// Given
		userGroupId := 1
		groupChat, _ := NewGroupChat(userGroupId)
		containerChat, _ := NewContainerChat(userGroupId, "container-123")
		privateChat, _ := NewPrivateChat()

		// Then
		assert.True(t, groupChat.IsGroupChat())
		assert.False(t, containerChat.IsGroupChat())
		assert.False(t, privateChat.IsGroupChat())
	})

	t.Run("when checking container chat type, Then return correct boolean", func(t *testing.T) {
		// Given
		userGroupId := 1
		groupChat, _ := NewGroupChat(userGroupId)
		containerChat, _ := NewContainerChat(userGroupId, "container-123")
		privateChat, _ := NewPrivateChat()

		// Then
		assert.False(t, groupChat.IsContainerChat())
		assert.True(t, containerChat.IsContainerChat())
		assert.False(t, privateChat.IsContainerChat())
	})

	t.Run("when checking private chat type, Then return correct boolean", func(t *testing.T) {
		// Given
		userGroupId := 1
		groupChat, _ := NewGroupChat(userGroupId)
		containerChat, _ := NewContainerChat(userGroupId, "container-123")
		privateChat, _ := NewPrivateChat()

		// Then
		assert.False(t, groupChat.IsPrivateChat())
		assert.False(t, containerChat.IsPrivateChat())
		assert.True(t, privateChat.IsPrivateChat())
	})
}

func TestChatStruct(t *testing.T) {
	t.Run("when creating chat struct directly, Then all fields can be set", func(t *testing.T) {
		// Given
		now := time.Now()
		userGroupId := 1
		containerId := "container-123"
		chat := Chat{
			Id:          "chat-123",
			Type:        ChatTypeContainer,
			UserGroupId: &userGroupId,
			ContainerId: &containerId,
			CreatedAt:   now,
		}

		// Then
		assert.Equal(t, "chat-123", chat.Id)
		assert.Equal(t, ChatTypeContainer, chat.Type)
		assert.Equal(t, &userGroupId, chat.UserGroupId)
		assert.Equal(t, &containerId, chat.ContainerId)
		assert.Equal(t, now, chat.CreatedAt)
	})

	t.Run("when modifying chat fields, Then changes are reflected", func(t *testing.T) {
		// Given
		userGroupId := 1
		chat, _ := NewGroupChat(userGroupId)
		originalCreatedAt := chat.CreatedAt

		// When
		chat.Id = "modified-id"
		chat.Type = ChatTypePrivate
		chat.UserGroupId = nil
		chat.ContainerId = nil

		// Then
		assert.Equal(t, "modified-id", chat.Id)
		assert.Equal(t, ChatTypePrivate, chat.Type)
		assert.Nil(t, chat.UserGroupId)
		assert.Nil(t, chat.ContainerId)
		assert.Equal(t, originalCreatedAt, chat.CreatedAt) // Should not change
	})
}

func TestValidateChatType(t *testing.T) {
	t.Run("when validating valid chat types, Then return no error", func(t *testing.T) {
		validTypes := []string{ChatTypePrivate, ChatTypeGroup, ChatTypeContainer}

		for _, chatType := range validTypes {
			err := validateChatType(chatType)
			assert.NoError(t, err, "Chat type %s should be valid", chatType)
		}
	})

	t.Run("when validating invalid chat types, Then return error", func(t *testing.T) {
		invalidTypes := []string{"invalid", "public", "direct", "", "GROUP"}

		for _, chatType := range invalidTypes {
			err := validateChatType(chatType)
			assert.Error(t, err, "Chat type %s should be invalid", chatType)
			assert.Contains(t, err.Error(), "invalid chat type")
		}
	})
}

func TestValidateChatData(t *testing.T) {
	t.Run("when validating private chat data, Then return no error for valid data", func(t *testing.T) {
		err := validateChatData(ChatTypePrivate, nil, nil)
		assert.NoError(t, err)
	})

	t.Run("when validating private chat data with user group, Then return error", func(t *testing.T) {
		userGroupId := 1
		err := validateChatData(ChatTypePrivate, &userGroupId, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "private chats should not have user group or container")
	})

	t.Run("when validating private chat data with container, Then return error", func(t *testing.T) {
		containerId := "container-123"
		err := validateChatData(ChatTypePrivate, nil, &containerId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "private chats should not have user group or container")
	})

	t.Run("when validating group chat data, Then return no error for valid data", func(t *testing.T) {
		userGroupId := 1
		err := validateChatData(ChatTypeGroup, &userGroupId, nil)
		assert.NoError(t, err)
	})

	t.Run("when validating group chat data without user group, Then return error", func(t *testing.T) {
		err := validateChatData(ChatTypeGroup, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group chats must have a user group")
	})

	t.Run("when validating group chat data with container, Then return error", func(t *testing.T) {
		userGroupId := 1
		containerId := "container-123"
		err := validateChatData(ChatTypeGroup, &userGroupId, &containerId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "group chats should not have container")
	})

	t.Run("when validating container chat data, Then return no error for valid data", func(t *testing.T) {
		userGroupId := 1
		containerId := "container-123"
		err := validateChatData(ChatTypeContainer, &userGroupId, &containerId)
		assert.NoError(t, err)
	})

	t.Run("when validating container chat data without user group, Then return error", func(t *testing.T) {
		containerId := "container-123"
		err := validateChatData(ChatTypeContainer, nil, &containerId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "container chats must have a user group")
	})

	t.Run("when validating container chat data without container, Then return error", func(t *testing.T) {
		userGroupId := 1
		err := validateChatData(ChatTypeContainer, &userGroupId, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "container chats must have a container")
	})
}
