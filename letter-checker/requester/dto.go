package requester

import "net/http"

type LettersDTO struct {
	Letters string `json:"letters"`
}

func (l *LettersDTO) Bind(r *http.Request) error {
	return nil
}