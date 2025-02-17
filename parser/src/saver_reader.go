package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/corpix/uarand"
	"github.com/rs/zerolog/log"
)

type Saver func(text string, html *string, wg *sync.WaitGroup)
type Checker func(text string, chWeb chan string)

// HTTP запросы
type UrlHtmlReqDTO struct {
	Url string `json:"url"`
	Html string `json:"html"`
}


func writeToDb(url string, html *string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Info().Msgf("Отправляем текст %s", url)
	data := UrlHtmlReqDTO{
		Url: "https://" + url + ".narod.ru",
		Html: *html,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msg("Ошибка маршалинга")
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST", 
		"http://localhost:8080/api/v1/narod/save", 
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
		log.Info().Msgf("код: %d, на текст %s запрос успешно\n", resp.StatusCode, url)
		resp.Body.Close() 
	} else {
		log.Info().Msgf("Получен код состояния: %d, на текст %s запрос не удалось выполнить успешно\n", resp.StatusCode, url)
	}

}


func checkInDb(text string, chWeb chan string) {
	chWeb <- text
}




// Работа с файлом
func writeToFile(text string, html []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.OpenFile("t.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Err(err).Msg("couldn't open file")
	}
	if _, err := f.Write([]byte("\n" + "https://" + text + ".narod.ru")); err != nil {
		log.Err(err).Msg("couldn't wrtie to file")
	}
	if err := f.Close(); err != nil {
		log.Err(err).Msg("couldn't close file")
	}
}



func checkIfStrInFile(text string, chWeb chan string) {
	b, err := os.ReadFile("t.txt")
	if err != nil {
		panic(err)
	}
	s := string(b)

	isContaints := strings.Contains(s, text)
	if isContaints { return }
	log.Printf("положили текст в очередь для веба %s", text)
	chWeb <- text
}
