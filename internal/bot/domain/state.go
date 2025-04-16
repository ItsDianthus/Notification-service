package domain

type ConversationState string

const (
	StateDefault         ConversationState = "DEFAULT"
	StateAwaitingLink    ConversationState = "AWAITING_LINK"
	StateAwaitingTags    ConversationState = "AWAITING_TAGS"
	StateAwaitingFilters ConversationState = "AWAITING_FILTERS"
)
