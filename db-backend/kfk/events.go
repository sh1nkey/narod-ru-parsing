package kfk

import (
	"time"

	"github.com/google/uuid"
)

// Topics
const (
	FailureTopic            = "failure"
	RequestedSaveUrlTopic   = "requested_save_url"
	RequestedMarkEmptyTopic = "requested_mark_empty"
	RequestedSetDescTopic   = "requested_set_description"
	SavedUrlTopic           = "saved_url"
	HtmlParsedTopic         = "html_parsed"
	AiSummarizedTopic       = "summarized"
)

// Possible Events Enum
const (
	LogFailureEvent         = "failure"
	RequestedSaveUrlEvent   = "requested_save_url"
	RequestedMarkEmptyEvent = "requested_mark_empty"
	RequestedSetDescEvent   = "requested_set_description"
	SavedUrlEvent           = "saved_url"
	HtmlParsedEvent         = "html_parsed"
	AiSummarizedEvent       = "summarized"
)

// Event Schemas
type BaseEventDTO struct {
	EventUuid       uuid.UUID `json:"event_uuid"`
	CorrelationUuid uuid.UUID `json:"correlation_uuid"`
	CreatedAt       time.Time `json:"created_at"`
}

func (be *BaseEventDTO) FillBaseData() {
	be.EventUuid = uuid.New()
	be.CreatedAt = time.Now()
}

type DbQueryId = uint64

type MarshalableEvent interface {
	RequestedMarkEmptyEventDTO | RequestedSetDescEventDTO | RequestedSaveUrlEventDTO | SavedUrlEventDTO | AiSummarizedEventDTO
}

// type LogFailureEventDTO struct {
// 	baseEventDTO
// 	FailedId uuid.UUID `json:"failed_at"`
// 	EventType uint8 `json:"event_type"`
// 	ErrorText string `json:"error_text"`
// }

type RequestedMarkEmptyEventDTO struct {
	BaseEventDTO
	Url string `json:"url"`
}

type RequestedSetDescEventDTO struct {
	BaseEventDTO
	Url         string  `json:"url"`
	Description string  `json:"description"`
}

type RequestedSaveUrlEventDTO struct {
	BaseEventDTO
	Url         string `json:"url"`
	HtmlContent string `json:"html_content"`
}

type SavedUrlEventDTO struct {
	BaseEventDTO
	HtmlContent string    `json:"html_content"`
	Url         string    `json:"url"`
}

// type HtmlParsedTopicEventDTO struct {
// 	baseEventDTO
// 	Url string `json:"url"`
// 	ParsedContent string `json:"parsed_content"`
// }

type AiSummarizedEventDTO struct {
	BaseEventDTO
	Url               string `json:"url"`
	SummarizedContent string    `json:"summarized_content"`
}
