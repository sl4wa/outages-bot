package subscription

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"outages-bot/internal/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testUserRepo struct {
	users     map[int64]*users.User
	findErr   error
	saveErr   error
	removeErr error
}

func newTestUserRepo() *testUserRepo {
	return &testUserRepo{users: make(map[int64]*users.User)}
}

func (r *testUserRepo) Find(chatID int64) (*users.User, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.users[chatID], nil
}

func (r *testUserRepo) Save(user *users.User) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.users[user.ID] = user
	return nil
}

func (r *testUserRepo) Remove(chatID int64) (bool, error) {
	if r.removeErr != nil {
		return false, r.removeErr
	}
	if _, ok := r.users[chatID]; !ok {
		return false, nil
	}
	delete(r.users, chatID)
	return true, nil
}

type testStreetRepo struct {
	streets []users.Street
}

func (r *testStreetRepo) GetAllStreets() []users.Street {
	return r.streets
}

func testStreets() []users.Street {
	return []users.Street{
		{ID: 1, Name: "Стрийська"},
		{ID: 2, Name: "Наукова"},
		{ID: 3, Name: "Стрілецька"},
	}
}

func newTestWorkflow(t *testing.T, repo *testUserRepo) (*Workflow, *testUserRepo) {
	t.Helper()
	if repo == nil {
		repo = newTestUserRepo()
	}
	wf := NewWorkflow(WorkflowConfig{
		UserRepo:   repo,
		StreetRepo: &testStreetRepo{streets: testStreets()},
	})
	return wf, repo
}

func addUser(t *testing.T, repo *testUserRepo, chatID int64, streetID int, streetName, building string) {
	t.Helper()
	addr, err := users.NewAddress(streetID, streetName, building)
	require.NoError(t, err)
	repo.users[chatID] = &users.User{ID: chatID, Address: addr}
}

func startSearch(t *testing.T, wf *Workflow, chatID int64) {
	t.Helper()
	response := wf.Handle(chatID, Command{Kind: CommandStart})
	require.Equal(t, messagePromptStreet, response.Text)
}

func selectStreet(t *testing.T, wf *Workflow, chatID int64, street string) {
	t.Helper()
	response := wf.Handle(chatID, Command{Kind: CommandText, Text: street})
	require.Equal(t, fmt.Sprintf(messagePromptBuilding, street), response.Text)
}

func TestServiceStartStopAndSubscription(t *testing.T) {
	wf, repo := newTestWorkflow(t, nil)

	response := wf.Handle(100, Command{Kind: CommandStart})
	assert.Equal(t, messagePromptStreet, response.Text)
	require.NotNil(t, wf.GetState(100))
	assert.Equal(t, StepSearchStreet, wf.GetState(100).Step)

	response = wf.Handle(100, Command{Kind: CommandStop})
	assert.Equal(t, messageNoSubscription, response.Text)
	assert.Nil(t, wf.GetState(100))

	addUser(t, repo, 100, 1, "Стрийська", "10")

	response = wf.Handle(100, Command{Kind: CommandSubscription})
	assert.Equal(t, "Ваша поточна підписка:\nВулиця: Стрийська\nБудинок: 10", response.Text)

	response = wf.Handle(100, Command{Kind: CommandStart})
	assert.Equal(t, "Ваша поточна підписка:\nВулиця: Стрийська\nБудинок: 10\n\nБудь ласка, введіть нову назву вулиці для оновлення підписки:", response.Text)

	response = wf.Handle(100, Command{Kind: CommandStop})
	assert.Equal(t, messageUnsubscribed, response.Text)
	assert.Nil(t, repo.users[100])
}

func TestServiceStreetSearch(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		wantText    string
		wantStep    StepKind
		wantStreet  int
		wantOptions []string
	}{
		{
			name:     "empty query",
			query:    "  ",
			wantText: messageEmptyStreetQuery,
			wantStep: StepSearchStreet,
		},
		{
			name:     "missing street",
			query:    "Невідома",
			wantText: messageStreetNotFound,
			wantStep: StepSearchStreet,
		},
		{
			name:       "single partial match",
			query:      "науков",
			wantText:   "Ви обрали вулицю: Наукова\nБудь ласка, введіть номер будинку:",
			wantStep:   StepSaveSubscription,
			wantStreet: 2,
		},
		{
			name:        "multiple matches",
			query:       "Стр",
			wantText:    messageStreetOptions,
			wantStep:    StepSearchStreet,
			wantOptions: []string{"Стрийська", "Стрілецька"},
		},
		{
			name:       "exact match",
			query:      "Наукова",
			wantText:   "Ви обрали вулицю: Наукова\nБудь ласка, введіть номер будинку:",
			wantStep:   StepSaveSubscription,
			wantStreet: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wf, _ := newTestWorkflow(t, nil)
			startSearch(t, wf, 100)

			response := wf.Handle(100, Command{Kind: CommandText, Text: tt.query})

			assert.Equal(t, tt.wantText, response.Text)
			assert.Equal(t, tt.wantOptions, response.StreetOptions)
			state := wf.GetState(100)
			require.NotNil(t, state)
			assert.Equal(t, tt.wantStep, state.Step)
			if tt.wantStreet != 0 {
				assert.Equal(t, tt.wantStreet, state.SelectedStreetID)
				assert.Equal(t, "Наукова", state.SelectedStreetName)
			}
		})
	}
}

func TestServiceInvalidBuildingKeepsState(t *testing.T) {
	wf, _ := newTestWorkflow(t, nil)
	startSearch(t, wf, 100)
	selectStreet(t, wf, 100, "Наукова")

	response := wf.Handle(100, Command{Kind: CommandText, Text: "bad"})

	assert.Equal(t, users.ErrInvalidBuildingFormat.Error(), response.Text)
	require.NotNil(t, wf.GetState(100))
	assert.Equal(t, StepSaveSubscription, wf.GetState(100).Step)
}

func TestServiceSuccessfulSaveClearsState(t *testing.T) {
	wf, repo := newTestWorkflow(t, nil)
	startSearch(t, wf, 100)
	selectStreet(t, wf, 100, "Наукова")

	response := wf.Handle(100, Command{Kind: CommandText, Text: "10"})

	assert.Equal(t, "Ви підписалися на сповіщення про відключення електроенергії для вулиці Наукова, будинок 10.", response.Text)
	require.NotNil(t, repo.users[100])
	assert.Equal(t, "Наукова", repo.users[100].Address.StreetName)
	assert.Equal(t, "10", repo.users[100].Address.Building)
	assert.Nil(t, wf.GetState(100))
}

func TestServiceRepositoryErrors(t *testing.T) {
	t.Run("save error", func(t *testing.T) {
		repo := newTestUserRepo()
		repo.saveErr = errors.New("disk error")
		wf, _ := newTestWorkflow(t, repo)
		startSearch(t, wf, 100)
		selectStreet(t, wf, 100, "Наукова")

		response := wf.Handle(100, Command{Kind: CommandText, Text: "10"})

		assert.Equal(t, messageGenericError, response.Text)
		assert.EqualError(t, response.Err, "disk error")
		assert.NotNil(t, wf.GetState(100))
	})

	t.Run("find error", func(t *testing.T) {
		repo := newTestUserRepo()
		repo.findErr = errors.New("disk error")
		wf, _ := newTestWorkflow(t, repo)

		response := wf.Handle(100, Command{Kind: CommandSubscription})

		assert.Equal(t, messageGenericError, response.Text)
		assert.EqualError(t, response.Err, "disk error")
	})

	t.Run("remove error", func(t *testing.T) {
		repo := newTestUserRepo()
		repo.removeErr = errors.New("disk error")
		wf, _ := newTestWorkflow(t, repo)

		response := wf.Handle(100, Command{Kind: CommandStop})

		assert.Equal(t, messageGenericError, response.Text)
		assert.EqualError(t, response.Err, "disk error")
	})
}

func TestServiceUnknownTextWithoutPendingStateIsIgnored(t *testing.T) {
	wf, _ := newTestWorkflow(t, nil)

	response := wf.Handle(100, Command{Kind: CommandText, Text: "Наукова"})

	assert.Empty(t, response.Text)
	assert.Empty(t, response.StreetOptions)
	assert.NoError(t, response.Err)
}

func TestServiceGetStateReturnsCopy(t *testing.T) {
	wf, _ := newTestWorkflow(t, nil)
	startSearch(t, wf, 100)

	state := wf.GetState(100)
	require.NotNil(t, state)
	state.Step = StepSaveSubscription
	state.SelectedStreetID = 99
	state.SelectedStreetName = "Змінена"

	state = wf.GetState(100)
	require.NotNil(t, state)
	assert.Equal(t, StepSearchStreet, state.Step)
	assert.Zero(t, state.SelectedStreetID)
	assert.Empty(t, state.SelectedStreetName)
}

func TestServiceStartRecoveryOnFindError(t *testing.T) {
	repo := newTestUserRepo()
	repo.findErr = errors.New("corrupt user file")
	wf, _ := newTestWorkflow(t, repo)

	response := wf.Handle(100, Command{Kind: CommandStart})

	assert.Equal(t, messagePromptStreet, response.Text)
	assert.EqualError(t, response.Err, "corrupt user file")
	state := wf.GetState(100)
	require.NotNil(t, state)
	assert.Equal(t, StepSearchStreet, state.Step)
}

func TestServiceExpiredPendingStateIsIgnored(t *testing.T) {
	wf, _ := newTestWorkflow(t, nil)
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	wf.now = func() time.Time { return base }
	wf.Handle(100, Command{Kind: CommandStart})

	wf.now = func() time.Time { return base.Add(31 * time.Minute) }
	response := wf.Handle(100, Command{Kind: CommandText, Text: "Наукова"})

	assert.Empty(t, response.Text)
	assert.Nil(t, wf.GetState(100))
}

func TestServiceActiveStateNotExpired(t *testing.T) {
	wf, _ := newTestWorkflow(t, nil)
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	wf.now = func() time.Time { return base }
	wf.Handle(100, Command{Kind: CommandStart})

	wf.now = func() time.Time { return base.Add(29 * time.Minute) }
	response := wf.Handle(100, Command{Kind: CommandText, Text: "Наукова"})

	assert.Equal(t, "Ви обрали вулицю: Наукова\nБудь ласка, введіть номер будинку:", response.Text)
	require.NotNil(t, wf.GetState(100))
	assert.Equal(t, StepSaveSubscription, wf.GetState(100).Step)
}

func TestServiceStreetToBuilding_PreservesStartedAt(t *testing.T) {
	wf, _ := newTestWorkflow(t, nil)
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	wf.now = func() time.Time { return base }
	wf.Handle(100, Command{Kind: CommandStart})

	wf.now = func() time.Time { return base.Add(5 * time.Minute) }
	wf.Handle(100, Command{Kind: CommandText, Text: "Наукова"})

	state := wf.GetState(100)
	require.NotNil(t, state)
	assert.Equal(t, base, state.StartedAt)
}

func TestServiceRestartResetsStartedAt(t *testing.T) {
	wf, _ := newTestWorkflow(t, nil)
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	wf.now = func() time.Time { return base }
	wf.Handle(100, Command{Kind: CommandStart})

	restart := base.Add(5 * time.Minute)
	wf.now = func() time.Time { return restart }
	wf.Handle(100, Command{Kind: CommandStart})

	state := wf.GetState(100)
	require.NotNil(t, state)
	assert.Equal(t, restart, state.StartedAt)
}

func TestNewWorkflow_HonorsConfigNow(t *testing.T) {
	fixed := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	wf := NewWorkflow(WorkflowConfig{
		UserRepo:   newTestUserRepo(),
		StreetRepo: &testStreetRepo{streets: testStreets()},
		Now:        func() time.Time { return fixed },
	})

	wf.Handle(100, Command{Kind: CommandStart})

	state := wf.GetState(100)
	require.NotNil(t, state)
	assert.Equal(t, fixed, state.StartedAt)
}
