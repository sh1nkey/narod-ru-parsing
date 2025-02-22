package main

import (
	"parser/code"
	"sync"
	"github.com/rs/zerolog/log"
)


func main() {
	log.Info().Msg("Начинаем работу")
	conf := checkParams{
		saver: code.WriteToDb,
		hostChecker: "localhost",
		hostWriter: "localhost",
		parser: code.CheckIfWebPageExist,
		checker: code.NewServiceCheck,
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++{
		go code.RandStringBytesMaskImprSrcUnsafe(i, &conf)
	}

	wg.Add(1)
	log.Info().Msg("Запустили проверку существования веб-страниц")
	wg.Wait()
}
