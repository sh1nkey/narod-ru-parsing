package interfaces

type Checker func(text string, params CheckParamser)
type Parse func(text string, params CheckParamser)
type Saver func(text string, html *string, host string)


type CheckParamser interface {
	Check(text string)
	Save(text string, html *string)
	Parse(text string)
	GetSaveHostUrl() string
	GetCheckHostUrl() string
}