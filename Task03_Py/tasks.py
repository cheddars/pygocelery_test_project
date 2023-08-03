from typing import Dict, Any

from celery import Celery
from kombu import Queue, Exchange

app = Celery('tasks.task03_process',
    broker='pyamqp://test:1234@localhost:5672//',
    backend='redis://localhost:6379',
    task_protocol=1
)
app.conf.task_queues = (
    Queue('task02_queue', Exchange('task02_ex'), routing_key='celeryx'),
)

@app.task
def task03_process(data: Dict[str, Any]) -> Dict[str, Any]:
    print(f"task03 {data}")

    result = {
        "id": data.get("id"),
        "status": "OK",
        "worker": "task03_Py"
    }
    return result
