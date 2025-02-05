package parsenarod

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/render"
)


type Url struct {
	Id uint64 `json:"id"`
	Url string `json:"url"`
	Description sql.NullString `json:"description"`
	IsEmpty sql.NullBool `json:"is_empty"`
	CreatedAt time.Time `json:"created_at"`
}
type CreateUrlReqDTO struct {
	Url string `json:"url"`
}


func (cur *CreateUrlReqDTO) Bind(r *http.Request) error {
	return cur.validate()
}

func (cur *CreateUrlReqDTO) validate() error {
	_, err := url.ParseRequestURI(cur.Url)
	if err != nil {
		return errors.New("поданная ссылка не является https url-ом")
	}

	if !strings.Contains(cur.Url, "narod.ru") {
		return errors.New("поданная ссылка не является ссылкой на сайт narod.ru")
	}
	return nil
}


type UrlsResponse struct {
	Url
}

func (rd *UrlsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewUrlResponse(url Url) *UrlsResponse {
	resp := &UrlsResponse{Url: url}
	return resp
}

func NewUrlListResponse(products []Url) []render.Renderer {
	list := make([]render.Renderer, 0)
	for _, url := range products {
		list = append(list, NewUrlResponse(url))
	}
	return list
}


type MarkAsEmptyReq struct {
	Id uint64 `json:"id"`
}

func (m *MarkAsEmptyReq) Bind(r *http.Request) error {
	return nil
}

type SetDescriptionReq struct {
	Id uint64 `json:"id"`
	Description string `json:"description"`
}

func (m *SetDescriptionReq) Bind(r *http.Request) error {
	return nil
}
