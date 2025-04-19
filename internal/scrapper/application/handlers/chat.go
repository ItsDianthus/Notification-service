package handlers

import (
	"github.com/labstack/echo/v4"
	"go-ItsDianthus-NotificationLink/internal/scrapper/infrastructure/repo"
	"net/http"
)

type ChatHandler struct {
	Repo *repo.SubscriptionRepo
}

func NewChatHandler(r *repo.SubscriptionRepo) *ChatHandler {
	return &ChatHandler{Repo: r}
}

func (h *ChatHandler) PostTgChatId(c echo.Context, id int64) error {
	h.Repo.RegisterChat(id)
	return c.NoContent(http.StatusCreated)
}

func (h *ChatHandler) DeleteTgChatId(c echo.Context, id int64) error {
	h.Repo.RemoveChat(id)
	return c.NoContent(http.StatusNoContent)
}
