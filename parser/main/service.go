package main

import "parser/interfaces"

type checkParams struct {
	checker     interfaces.Checker
	saver       interfaces.Saver
	parser      interfaces.Parse
	hostChecker string
	hostWriter  string
}

func (c *checkParams) Check(text string) {
	c.checker(text, c)
}

func (c *checkParams) Save(text string, html *string) {
	c.saver(text, html, c.hostChecker)
}

func (c *checkParams) Parse(text string) {
	c.parser(text, c)
}

func (c *checkParams) GetCheckHostUrl() string {
	return "http://" + c.hostChecker + ":8083/api/v1/check"
}

func (c *checkParams) GetSaveHostUrl() string {
	return "http://" + c.hostWriter + ":8080/api/v1/check"
}