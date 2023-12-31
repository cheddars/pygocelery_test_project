package main

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
)

const (
	redisHostEnvVar = "REDIS_HOST"
	taskNameEnvVar  = "TASK_NAME"
	smtpServer      = "localhost:1025"
	taskName        = "tasks.add"
)

const fromEmail = "admin@foo.com"

const emailBodyTemplate = "Hi %s!!\n\nHere is your auto-generated password %s. Visit https://foobar.com/login to login update your password.\n\nCheers,\nTeam FooBar.\n\n[processed by %s]"

const autogenPassword = "foobarbaz_foobarbaz"

const emailHeaderTemplate = "From: %s" + "\n" +
	"To: %s" + "\n" +
	"Subject: Welcome to FooBar! Here are your login instructions\n\n" +
	"%s"

var (
	redisHost string
	workerID  string
)

func init() {
	redisHost = os.Getenv(redisHostEnvVar)
	if redisHost == "" {
		redisHost = "localhost:6379"
	}

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
		gocelery.NewAMQPCeleryBroker("amqp://test:1234@localhost:5672//"),
		// gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		1,
	)

	if err != nil {
		log.Fatal("failed to create celery client ", err)
	}

	sendEmail := func(num1 int, num2 int) {
		fmt.Println("num1 : ", num1)
		fmt.Println("num2 : ", num2)

		registrationEvent := "test"
		name := strings.Split(registrationEvent, ",")[0]
		userEmail := strings.Split(registrationEvent, ",")[1]

		fmt.Println("user registration info:", name, userEmail)

		sleepFor := rand.Intn(9) + 1
		time.Sleep(time.Duration(sleepFor) * time.Second)

		body := fmt.Sprintf(emailBodyTemplate, name, autogenPassword, workerID)
		msg := fmt.Sprintf(emailHeaderTemplate, fromEmail, userEmail, body)

		err := smtp.SendMail(smtpServer, nil, "test@localhost", []string{"foo@bar.com"}, []byte(msg))
		if err != nil {
			log.Fatal("failed to send email - ", err)
		}

		fmt.Println("sent email to", userEmail)

	}

	celeryClient.Register(taskName, sendEmail)

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
