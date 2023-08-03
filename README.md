# 복합 Celery 구조 테스트

- 기존의 celery 기반 프로젝트가 있는데 (python 기반) process chain 을 celery 로 관리하고 있다
- 하나의 작업은 다수의 처리 절차가 필요하고 각 단계를 celery 로 관리하고 있다

Task01 -> Task02 -> Task03


- Task01 -> Task02, Task02 -> Task03 각각의 단계에서 celery 를 이용해 데이터를 전달하는데 이중 Task03 을 golang 기반으로 변경하려고 한다

- Task03 Worker 를 golang 기반으로도 돌리고 기존의 Python 기반으로도 돌리면서 당분간 운영해보려고 한다
- 이때 문제되는 부분으로 gocelery 모듈은 celery protocol 2 지원하는 않는다는 것이 있다
- 전체 환경은 protocol 2 기반으로 돌리고 Task02 -> Task03 부분만 protocol 1 기반으로 변경가능한지 확인 하는 것이 이 프로젝트의 목적이다
- Task02 는 celery worker 이면서 동시에 client 이다
- client 로서는 protocol 1 을 사용해 Task03 으로 데이터를 전달
- worker 로서는 protocol 2 를 사용해 기존의 Task01 로 부터 오는 데이터 처리

구현하려는 구조는 아래와 같다

```
           protocol v2              protocol v1
Task01_Py ------------> Task02_Py -------------> Task03_Go
                                  \
                                   \--------> Task03_Py
```

### 테스트 해보니...

gocelery 는 amqp 의 기본 exchange(default) 설정되어 있고 다른 exchange 와 queue 를 사용하려면 수정을 좀 많이 해야할듯 하다
기본 exchange 는 잘 동작되는게 확인된다
