package main

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type Configuration struct {
	saver  Saver
	checker Checker
}

func main() {
	log.Info().Msg("Начинаем работу")
	conf := Configuration{
		saver: writeToDb,
		checker: checkInDb,
	}

	var wg sync.WaitGroup
	chWeb := make(chan string)

	wg.Add(1)
	for i := 0; i <= 10; i++ {
		go randStringBytesMaskImprSrcUnsafe(3, &wg, chWeb, conf.checker)
	}
	wg.Add(1)
	go checkIfWebPageExist(chWeb, &wg, conf.saver)

	log.Info().Msg("Запустили проверку существования веб-страниц")
	wg.Wait()
}
