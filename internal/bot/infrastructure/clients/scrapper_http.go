package clients

import (
	"net/http"
)

type scrapperHTTPClient struct {
	baseURL string
	client  *http.Client
}

//
//func (c *scrapperHTTPClient) AddLink(ctx context.Context, chatID int64, link string, tags, filters []string) (*bot_api.LinkUpdate, error) {
//	// 1) Собрать URL: c.baseURL + "/links"
//	// 2) Подготовить тело запроса (JSON { link, tags, filters })
//	// 3) Сделать POST-запрос с контекстом: req = http.NewRequestWithContext(ctx, "POST", url, body)
//	// 4) Проверить статус-код, распарсить ответ в Link, вернуть
//}
//
//func (c *scrapperHTTPClient) RemoveLink(ctx context.Context, chatID int64, link string) error {
//	// Аналогично, но DELETE /links с телом { link }
//}
//
//func (c *scrapperHTTPClient) ListLinks(ctx context.Context, chatID int64) ([]bot_api.LinkUpdate, error) {
//	// GET /links, заголовок Tg-Chat-Id: chatID, распарсить массив LinkResponse → []Link
//}
//
//// NewScrapperHTTPClient возвращает HTTP‑реализацию этого интерфейса.
//func NewScrapperHTTPClient(baseURL string, timeout time.Duration) clients.ScrapperClient {
//	return &scrapperHTTPClient{
//		baseURL: baseURL,
//		client:  &http.Client{Timeout: timeout},
//	}
//}
