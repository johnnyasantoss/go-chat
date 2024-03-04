package server

import (
	"context"

	"johnnyasantos.com/chat/shared"
)

type ChatContext struct {
	Room *shared.Room
}

type contextKey string

var chatContextKey = contextKey("chatContext")

func NewChatContext(ctx context.Context, room *shared.Room) context.Context {
	return context.WithValue(ctx, chatContextKey, &ChatContext{room})
}

func ChatContextFromContext(ctx context.Context) (*ChatContext, bool) {
	chatCtx, ok := ctx.Value(chatContextKey).(*ChatContext)
	return chatCtx, ok
}
