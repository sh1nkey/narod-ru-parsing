package code

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"parser/interfaces"
	"strings"

	"github.com/rs/zerolog/log"
)



func checkIfStrInFile(text string, params interfaces.CheckParamser) {
	b, err := os.ReadFile("t.txt")
	if err != nil {
		panic(err)
	}
	s := string(b)

	isContaints := strings.Contains(s, text)
	if isContaints {
		return
	}
	log.Printf("положили текст в очередь для веба %s", text)
	go params.Parse(text)
}

func NewServiceCheck(text string, params interfaces.CheckParamser) {
	data := LettersDTO{
		Letters: text,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msg("Ошибка маршалинга")
		return
	}

	url := params.GetCheckHostUrl()
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Err(err).Msg("Ошибка создания POST-запроса")
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("Ошибка выполнения POST-запроса")
		return
	}

	if resp.StatusCode == http.StatusNoContent {
		log.Info().Msgf("код: %d, на текст %s запрос успешно\n", resp.StatusCode, url)
		resp.Body.Close()

		go params.Parse(text)
	} else {
		log.Info().Msgf("Получен код состояния: %d, на текст %s запрос не удалось выполнить успешно\n", resp.StatusCode, url)
	}
}
