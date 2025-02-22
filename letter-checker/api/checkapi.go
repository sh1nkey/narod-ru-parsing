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
		r.Get("/read", ca.Read)
	})
}



func (ca *CheckApi) Read(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

    letters := queryParams.Get("letters") 
	
	log.Info().Msgf("получили запрос для букв %s", letters)

	responseStatusCode := ca.checker.ReadOne(letters)

	w.WriteHeader(responseStatusCode)
}


func (ca *CheckApi) Create(w http.ResponseWriter, r *http.Request) {
	data := &requester.LettersDTO{}
	if err := render.Bind(r, data); err != nil {
		log.Error().Err(err).Msg("Не смогли обработать данные запроса")
		Render(w, r, ErrInvalidRequest(err))
		return
	}

	w.WriteHeader(202)
	
	go ca.checker.WriteOne(data.Letters)
}







