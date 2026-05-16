package subscription

import "errors"

var (
	ErrEmptyStreetQuery = errors.New("empty street query")
	ErrStreetNotFound   = errors.New("street not found")
)
