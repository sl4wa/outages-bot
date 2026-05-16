package subscription

import (
	"github.com/sl4wa/outages-bot/internal/outage/users"
	"time"
)

const defaultPendingTTL = 30 * time.Minute

// CommandKind identifies a subscription command.
type CommandKind int

const (
	CommandText CommandKind = iota
	CommandStart
	CommandStop
	CommandSubscription
)

// Command is an application-level subscription command.
type Command struct {
	Kind CommandKind
	Text string
}

// Response is ready to send by adapters, with optional street options for reply keyboards.
type Response struct {
	Text          string
	StreetOptions []string
	Err           error
}

// StepKind identifies a step in the subscription conversation.
type StepKind int

const (
	_ StepKind = iota
	StepSearchStreet
	StepSaveSubscription
)

// State holds the state of a user's subscription conversation.
type State struct {
	Step               StepKind
	SelectedStreetID   int
	SelectedStreetName string
	StartedAt          time.Time
}

// UserRepository provides the user persistence operations required by subscription actions.
type UserRepository interface {
	Find(chatID int64) (*users.User, error)
	Save(user *users.User) error
	Remove(chatID int64) (bool, error)
}

// StreetRepository provides the street data required by street search.
type StreetRepository interface {
	GetAllStreets() []users.Street
}

// Workflow owns the subscription conversation workflow.
type Workflow struct {
	userRepo   UserRepository
	streetRepo StreetRepository
	pending    map[int64]State
	ttl        time.Duration
	now        func() time.Time
}

// WorkflowConfig holds configuration for Workflow.
type WorkflowConfig struct {
	UserRepo   UserRepository
	StreetRepo StreetRepository
	TTL        time.Duration
	Now        func() time.Time
}

// NewWorkflow creates a new subscription conversation workflow.
func NewWorkflow(cfg WorkflowConfig) *Workflow {
	ttl := cfg.TTL
	if ttl == 0 {
		ttl = defaultPendingTTL
	}
	now := cfg.Now
	if now == nil {
		now = time.Now
	}
	return &Workflow{
		userRepo:   cfg.UserRepo,
		streetRepo: cfg.StreetRepo,
		pending:    make(map[int64]State),
		ttl:        ttl,
		now:        now,
	}
}

// Handle applies a command to the subscription workflow.
func (w *Workflow) Handle(chatID int64, cmd Command) Response {
	switch cmd.Kind {
	case CommandStart:
		return w.handleStart(chatID)
	case CommandStop:
		return w.handleStop(chatID)
	case CommandSubscription:
		return w.handleSubscription(chatID)
	case CommandText:
		return w.handleText(chatID, cmd.Text)
	default:
		return ignoredResponse()
	}
}

// GetState returns a copy of the conversation state for tests and adapters.
func (w *Workflow) GetState(chatID int64) *State {
	state, ok := w.pending[chatID]
	if !ok {
		return nil
	}
	return &state
}
