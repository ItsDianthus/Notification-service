package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-ItsDianthus-NotificationLink/internal/api/openapi/scrapper_api"
	"go-ItsDianthus-NotificationLink/internal/scrapper/infrastructure/repo"
)

type LinksHandler struct {
	Repo *repo.SubscriptionRepo
}

func NewLinksHandler(r *repo.SubscriptionRepo) *LinksHandler {
	return &LinksHandler{Repo: r}
}

func (h *LinksHandler) PostLinks(c echo.Context, params scrapper_api.PostLinksParams) error {
	var body scrapper_api.PostLinksJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, scrapper_api.ApiErrorResponse{
			Description: scrapper_api.PtrString("invalid body"),
		})
	}
	chatID := params.TgChatId
	if body.Link != nil {
		h.Repo.AddLink(chatID, *body.Link)
	}
	return c.NoContent(http.StatusCreated)
}

func (h *LinksHandler) GetLinks(c echo.Context, params scrapper_api.GetLinksParams) error {
	chatID := params.TgChatId
	urls := h.Repo.ListLinks(chatID)

	resp := scrapper_api.ListLinksResponse{
		Size: scrapper_api.PtrInt32(int32(len(urls))),
	}
	list := make([]scrapper_api.LinkResponse, 0, len(urls))
	for i, url := range urls {
		list = append(list, scrapper_api.LinkResponse{
			Id:  scrapper_api.PtrInt64(int64(i)),
			Url: scrapper_api.PtrString(url),
		})
	}
	resp.Links = &list
	return c.JSON(http.StatusOK, resp)
}

func (h *LinksHandler) DeleteLinks(c echo.Context, params scrapper_api.DeleteLinksParams) error {
	var body scrapper_api.DeleteLinksJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, scrapper_api.ApiErrorResponse{
			Description: scrapper_api.PtrString("invalid body"),
		})
	}
	chatID := params.TgChatId
	if body.Link != nil {
		h.Repo.RemoveLink(chatID, *body.Link)
	}
	return c.NoContent(http.StatusNoContent)
}
