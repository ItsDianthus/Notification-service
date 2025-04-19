package domain

type Subscription struct {
	ChatID int64
	Links  map[string]struct{}
}
