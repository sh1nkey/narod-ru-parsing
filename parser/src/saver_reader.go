package main

import (
	"os"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

type Saver func(text string, wg *sync.WaitGroup)
type Checker func(text string, chWeb chan string)


func writeToFile(text string, wg *sync.WaitGroup) {
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
