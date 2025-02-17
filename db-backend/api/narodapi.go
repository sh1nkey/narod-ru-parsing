package api

import (
	"context"
	"data-sender/core"
	"data-sender/core/parsenarod"
	"data-sender/kfk"
	"os"

	"fmt"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type NarodParseService interface {
	Create(ctx context.Context, url *parsenarod.UrlReqDTO) error
	GetAll(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]parsenarod.Url, error)
	MarkAsEmpty(ctx context.Context, url string, options ...core.UpdateOptions) error
	SetDescription(ctx context.Context, url string, description string, options ...core.UpdateOptions) error
}

const (
	CtxKeyNarodParse CtxKey = "narodparse"
)

type NarodParseApi struct {
	service NarodParseService
	kfkProd *sarama.AsyncProducer
}

func NewNarodParseApi(service NarodParseService, kfkProd *sarama.AsyncProducer) *NarodParseApi {
	return &NarodParseApi{
		service: service,
		kfkProd: kfkProd,
	}
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
	data := &parsenarod.UrlHtmlReqDTO{}
	if err := render.Bind(r, data); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		Render(w, r, ErrInvalidRequest(err))
		return
	}
	
	event := kfk.RequestedSaveUrlEventDTO{
		Url: data.Url, 
		HtmlContent: data.Html,
	}

	event.FillBaseData()
	event.BaseEventDTO.CorrelationUuid = uuid.New()

	sigCh := make(chan *os.Signal)
	go kfk.ProduceMessage(
		a.kfkProd,
		kfk.RequestedSaveUrlTopic, 
		&event, 
		&sigCh,
	)
	
	w.WriteHeader(http.StatusAccepted)
	log.Info().Msgf("запись с url %s успешно отправлена в Kafka", data.Url)
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
	data := &parsenarod.UrlReqDTO{}
	if err := render.Bind(r, data); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		Render(w, r, ErrInvalidRequest(err))
		return
	}

	event := &kfk.RequestedMarkEmptyEventDTO{Url: data.Url}
	event.FillBaseData()
	event.BaseEventDTO.CorrelationUuid = uuid.New()
	log.Info().Msgf("Получили запрос с url=%s", data.Url)
	
	sigCh := make(chan *os.Signal)
	go kfk.ProduceMessage(
		a.kfkProd,
		kfk.RequestedMarkEmptyTopic,
		event,
		&sigCh,
	)
	
	error := a.service.MarkAsEmpty(r.Context(), data.Url)
	if error != nil {
		log.Error().Err(error).Interface("MarkAsEmptyRequest", data).Msg("failed to reserve")
		Render(w, r, ErrInternalServer)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Info().Msg(fmt.Sprintf("запись с url=%s успешно помечена пустой", data.Url))
}

func (a *NarodParseApi) SetDescription(w http.ResponseWriter, r *http.Request) {
	data := &parsenarod.SetDescriptionReq{}
	if err := render.Bind(r, data); err != nil {
		log.Error().Err(err).Msg("Failed to bind request data")
		Render(w, r, ErrInvalidRequest(err))
		return
	}

	event := &kfk.RequestedSetDescEventDTO{Url: data.Url, Description: data.Description}
	event.FillBaseData()
	event.BaseEventDTO.CorrelationUuid = uuid.New()
	log.Info().Msgf("Получили запрос с url=%s", data.Url)
	
	sigCh := make(chan *os.Signal)
	go kfk.ProduceMessage(
		a.kfkProd,
		kfk.RequestedSetDescEvent,
		event,
		&sigCh,
	)

	err := a.service.SetDescription(r.Context(), data.Url, data.Description)
	if err != nil {
		log.Error().Err(err).Interface("SetDescriptionRequest", data).Msg("failed to reserve")
		Render(w, r, ErrInternalServer)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Info().Msg(fmt.Sprintf("запись с id %d успешно помечена пустой", data.Url))
}
