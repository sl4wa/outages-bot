package subscription

import (
	"errors"
	"outages-bot/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowSubscription_NewUser(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewShowSubscriptionService(repo)
	msg := svc.Handle(12345)
	assert.Equal(t, "Будь ласка, введіть назву вулиці:", msg)
}

func TestShowSubscription_ExistingUser(t *testing.T) {
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[12345] = &domain.User{ID: 12345, Address: addr}
	svc := NewShowSubscriptionService(repo)
	msg := svc.Handle(12345)
	assert.Contains(t, msg, "Ваша поточна підписка:")
	assert.Contains(t, msg, "Вулиця: Стрийська")
	assert.Contains(t, msg, "Будинок: 10")
	assert.Contains(t, msg, "Будь ласка, введіть нову назву вулиці для оновлення підписки:")
}

func TestShowSubscription_DifferentChatIDs(t *testing.T) {
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Наукова", "5")
	repo.users[111] = &domain.User{ID: 111, Address: addr}
	svc := NewShowSubscriptionService(repo)

	msg1 := svc.Handle(111)
	assert.Contains(t, msg1, "Наукова")

	msg2 := svc.Handle(222)
	assert.Equal(t, "Будь ласка, введіть назву вулиці:", msg2)
}

func TestShowSubscription_DifferentStreetAndBuilding(t *testing.T) {
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(2, "Молдавська", "25-А")
	repo.users[333] = &domain.User{ID: 333, Address: addr}
	svc := NewShowSubscriptionService(repo)
	msg := svc.Handle(333)
	assert.Contains(t, msg, "Вулиця: Молдавська")
	assert.Contains(t, msg, "Будинок: 25-А")
}

func TestShowSubscription_CyrillicLabelsPresent(t *testing.T) {
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[12345] = &domain.User{ID: 12345, Address: addr}
	svc := NewShowSubscriptionService(repo)
	msg := svc.Handle(12345)
	assert.Contains(t, msg, "Вулиця:")
	assert.Contains(t, msg, "Будинок:")
}

func TestShowSubscription_CorruptedDataFallsBack(t *testing.T) {
	repo := newMockUserRepo()
	repo.findErr = errors.New("corrupted data")
	svc := NewShowSubscriptionService(repo)
	msg := svc.Handle(12345)
	assert.Equal(t, "Будь ласка, введіть назву вулиці:", msg)
}

func TestShowCurrent_ExistingUser(t *testing.T) {
	repo := newMockUserRepo()
	addr, _ := domain.NewUserAddress(1, "Стрийська", "10")
	repo.users[12345] = &domain.User{ID: 12345, Address: addr}
	svc := NewShowSubscriptionService(repo)
	msg, err := svc.ShowCurrent(12345)
	assert.NoError(t, err)
	assert.Equal(t, "Ваша поточна підписка:\nВулиця: Стрийська\nБудинок: 10", msg)
	assert.NotContains(t, msg, "введіть нову назву вулиці")
}

func TestShowCurrent_NoUser(t *testing.T) {
	repo := newMockUserRepo()
	svc := NewShowSubscriptionService(repo)
	msg, err := svc.ShowCurrent(12345)
	assert.NoError(t, err)
	assert.Equal(t, "Ви не маєте активної підписки.", msg)
}

func TestShowCurrent_RepoError(t *testing.T) {
	repo := newMockUserRepo()
	repo.findErr = errors.New("disk error")
	svc := NewShowSubscriptionService(repo)
	_, err := svc.ShowCurrent(12345)
	assert.Error(t, err)
}
