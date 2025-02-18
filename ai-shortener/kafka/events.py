from datetime import datetime
import msgspec
from uuid import UUID


HTML_PARSED_TOPIC = "html_parsed"
AI_SUMMARIZED_TOPIC = "summarized"

class _BaseEventDTO(msgspec.Struct, gc=False):
    event_uuid: UUID
    correlation_uuid: UUID
    created_at: str

class HtmlParsedTopicEventDTO(_BaseEventDTO):
    url: str
    parsed_content: str


class AiSummarizedEventDTO(_BaseEventDTO):
    url: str
    summarized_content: str


