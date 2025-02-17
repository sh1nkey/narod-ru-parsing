package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"data-sender/config"
)


const (
	HealthEndpoint  = "/health"

	ApiPath         = "/api/v1"
	NarodParsePath = "/narod"

)



func ConfigureRouter(cfg *config.Config, integServ NarodParseService) chi.Router {
	log.Info().Msg("configuring router...")
	r := chi.NewRouter()

	r.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins: []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: true,
				MaxAge:           300,
			},
		),
	)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.HandleFunc(ApiPath, func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("UP"))
		},
	)

	r.Route(ApiPath, func(r chi.Router) {
		r.Route(NarodParsePath, NewNarodParseApi(integServ).ConfigureRouter)
	})

	return r

}


func Render(w http.ResponseWriter, r *http.Request, rnd render.Renderer) {
	if err := render.Render(w, r, rnd); err != nil {
		log.Warn().Err(err).Msg("failed to render")
	}
}

func RenderList(w http.ResponseWriter, r *http.Request, l []render.Renderer) {
	if err := render.RenderList(w, r, l); err != nil {
		log.Warn().Err(err).Msg("failed to render")
	}
}