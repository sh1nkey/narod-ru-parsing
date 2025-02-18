


import asyncio
from time import sleep
from kafka.consumer import consume
from kafka.producer import start_producer
from loguru import logger


async def start_kafka():
    #sleep(30)
    logger.info("Запускаем продюссер и консюмер")
    producer = await start_producer('localhost:9092')
    logger.info("Запускаем консюмер")
    await consume(producer, "localhost:9092")

asyncio.run(start_kafka())