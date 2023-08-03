package main

import (
	"fmt"
	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
	"github.com/hashicorp/go-uuid"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	taskName  = "tasks.task03_process"
	redisHost = "localhost:6379"
	workerID  string
)

func init() {
	rnd, _ := uuid.GenerateUUID()
	workerID = "worker-" + rnd
}

func main() {
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	celeryClient, err := gocelery.NewCeleryClient(
		NewAMQPCeleryBrokerWithExchangeAndQueue("amqp://test:1234@localhost:5672//", "task02_ex", "task02_queue"),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		3,
	)

	if err != nil {
		log.Fatal("failed to create celery client ", err)
	}

	task03_process := func(data map[string]interface{}) map[string]interface{} {
		fmt.Println("data : ", data)
		fmt.Println("data id : ", data["id"])
		fmt.Println("data items : ", data["items"])
		m := make(map[string]interface{})
		m["id"] = data["id"]
		m["worker"] = "task03_Go"
		return m
	}

	celeryClient.Register(taskName, task03_process)

	go func() {
		celeryClient.StartWorker()
		fmt.Println("celery worker started. worker ID", workerID)
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit
	fmt.Println("exit signalled")

	celeryClient.StopWorker()
	fmt.Println("celery worker stopped")
}

func NewAMQPCeleryBrokerWithExchangeAndQueue(host string, exchange string, queue string) *gocelery.AMQPCeleryBroker {
	conn, channel := gocelery.NewAMQPConnection(host)
	broker := &gocelery.AMQPCeleryBroker{
		Channel:    channel,
		Connection: conn,
		Exchange:   NewAMQPExchange(exchange, "direct"),
		Queue:      gocelery.NewAMQPQueue(queue),
		Rate:       4,
	}
	if err := broker.CreateExchange(); err != nil {
		panic(err)
	}
	if err := broker.CreateQueue(); err != nil {
		panic(err)
	}
	if err := broker.Qos(broker.Rate, 0, false); err != nil {
		panic(err)
	}
	if err := broker.StartConsumingChannel(); err != nil {
		panic(err)
	}
	return broker
}

func NewAMQPExchange(name string, typeName string) *gocelery.AMQPExchange {
	return &gocelery.AMQPExchange{
		Name:       name,
		Type:       typeName,
		Durable:    true,
		AutoDelete: false,
	}
}
