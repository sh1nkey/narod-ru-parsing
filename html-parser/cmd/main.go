package main

import (
	"html-parser/kfk"
	"sync"

	"github.com/rs/zerolog/log"
)

const host = "makafka:9092"

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	group := kfk.ConfigKfk(CleanText, kfk.ServeHtmlParsed, host)
	defer func() {
		if err := group.Close(); err != nil {
			log.Error().Err(err).Msg("Error stopping consumer")
			return
		}
		log.Info().Msg("Успешно стопнули консюмер группу")
	}()

	wg.Wait()


}