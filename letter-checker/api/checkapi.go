package api

import (
	"letter-checker/requester"

	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)





type CheckApi struct {
	checker requester.Checker
}

func NewCheckApi(service requester.Checker) *CheckApi {
	return &CheckApi{
		checker: service,
	}
}

func (ca *CheckApi) ConfigureRouter(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Post("/save", ca.Create)
	})
}

func (ca *CheckApi) Create(w http.ResponseWriter, r *http.Request) {
	data := &requester.LettersDTO{}
	if err := render.Bind(r, data); err != nil {
		log.Error().Err(err).Msg("Не смогли обработать данные запроса")
		Render(w, r, ErrInvalidRequest(err))
		return
	}
	
	responseStatusCode := ca.checker.Check(data.Letters)

	w.WriteHeader(responseStatusCode)
}

