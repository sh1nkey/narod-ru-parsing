

import uuid
from aiokafka import AIOKafkaProducer
import datetime

import msgspec
from kafka.events import AI_SUMMARIZED_TOPIC, AiSummarizedEventDTO, HtmlParsedTopicEventDTO
from loguru import logger

async def start_producer(server: str) -> AIOKafkaProducer:
    logger.info("Запускаем продюссер")
    producer = AIOKafkaProducer(
        bootstrap_servers=server,
    )
    await producer.start()
    return producer

async def send_one(producer: AIOKafkaProducer, value: str, data: HtmlParsedTopicEventDTO):
    send_data = AiSummarizedEventDTO(
        event_uuid=uuid.uuid4(),
        correlation_uuid=data.correlation_uuid,
        url=data.url,
        summarized_content=value,
        created_at=datetime.datetime.now(datetime.timezone.utc).isoformat()
    )
    await producer.send_and_wait(AI_SUMMARIZED_TOPIC, msgspec.json.encode(send_data))
