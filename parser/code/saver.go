package code

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/corpix/uarand"
	"github.com/rs/zerolog/log"
)

// HTTP запросы
type UrlHtmlReqDTO struct {
	Url string `json:"url"`
	Html string `json:"html"`
}

// HTTP запросы
type LetterCheckDTO struct {
	Url string `json:"url"`
	Html string `json:"html"`
}

type LettersDTO struct {
	Letters string `json:"letters"`
}


func WriteToDb(url string, html *string, host string) {

	//log.Info().Msgf("Отправляем текст %s", url)
	data := UrlHtmlReqDTO{
		Url: "https://" + url + ".narod.ru",
		Html: *html,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msg("Ошибка маршалинга")
		return
	}

	urlSave := "http://" + host + ":8080/api/v1/narod/save"
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST", 
		urlSave, 
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Err(err).Msg("Ошибка создания GET-запроса")
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("Ошибка выполнения GET-запроса")
	}
	
	if resp.StatusCode == http.StatusAccepted  {
		log.Info().Msgf("Добавление в БД. %s %s", resp.StatusCode, url)
		resp.Body.Close() 
	} else {
		log.Info().Msgf("Получен код состояния: %d, на текст %s запрос не удалось выполнить успешно\n", resp.StatusCode, url)
	}

}


// // Работа с файлом
// func writeToFile(text string, html *string, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	f, err := os.OpenFile("t.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		log.Err(err).Msg("couldn't open file")
// 	}
// 	if _, err := f.Write([]byte("\n" + "https://" + text + ".narod.ru")); err != nil {
// 		log.Err(err).Msg("couldn't wrtie to file")
// 	}
// 	if err := f.Close(); err != nil {
// 		log.Err(err).Msg("couldn't close file")
// 	}
// }
