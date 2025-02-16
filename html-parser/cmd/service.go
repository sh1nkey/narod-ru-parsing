package main

import (
	"regexp"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)




func CleanText(htmlContent *string) (string, error) {
	log.Debug().Msgf("Получили текст для парсинга: %s", *htmlContent)

	log.Info().Msg("Парсим текст...")
	doc, err := html.Parse(strings.NewReader(*htmlContent))
	if err != nil {
		log.Err(err).Msg("Ошибка при открытии читателя при парсинге HTML-строки, event: ServeHtmlParsed")
		return "", err
	}

	var textContent strings.Builder

	var wg sync.WaitGroup
	extractText(doc, &textContent, &wg)
	wg.Wait()

	cleanedText := strings.ReplaceAll(textContent.String(), "\n", "")
	log.Debug().Msgf("Спарсили текст %s", cleanedText)
	
	re := regexp.MustCompile(`\s+`) // регулярное выражение для последовательностей пробелов
	cleanedText = re.ReplaceAllString(cleanedText, " ")

	// Удаляем пробелы в начале и в конце строки
	cleanedText = strings.Join(strings.Fields(cleanedText), " ")
	return cleanedText, nil
}


func extractText(n *html.Node, textContent *strings.Builder, wg *sync.WaitGroup) {
	if n.Type == html.TextNode {
		wg.Add(1)
		go func() {
			defer wg.Done() 
			_, err := textContent.WriteString(n.Data)
			if err != nil {
				log.Err(err).Msg("Ошибка при записи в строку, event: ServeHtmlParsed")
			}
		}()

	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, textContent, wg)
	}
}