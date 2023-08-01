from celery import Celery

app = Celery('tasks',
    broker='pyamqp://test:1234@localhost:5672//',
    backend='redis://localhost:6379'
)

@app.task
def add(x, y):
    return x + y

if __name__ == '__main__':
    ar = add.apply_async((5456, 2878), serializer='json')
    print(ar.get())