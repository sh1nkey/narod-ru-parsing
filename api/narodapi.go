package api

import (
	"context"
	"data-sender/core"
	"data-sender/core/parsenarod"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type NarodParseService interface {
	Create(ctx context.Context, url *parsenarod.CreateUrlReqDTO) error
	GetAll(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]parsenarod.Url, error)
	MarkAsEmpty(ctx context.Context, id uint64, options ...core.UpdateOptions) error
	SetDescription(ctx context.Context, id uint64, description string, options ...core.UpdateOptions) error
}

const (
	CtxKeyNarodParse CtxKey = "narodparse"
)

type NarodParseApi struct {
	service NarodParseService
}

func NewNarodParseApi(service NarodParseService) *NarodParseApi {
	return &NarodParseApi{service: service}
}

func (ra *NarodParseApi) ConfigureRouter(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.With(Paginate).Get("/", ra.List)
		r.Post("/save", ra.Create)
		r.Patch("/mark-empty", ra.MarkAsEmpty)
		r.Patch("/set-description", ra.SetDescription)
	})
}


func (a *NarodParseApi) Create(w http.ResponseWriter, r *http.Request) {
    data := &parsenarod.CreateUrlReqDTO{}
	if err := render.Bind(r, data); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		Render(w, r, ErrInvalidRequest(err))
		return
	}

    err := a.service.Create(r.Context(), data)
    if err != nil {
        log.Error().Err(err).Interface("reservationRequest", data).Msg("failed to reserve")
        Render(w, r, ErrInternalServer)
        return
    }

	w.WriteHeader(http.StatusCreated)
	log.Info().Msg("успешно добавлена запись с URL " + data.Url)
}



func (ra *NarodParseApi) List(w http.ResponseWriter, r *http.Request) {
	limit := r.Context().Value(CtxKeyLimit).(int)
	offset := r.Context().Value(CtxKeyOffset).(int)

	products, err := ra.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		log.Err(err).Send()
		Render(w, r, ErrInternalServer)
		return
	}
	render.Status(r, http.StatusOK)
	RenderList(w, r, parsenarod.NewUrlListResponse(products))
}

func (a *NarodParseApi) MarkAsEmpty(w http.ResponseWriter, r *http.Request) {
    data := &parsenarod.MarkAsEmptyReq{}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		err := errors.New("please, give me a query param called id")
		log.Error().Err(err).Msg("Failed to bind request data")
		Render(w, r, ErrInvalidRequest(err))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		Render(w, r, ErrInvalidRequest(err))
		return
    }
	log.Info().Msg(fmt.Sprintf(`id param gotten, id=$1`, id))



    error := a.service.MarkAsEmpty(r.Context(), id)
    if error != nil {
        log.Error().Err(error).Interface("MarkAsEmptyRequest", data).Msg("failed to reserve")
        Render(w, r, ErrInternalServer)
        return
    }

	w.WriteHeader(http.StatusNoContent)
	log.Info().Msg(fmt.Sprintf("запись с id %d успешно помечена пустой", data.Id))
}



func (a *NarodParseApi) SetDescription (w http.ResponseWriter, r *http.Request) {
    data := &parsenarod.SetDescriptionReq{}
	if err := render.Bind(r, data); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		Render(w, r, ErrInvalidRequest(err))
		return
	}

    err := a.service.SetDescription(r.Context(), data.Id, data.Description)
    if err != nil {
        log.Error().Err(err).Interface("SetDescriptionRequest", data).Msg("failed to reserve")
        Render(w, r, ErrInternalServer)
        return
    }

	w.WriteHeader(http.StatusNoContent)
	log.Info().Msg(fmt.Sprintf("запись с id %d успешно помечена пустой", data.Id))
}







// type MarkAsEmptyReq struct {
// 	id uint64
// }

// type SetDescriptionReq struct {
// 	id uint64
// 	description string
// }





