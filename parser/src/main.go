package main

import (
	"sync"
	"time"
)

type Configuration struct {
	saver  Saver
	reader Checker
}

func main() {
	conf := Configuration{
		saver: writeToFile,
		reader: checkIfStrInFile,
	}

	var wg sync.WaitGroup
	chWeb := make(chan string)

	wg.Add(1)
	for i := 0; i <= 10; i++ {
		time.Sleep(1 * time.Second)
		go randStringBytesMaskImprSrcUnsafe(3, &wg, chWeb, conf.reader)
	}
	wg.Add(1)
	go checkIfWebPageExist(chWeb, &wg, conf.saver)

	wg.Wait()
}
