package subscription

import "github.com/sl4wa/outages-bot/internal/outage/users"

func (w *Workflow) handleStart(chatID int64) Response {
	w.pending[chatID] = State{Step: StepSearchStreet, StartedAt: w.now()}

	current, err := w.userRepo.Find(chatID)
	if err != nil {
		resp := promptStreetResponse(nil)
		resp.Err = err
		return resp
	}
	return promptStreetResponse(current)
}

func (w *Workflow) handleStop(chatID int64) Response {
	delete(w.pending, chatID)

	removed, err := w.userRepo.Remove(chatID)
	if err != nil {
		return errorResponse(err)
	}
	if removed {
		return textResponse(messageUnsubscribed)
	}
	return textResponse(messageNoSubscription)
}

func (w *Workflow) handleSubscription(chatID int64) Response {
	user, err := w.userRepo.Find(chatID)
	if err != nil {
		return errorResponse(err)
	}
	if user == nil {
		return textResponse(messageNoSubscription)
	}
	return currentSubscriptionResponse(user)
}

func (w *Workflow) handleText(chatID int64, text string) Response {
	state, ok := w.pending[chatID]
	if !ok {
		return ignoredResponse()
	}

	if w.now().Sub(state.StartedAt) > w.ttl {
		delete(w.pending, chatID)
		return ignoredResponse()
	}

	switch state.Step {
	case StepSearchStreet:
		return w.handleSearchStreet(chatID, text)
	case StepSaveSubscription:
		return w.handleSaveSubscription(chatID, text, state)
	}
	return ignoredResponse()
}

func (w *Workflow) handleSearchStreet(chatID int64, text string) Response {
	result, err := w.searchStreet(text)
	if err != nil {
		return invalidInputResponse(err)
	}
	if len(result.options) > 0 {
		return streetOptionsResponse(result.options)
	}

	existing := w.pending[chatID]
	w.pending[chatID] = State{
		Step:               StepSaveSubscription,
		SelectedStreetID:   result.street.ID,
		SelectedStreetName: result.street.Name,
		StartedAt:          existing.StartedAt,
	}
	return promptBuildingResponse(result.street.Name)
}

func (w *Workflow) handleSaveSubscription(chatID int64, text string, state State) Response {
	addr, err := users.NewAddress(state.SelectedStreetID, state.SelectedStreetName, text)
	if err != nil {
		return invalidInputResponse(err)
	}

	user := &users.User{ID: chatID, Address: addr}
	if err := w.userRepo.Save(user); err != nil {
		return errorResponse(err)
	}

	delete(w.pending, chatID)
	return savedSubscriptionResponse(user)
}
