package code

import (
	"io"
	"net/http"
	"time"

	"github.com/corpix/uarand"
	"github.com/rs/zerolog/log"
)



func CheckIfWebPageExist(text string) (string, string, bool) {
	log.Info().Msgf("забрали текст %s", text)

	time.Sleep(300 * time.Millisecond)
	url := "https://" + text + ".narod.ru"

	time.Sleep(150 * time.Millisecond)


	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Err(err).Msg("Ошибка создания GET-запроса")
	}
	req.Header.Set("User-Agent", uarand.GetRandom())

	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msg("Ошибка выполнения GET-запроса")
	}
	time.Sleep(150 * time.Millisecond)
	
	if resp.StatusCode == http.StatusOK {
		html, err := io.ReadAll(resp.Body)
		defer resp.Body.Close() 
		if err != nil {
			log.Err(err).Msg("Ошибка чтения тела ответа")
			return url, "", false
		}
		htmlData := string(html)


		log.Info().Msgf("Проверка существования. %s %s", resp.StatusCode, text)
		return url, htmlData, true

	} 

	time.Sleep(150 * time.Millisecond)
	log.Printf("Получен код состояния: %d, на текст %s запрос не удалось выполнить успешно\n", resp.StatusCode, text)

	return url, "", false

}
