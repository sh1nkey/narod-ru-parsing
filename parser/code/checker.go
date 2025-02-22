package code

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)



// func checkIfStrInFile(text string) {
// 	b, err := os.ReadFile("t.txt")
// 	if err != nil {
// 		panic(err)
// 	}
// 	s := string(b)

// 	isContaints := strings.Contains(s, text)
// 	if isContaints {
// 		return
// 	}
// 	log.Printf("положили текст в очередь для веба %s", text)

// }

func NewServiceCheck(text string, host string, chanLet chan string) {

	//log.Info().Msgf("Отправляем на проверку текст %s", text)

	url := "http://" + host + ":8083/api/v1/check/read?letters=" + text
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		log.Err(err).Msg("Ошибка создания GET-запроса")
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("Ошибка выполнения GET-запроса")
		return
	}

	if resp.StatusCode == http.StatusNoContent {
		log.Info().Msgf("Проверка в БД. %d %s", resp.StatusCode, text)
		chanLet <- text

	}
	resp.Body.Close()
		
	//log.Info().Msgf("Получен код состояния: %d, на текст %s запрос не удалось выполнить успешно\n", resp.StatusCode, url)

}



func NewServiceWrite(text string, host string) {
	data := LettersDTO{
		Letters: text,
	}
	log.Info().Msgf("Отправляем на проверку текст %s", text)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msg("Ошибка маршалинга")
		return
	}

	url := "http://" + host + ":8083/api/v1/check/save"
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
		log.Info().Msgf("Проверка в БД. %s %s", resp.StatusCode, text)
		return
	}
	resp.Body.Close()
		
	log.Info().Msgf("Получен код состояния: %d, на текст %s запрос не удалось выполнить успешно\n", resp.StatusCode, url)
}