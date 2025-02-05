package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"unsafe"
	"github.com/corpix/uarand"
)

func main() {
	var wg sync.WaitGroup
	chWeb := make(chan string)

	wg.Add(1)
	for i := 0; i <= 10; i++ {
		// time.Sleep(1 * time.Second)
		go randStringBytesMaskImprSrcUnsafe(3, &wg, chWeb)
	}
	wg.Add(1)
	go checkIfWebPageExist(chWeb, &wg)

	wg.Wait()
}



const letterBytes = "abcdefghijklmnopqrstuvwxyz123456789"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())



func randStringBytesMaskImprSrcUnsafe(n int, wg *sync.WaitGroup, chWeb chan string) {
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
		go checkIfStrInFile(*(*string)(unsafe.Pointer(&b)), chWeb)
	}
}


func checkIfStrInFile(text string, chWeb chan string) {
	b, err := os.ReadFile("t.txt")
	if err != nil {
		panic(err)
	}
	s := string(b)
	// //check whether s contains substring text
	isContaints := strings.Contains(s, text)
	if isContaints { return }
	log.Printf("положили текст в очередь для веба %s", text)
	chWeb <- text
}

func checkIfWebPageExist(chIn chan string, wg *sync.WaitGroup) {
	for {
		text := <- chIn

		log.Printf("забрали текст %s", text)

		time.Sleep(300 * time.Millisecond)
		url := "https://" + text + ".narod.ru" // Замените на нужный URL

		// Отправляем GET-запрос
		time.Sleep(150 * time.Millisecond)


		client := &http.Client{}
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                log.Fatalln(err)
        }
        req.Header.Set("User-Agent", uarand.GetRandom())
        resp, err := client.Do(req)
        if err != nil {
                log.Fatalln(err)
        }

		time.Sleep(150 * time.Millisecond)
		if err != nil {
			log.Fatal("Ошибка выполнения GET-запроса:", err)
		}
		// Проверяем статус-код ответа
		
		if resp.StatusCode == http.StatusOK {
			log.Printf("код: %d, на текст %s запрос успешно\n", resp.StatusCode, text)
			resp.Body.Close() 
			go writeToFile(text, wg)
		} else {
			time.Sleep(150 * time.Millisecond)
			log.Printf("Получен код состояния: %d, на текст %s запрос не удалось выполнить успешно\n", resp.StatusCode, text)
		}
	}
}

func writeToFile(text string, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.OpenFile("t.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("couldn't open file")
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("\n" + "https://" + text + ".narod.ru")); err != nil {
		log.Fatal("couldn't wrtie to file")
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal("couldn't close file")
		log.Fatal(err)
	}
	
}