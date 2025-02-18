package kfk

import (
	"time"

	"github.com/google/uuid"
)

// Topics
const (
	//FailureTopic = "failure"
	// RequestedSaveUrlTopic   = "requested_save_url"
	// RequestedMarkEmptyTopic = "requested_mark_empty"
	// RequestedSetDescTopic   = "requested_set_description"
	SavedUrlTopic   = "saved_url"
	HtmlParsedTopic = "html_parsed"
	// AiSummarizedTopic       = "summarized"
)

// Possible Events Enum
const (
	LogFailureEvent = "failure"
	// RequestedSaveUrlEvent   = "requested_save_url"
	// RequestedMarkEmptyEvent = "requested_mark_empty"
	// RequestedSetDescEvent   = "requested_set_description"
	SavedUrlEvent   = "saved_url"
	HtmlParsedEvent = "html_parsed"
	// AiSummarizedEvent      = "summarized"
)

// Event Schemas
type baseEventDTO struct {
	EventUuid       uuid.UUID `json:"event_uuid"`
	CorrelationUuid uuid.UUID `json:"correlation_uuid"`
	CreatedAt       time.Time `json:"created_at"`
}

func (be *baseEventDTO) FillBaseData() {
	be.EventUuid = uuid.New()
	be.CreatedAt = time.Now()
}

type DbQueryId = uint64

// type LogFailureEventDTO struct {
// 	baseEventDTO
// 	FailedId uuid.UUID
// 	EventType uint8
// 	ErrorText string
// }

type SavedUrlEventDTO struct {
	baseEventDTO
	HtmlContent string `json:"html_content"`
	Url          string `json:"url"`
}

type HtmlParsedTopicEventDTO struct {
	baseEventDTO
	Url           string `json:"url"`
	ParsedContent string `json:"parsed_content"`
}

type Marshalable interface {
	HtmlParsedTopicEventDTO | SavedUrlEventDTO
}