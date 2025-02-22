package main

import (
	"parser/code"
	"sync"
	"github.com/rs/zerolog/log"
)


func main() {
	log.Info().Msg("Начинаем работу")
	conf := checkParams{
		check: code.NewServiceCheck,
		parse: code.CheckIfWebPageExist,
		saveUnexisting: code.NewServiceWrite,
		saveExisting: code.WriteToDb,
		hostCheck: "localhost",
		hostExistWrite: "localhost",
	}

	var wg sync.WaitGroup

	log.Info().Msg("Запустили проверку существования веб-страниц")

	chanLet := make(chan string)

	for i := 1; i < 10; i++{
		go RandStringBytesMaskImprSrcUnsafe(i, &conf, chanLet)
	}
	go conf.ManageParsed(chanLet)
	wg.Add(1)
	wg.Wait()

}
