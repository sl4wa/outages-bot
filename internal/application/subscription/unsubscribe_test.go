package subscription

import (
	"errors"
	"outages-bot/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsubscribe_ExistingUser(t *testing.T) {
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[12345] = &domain.User{ID: 12345, Address: addr}

	svc := NewUnsubscribeService(repo)
	result := svc.Handle(12345)

	assert.NoError(t, result.Err)
	assert.Equal(t, "Ви успішно відписалися від сповіщень про відключення електроенергії.", result.Message)
	assert.Equal(t, []int64{12345}, repo.removed)
}

func TestUnsubscribe_NoSubscription(t *testing.T) {
	repo := newMockUserRepo()

	svc := NewUnsubscribeService(repo)
	result := svc.Handle(12345)

	assert.NoError(t, result.Err)
	assert.Equal(t, "Ви не маєте активної підписки.", result.Message)
}

func TestUnsubscribe_RepoError(t *testing.T) {
	repo := newMockUserRepo()
	repo.removeErr = errors.New("disk error")

	svc := NewUnsubscribeService(repo)
	result := svc.Handle(12345)

	assert.EqualError(t, result.Err, "disk error")
	assert.Equal(t, "Сталася помилка. Спробуйте пізніше.", result.Message)
}
