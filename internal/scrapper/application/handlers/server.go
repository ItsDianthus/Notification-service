package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-ItsDianthus-NotificationLink/internal/api/openapi/scrapper_api"
	"go-ItsDianthus-NotificationLink/internal/scrapper/infrastructure/repo"
)

type Server struct {
	Repo *repo.SubscriptionRepo
}

func NewServer(repo *repo.SubscriptionRepo) *Server {
	return &Server{Repo: repo}
}

func (s *Server) PostTgChatId(c echo.Context, id int64) error {
	s.Repo.RegisterChat(id)
	return c.NoContent(http.StatusCreated)
}

func (s *Server) DeleteTgChatId(c echo.Context, id int64) error {
	s.Repo.RemoveChat(id)
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) PostLinks(c echo.Context, params scrapper_api.PostLinksParams) error {
	var body scrapper_api.PostLinksJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, scrapper_api.ApiErrorResponse{
			Description: scrapper_api.PtrString("invalid body"),
		})
	}
	if body.Link != nil {
		s.Repo.AddLink(params.TgChatId, *body.Link)
	}
	return c.NoContent(http.StatusCreated)
}

func (s *Server) GetLinks(c echo.Context, params scrapper_api.GetLinksParams) error {
	urls := s.Repo.ListLinks(params.TgChatId)
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

func (s *Server) DeleteLinks(c echo.Context, params scrapper_api.DeleteLinksParams) error {
	var body scrapper_api.DeleteLinksJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, scrapper_api.ApiErrorResponse{
			Description: scrapper_api.PtrString("invalid body"),
		})
	}
	if body.Link != nil {
		s.Repo.RemoveLink(params.TgChatId, *body.Link)
	}
	return c.NoContent(http.StatusNoContent)
}
