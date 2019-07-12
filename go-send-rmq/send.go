package main

import (
	"log"
	"fmt"
	"os"
	"flag"
	"strings"
	"github.com/streadway/amqp"
	"github.com/krionbsd/go-amqp-reconnect/rabbitmq"
)

var exchangeName = flag.String("e", "pool1", "AMQP Exchange name to send")
var queueName = flag.String("q", "vm_action", "AMQP Queue name in Exchange to send")

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog 'payload-to-send'\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	rabbitmq.Debug = false

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("argv1 is missing");
		os.Exit(1);
	}

	//todo: to cfg
	conn, err := rabbitmq.Dial("amqp://vi_user:vi@10.64.67.240:5672/vi_vhost")

	if err != nil {
		log.Panic(err)
	}

	sendCh, err := conn.Channel()
	if err != nil {
		log.Panic(err)
	}

	err = sendCh.ExchangeDeclare(*exchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		log.Panic(err)
	}

	_, err = sendCh.QueueDeclare(*queueName, true, false, false, false, nil)
	if err != nil {
		log.Panic(err)
	}

	if err := sendCh.QueueBind(*queueName, "", *exchangeName, false, nil); err != nil {
		log.Panic(err)
	}

//	body1 := args[0]
	body1 := os.Args[1:]
//	fmt.Printf("-> [%s]\n",body)
//	fmt.Printf("-> %T [%s]\n",body1, body1)
	body := strings.Join(body1," ")
	fmt.Printf("-> Send: %s\n",body)
//	os.Exit(1)
//	body := string(os.Args[1:])

	err = sendCh.Publish(*exchangeName, "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(body),
	})
	log.Printf(" [x] Sent to %s:%s, payload[%s]<\n", *exchangeName, *queueName,body)
	failOnError(err, "Failed to publish a message")
}
