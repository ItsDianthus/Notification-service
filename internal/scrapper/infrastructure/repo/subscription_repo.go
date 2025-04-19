package repo

import "sync"

type SubscriptionRepo struct {
	mu            sync.RWMutex
	subscriptions map[int64]map[string]struct{}
}

func NewSubscriptionRepo() *SubscriptionRepo {
	return &SubscriptionRepo{
		subscriptions: make(map[int64]map[string]struct{}),
	}
}

func (r *SubscriptionRepo) RegisterChat(chatID int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.subscriptions[chatID]; !ok {
		r.subscriptions[chatID] = make(map[string]struct{})
	}
}

func (r *SubscriptionRepo) RemoveChat(chatID int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.subscriptions, chatID)
}

func (r *SubscriptionRepo) AddLink(chatID int64, url string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.subscriptions[chatID]; !ok {
		r.subscriptions[chatID] = make(map[string]struct{})
	}
	r.subscriptions[chatID][url] = struct{}{}
}

func (r *SubscriptionRepo) RemoveLink(chatID int64, url string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if links, ok := r.subscriptions[chatID]; ok {
		delete(links, url)
	}
}

func (r *SubscriptionRepo) ListLinks(chatID int64) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	links := []string{}
	if set, ok := r.subscriptions[chatID]; ok {
		for url := range set {
			links = append(links, url)
		}
	}
	return links
}
