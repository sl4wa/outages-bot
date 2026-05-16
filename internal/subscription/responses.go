package subscription

import (
	"errors"
	"fmt"

	"outages-bot/internal/users"
)

const (
	messagePromptStreet       = "Будь ласка, введіть назву вулиці:"
	messageNoSubscription     = "Ви не маєте активної підписки."
	messageStreetOptions      = "Будь ласка, оберіть вулицю:"
	messageUnsubscribed       = "Ви успішно відписалися від сповіщень про відключення електроенергії."
	messageGenericError       = "Сталася помилка. Спробуйте пізніше."
	messageEmptyStreetQuery   = "Введіть назву вулиці."
	messageStreetNotFound     = "Вулицю не знайдено. Спробуйте ще раз."
	messagePromptStreetUpdate = "Ваша поточна підписка:\nВулиця: %s\nБудинок: %s\n\nБудь ласка, введіть нову назву вулиці для оновлення підписки:"
	messageCurrent            = "Ваша поточна підписка:\nВулиця: %s\nБудинок: %s"
	messagePromptBuilding     = "Ви обрали вулицю: %s\nБудь ласка, введіть номер будинку:"
	messageSaved              = "Ви підписалися на сповіщення про відключення електроенергії для вулиці %s, будинок %s."
)

func ignoredResponse() Response {
	return Response{}
}

func textResponse(text string) Response {
	return Response{Text: text}
}

func promptStreetResponse(current *users.User) Response {
	if current == nil {
		return textResponse(messagePromptStreet)
	}

	return textResponse(fmt.Sprintf(
		messagePromptStreetUpdate,
		current.Address.StreetName,
		current.Address.Building,
	))
}

func currentSubscriptionResponse(user *users.User) Response {
	return textResponse(fmt.Sprintf(
		messageCurrent,
		user.Address.StreetName,
		user.Address.Building,
	))
}

func streetOptionsResponse(options []users.Street) Response {
	names := make([]string, len(options))
	for i, opt := range options {
		names[i] = opt.Name
	}
	return Response{Text: messageStreetOptions, StreetOptions: names}
}

func promptBuildingResponse(streetName string) Response {
	return textResponse(fmt.Sprintf(messagePromptBuilding, streetName))
}

func savedSubscriptionResponse(user *users.User) Response {
	return textResponse(fmt.Sprintf(
		messageSaved,
		user.Address.StreetName,
		user.Address.Building,
	))
}

func invalidInputResponse(err error) Response {
	switch {
	case errors.Is(err, ErrEmptyStreetQuery):
		return textResponse(messageEmptyStreetQuery)
	case errors.Is(err, ErrStreetNotFound):
		return textResponse(messageStreetNotFound)
	default:
		return textResponse(err.Error())
	}
}

func errorResponse(err error) Response {
	return Response{Text: messageGenericError, Err: err}
}
