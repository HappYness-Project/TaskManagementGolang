package repository

const (
	GetAllChats                = `SELECT * FROM public.chat`
	GetChatsByUserGroupIdQuery = `
		SELECT id, type, usergroup_id, container_id, created_at
		FROM chat
		WHERE usergroup_id = $1
		ORDER BY created_at DESC
	`

	GetChatByIdQuery = `
		SELECT id, type, usergroup_id, container_id, created_at
		FROM chat
		WHERE id = $1
	`

	CreateChatQuery = `
		INSERT INTO chat (id, type, usergroup_id, container_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	DeleteChatQuery = `
		DELETE FROM chat
		WHERE id = $1
	`
)
