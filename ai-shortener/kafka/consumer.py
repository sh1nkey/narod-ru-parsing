import asyncio
from aiokafka import AIOKafkaConsumer, AIOKafkaProducer
import msgspec
from ai.shorten_text import summarize_text

from .producer import send_one
from .events import HTML_PARSED_TOPIC, HtmlParsedTopicEventDTO

from loguru import logger

async def consume(producer: AIOKafkaProducer, host: str) -> None:
    logger.info("Запускаем консюмер")
    consumer = AIOKafkaConsumer(
        HTML_PARSED_TOPIC, 
        bootstrap_servers=host,
        group_id="ai_sum-0",
        enable_auto_commit=False,
        max_poll_interval_ms=1000000,
    )

    await consumer.start()
    logger.info("Запустили консюмер")
    try:
        # Consume messages
        async for msg in consumer:
            logger.info("Запуcтили")
            print("consumed: ", msg.topic, msg.partition, msg.offset, msg.key, msg.value, msg.timestamp)

            data: HtmlParsedTopicEventDTO = msgspec.json.decode(msg.value, type=HtmlParsedTopicEventDTO) # pyright: ignore
            shortened_text = summarize_text(data.parsed_content)

            await send_one(producer, shortened_text, data)
            await consumer.commit()
    
    finally:
        await consumer.stop()

