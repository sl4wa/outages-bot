package telegram

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserNotifierSendsInChatIDOrder(t *testing.T) {
	sender := &fakeSender{}
	store := fakeStore{ids: []int64{111, 222}}

	notified, err := UserNotifier{Sender: sender, Subscribers: store}.Notify(context.Background(), "msg")

	require.NoError(t, err)
	assert.True(t, notified)
	assert.Equal(t, []int64{111, 222}, sender.chatIDs)
	assert.Equal(t, []string{"msg", "msg"}, sender.messages)
}

func TestUserNotifierPerUserFailureDoesNotStopOthers(t *testing.T) {
	sender := &fakeSender{failFor: map[int64]error{111: errors.New("forbidden")}}
	store := fakeStore{ids: []int64{111, 222}}

	notified, err := UserNotifier{Sender: sender, Subscribers: store}.Notify(context.Background(), "msg")

	require.NoError(t, err)
	assert.True(t, notified)
	assert.Equal(t, []int64{111, 222}, sender.chatIDs)
}

func TestUserNotifierEmptyChatIDs(t *testing.T) {
	sender := &fakeSender{}

	notified, err := UserNotifier{Sender: sender, Subscribers: fakeStore{}}.Notify(context.Background(), "msg")

	require.NoError(t, err)
	assert.False(t, notified)
	assert.Empty(t, sender.chatIDs)
}

func TestUserNotifierPropagatesStoreError(t *testing.T) {
	sender := &fakeSender{}
	storeErr := errors.New("boom")

	notified, err := UserNotifier{Sender: sender, Subscribers: fakeStore{err: storeErr}}.Notify(context.Background(), "msg")

	require.ErrorIs(t, err, storeErr)
	assert.False(t, notified)
	assert.Empty(t, sender.chatIDs)
}

type fakeStore struct {
	ids []int64
	err error
}

func (s fakeStore) ChatIDs() ([]int64, error) { return s.ids, s.err }

type fakeSender struct {
	chatIDs  []int64
	messages []string
	failFor  map[int64]error
}

func (s *fakeSender) SendHTML(ctx context.Context, chatID int64, text string) error {
	_ = ctx
	s.chatIDs = append(s.chatIDs, chatID)
	s.messages = append(s.messages, text)
	if s.failFor != nil {
		return s.failFor[chatID]
	}
	return nil
}
