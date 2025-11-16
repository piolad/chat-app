package service

import (
	"context"

	"message-data-centre/server/data"
)

type MessageStore interface {
	Save(ctx context.Context, m *data.Message) error
	FetchByConversation(ctx context.Context, conversationID string, start, count int32) ([]*data.Message, int32, error)
}

type ConversationStore interface {
	GetOrCreate(ctx context.Context, sender, receiver string) (string, error)
	FetchByMember(ctx context.Context, member string, start, count int32) ([]*data.Conversation, int32, error)
}

type MessageService struct {
	msgStore  MessageStore
	convStore ConversationStore
}

func NewMessageService(m MessageStore, c ConversationStore) *MessageService {
	return &MessageService{msgStore: m, convStore: c}
}

func (s *MessageService) SendMessage(
	ctx context.Context, msg *data.Message, receiver string,
) error {
	convID, err := s.convStore.GetOrCreate(ctx, msg.Sender, receiver)
	if err != nil {
		return err
	}
	msg.ConversationID = convID
	return s.msgStore.Save(ctx, msg)
}

func (s *MessageService) FetchLastMessages(
	ctx context.Context, sender, receiver string, start, count int32,
) ([]*data.Message, bool, error) {
	convID, err := s.convStore.GetOrCreate(ctx, sender, receiver)
	if err != nil {
		return nil, false, err
	}

	msgs, total, err := s.msgStore.FetchByConversation(ctx, convID, start, count)
	if err != nil {
		return nil, false, err
	}

	hasMore := (start + int32(len(msgs))) < total
	return msgs, hasMore, nil
}

func (s *MessageService) FetchLastConversations(
	ctx context.Context, member string, start, count int32,
) ([]*data.Conversation, bool, error) {
	convs, total, err := s.convStore.FetchByMember(ctx, member, start, count)
	if err != nil {
		return nil, false, err
	}
	hasMore := (start + int32(len(convs))) < total
	return convs, hasMore, nil
}
