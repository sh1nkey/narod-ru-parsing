package main

import (
	"html-parser/kfk"
	"os"
	"sync"
	"time"

	//"github.com/jcmturner/gokrb5/v8/messages"
	"github.com/rs/zerolog/log"
)

const host = "localhost:9092"

func main() {

	producer, err := kfk.SetupProducer(host)
	if err != nil {
		log.Err(err).Msg("Не смогли запустить kafkaproducer")
	}
	message := kfk.SavedUrlEventDTO{
		Id: 123,
		HtmlContent: `        <div class="initials" style="line-height: 42px">
АК
        </div>

       </div>

      </div>

      <div class="body">

       <div class="pull_right date details" title="08.02.2024 07:32:42 UTC+03:00">
07:32
       </div>

       <div class="from_name">
Алексей Куделько
       </div>

       <div class="text">
где посты в паблике? 🤨<br>уже как 5 дней ничего не было
       </div>

       <div class="reactions">

        <div class="reaction">

         <div class="emoji">
🤣
         </div>
`,
	}
	message.FillBaseData()

	sign := make(chan *os.Signal)
	go func() {
		for {
			go kfk.ProduceMessage(producer, kfk.SavedUrlEvent, &message, &sign)
			time.Sleep(1 * time.Second)
		}
	}()


	

	time.Sleep(1)
	var wg sync.WaitGroup

	wg.Add(1)
	group := kfk.ConfigKfk(CleanText, kfk.ServeHtmlParsed, host)
	defer func() {
		log.Info().Err(err).Msg("Стопаем консюмер-группу")
		if err := group.Close(); err != nil {
			log.Error().Err(err).Msg("Error stopping consumer")
			return
		}
		log.Info().Err(err).Msg("Успешно стопнули консюмер группу")
	}()

	wg.Wait()


	// producer, err := kfk.SetupProducer(host)
	// if err != nil {
	// 	log.Err(err).Msg("Не смогли запустить kafkaproducer")
	// }
	// message := kfk.HtmlParsedTopicEventDTO{
	// 	Id: 123,
	// 	ParsedContent: "123",
	// }
	// message.FillBaseData()

	// sign := make(chan *os.Signal)
	// kfk.ProduceMessage(producer, kfk.HtmlParsedEvent, &message, &sign)

}