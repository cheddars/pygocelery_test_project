from datetime import datetime
from typing import Any, Dict

from celery import Celery
from kombu import Queue, Exchange

app = Celery('tasks',
    broker='pyamqp://test:1234@localhost:5672//',
    backend='redis://localhost:6379',
    task_protocol=2
)
app.conf.task_queues = (
    Queue('task01_queue', Exchange('task01_ex'), routing_key='celery'),
)

app_v1 = Celery('tasks',
    broker='pyamqp://test:1234@localhost:5672//',
    backend='redis://localhost:6379',
    task_protocol=1
)
app_v1.conf.task_queues = (
    Queue('task02_queue', Exchange('task02_ex'), routing_key='celeryx'),
)

@app.task
def task02_process(data: Dict[str, Any]) -> Dict[str, Any]:
    print(data)
    data["task02"] = datetime.now().strftime("%Y-%m-%d %H:%M:%S.%f")[:-3]
    arr = task03_process.apply_async(args=(data,), exchange="task02_ex", routing_key='celeryx', serializer='json')
    arr.forget()

    result = {
        "id": data.get("id"),
        "status": "OK"
    }
    return result

@app_v1.task
def task03_process(data: Dict[str, Any]) -> Dict[str, Any]:
    print(data)
    result = {
        "id": data.get("id"),
        "status": "OK"
    }
    return result