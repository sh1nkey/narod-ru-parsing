package main

import "strings"

type checker func(text string, host string, chanLet chan string)
type parser func(text string) (url string, html string, exists bool)
type saverUnexisting func(text string, host string)
type saverExisting func(url string, html *string, host string)

type checkParams struct {
	check          checker
	parse          parser
	saveUnexisting saverUnexisting
	saveExisting   saverExisting
	hostCheck      string
	hostExistWrite string
}

func (cp *checkParams) Check(text string, chanLet chan string) {
	go cp.check(text, cp.hostCheck, chanLet)
}

func (cp *checkParams) ManageParsed(chanLet chan string) {
	letters := <- chanLet

	url, html, exists := cp.parse(letters)

	if exists {
		go cp.saveExisting(url, &html, cp.hostExistWrite)
		return
	}
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, ".narod.ru")

	go cp.saveUnexisting(url, cp.hostCheck)
}