package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/happYness-Project/taskManagementGolang/internal/chat/model"
)

const dbTimeout = time.Second * 5

type ChatRepository interface {
	GetAllChats() ([]model.Chat, error)
	GetChatsByUserGroupId(ctx context.Context, userGroupId int) ([]model.Chat, error)
	GetChatById(ctx context.Context, chatId string) (*model.Chat, error)
	CreateChat(ctx context.Context, chat *model.Chat) error
	DeleteChat(ctx context.Context, chatId string) error
}

type ChatRepo struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepo {
	return &ChatRepo{
		DB: db,
	}
}

func (r *ChatRepo) GetAllChats() ([]model.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := r.DB.QueryContext(ctx, GetAllChats)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chats []model.Chat
	for rows.Next() {
		chat, err := scanRowsIntoChat(rows)
		if err != nil {
			return nil, err
		}
		chats = append(chats, *chat)
	}
	return chats, nil
}

func (r *ChatRepo) GetChatsByUserGroupId(ctx context.Context, userGroupId int) ([]model.Chat, error) {
	rows, err := r.DB.QueryContext(ctx, GetChatsByUserGroupIdQuery, userGroupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []model.Chat
	for rows.Next() {
		var chat model.Chat
		var userGroupIdPtr *int
		var containerIdPtr *string

		err := rows.Scan(
			&chat.Id,
			&chat.Type,
			&userGroupIdPtr,
			&containerIdPtr,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		chat.UserGroupId = userGroupIdPtr
		chat.ContainerId = containerIdPtr
		chats = append(chats, chat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *ChatRepo) GetChatById(ctx context.Context, chatId string) (*model.Chat, error) {
	var chat model.Chat
	var userGroupIdPtr *int
	var containerIdPtr *string

	err := r.DB.QueryRowContext(ctx, GetChatByIdQuery, chatId).Scan(
		&chat.Id,
		&chat.Type,
		&userGroupIdPtr,
		&containerIdPtr,
		&chat.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	chat.UserGroupId = userGroupIdPtr
	chat.ContainerId = containerIdPtr

	return &chat, nil
}

func (r *ChatRepo) CreateChat(ctx context.Context, chat *model.Chat) error {
	_, err := r.DB.ExecContext(ctx, CreateChatQuery,
		chat.Id,
		chat.Type,
		chat.UserGroupId,
		chat.ContainerId,
		chat.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ChatRepo) DeleteChat(ctx context.Context, chatId string) error {
	result, err := r.DB.ExecContext(ctx, DeleteChatQuery, chatId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		errors.New("No chat found to delete")
	}

	return nil
}
func scanRowsIntoChat(rows *sql.Rows) (*model.Chat, error) {
	chat := new(model.Chat)
	err := rows.Scan(
		&chat.Id,
		&chat.Type,
		&chat.UserGroupId,
		&chat.ContainerId,
		&chat.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return chat, nil
}
