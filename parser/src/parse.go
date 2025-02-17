package main

import (
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"unsafe"

	"github.com/corpix/uarand"
	"github.com/rs/zerolog/log"
)


const letterBytes = "abcdefghijklmnopqrstuvwxyz123456789"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())



func randStringBytesMaskImprSrcUnsafe(n int, wg *sync.WaitGroup, chWeb chan string, check Checker) {
	log.Info().Msg("Запускаем генерацию случайных строк")
	
	for i := 0; i <= 1_000; i++ {
		wg.Add(1)
		time.Sleep(300 * time.Millisecond)
		b := make([]byte, n)
		for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
			time.Sleep(1 * time.Second)
			if remain == 0 {
				cache, remain = src.Int63(), letterIdxMax
			}
			if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
				b[i] = letterBytes[idx]
				i--
			}
			cache >>= letterIdxBits
			remain--
		}
		go check(*(*string)(unsafe.Pointer(&b)), chWeb)
	}
}


func checkIfWebPageExist(chIn chan string, wg *sync.WaitGroup, save Saver) {
	log.Info().Msg("Запускаем проверку существования веб-страниц")
	for {
		text := <- chIn
		log.Info().Msgf("забрали текст %s", text)


		time.Sleep(300 * time.Millisecond)
		url := "https://" + text + ".narod.ru" // Замените на нужный URL

		// Отправляем GET-запрос
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
			if err != nil {
				log.Err(err).Msg("Ошибка чтения тела ответа")
				return
			}
			htmlData := string(html)
			go save(text, &htmlData, wg)

			log.Info().Msgf("код: %d, на текст %s запрос успешно\n", resp.StatusCode, text)
			//log.Info().Msgf("HTML-содержаие: %s", string(html))

			resp.Body.Close() 
		} else {
			time.Sleep(150 * time.Millisecond)
			log.Printf("Получен код состояния: %d, на текст %s запрос не удалось выполнить успешно\n", resp.StatusCode, text)
		}
	}
}
