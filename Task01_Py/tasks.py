import logging
import time
from typing import Any, Dict

from celery import Celery
from kombu import Queue, Exchange

from helper import generate_data

app = Celery('tasks',
    broker='pyamqp://test:1234@localhost:5672//',
    backend='redis://localhost:6379',
    task_protocol=2
)
app.conf.task_queues = (
    Queue('task01_queue', Exchange('task01_ex'), routing_key='celery'),
)

logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(levelname)s %(message)s')
logger = logging.getLogger()


@app.task
def task02_process(data: Dict[str, Any]) -> Dict[str, Any]:
    pass


def result_handler(id, result: Dict[str, Any]) -> None:
    logger.info(f"id : {id}")
    logger.info(f"data : {result}")


if __name__ == '__main__':
    iterations = 10000
    while iterations > 0:
        data: Dict[str, Any] = generate_data()
        logger.info(f"client side {data}")
        try:
            ar = task02_process.apply_async(args=(data,), exchange='task01_ex', routing_key='celery', serializer='json')
            ar.forget()

        except Exception as ex:
            logger.exception(ex)

        iterations -= 1