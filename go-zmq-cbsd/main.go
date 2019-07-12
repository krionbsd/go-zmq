package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/krionbsd/go-amqp-reconnect/rabbitmq"
	"github.com/streadway/amqp"
	"log"
	"os"
	"sync"
	"time"
)

var exchangeName = flag.String("e", "pool1", "AMQP Exchange name to send")
var queueName = flag.String("q", "vm_action", "AMQP Queue name in Exchange to send")

func usage() {
	_, err := fmt.Fprintf(os.Stderr, "usage: myprog [-e exchange] [-q queue]'\n")
	if err != nil {
		panic(err)
	}
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	rabbitmq.Debug = false

	flag.Usage = usage
	flag.Parse()

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

	go func() {
		for {
			err := sendCh.Publish(*exchangeName, "", false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:	[]byte("ping"),
			})
			if err != nil {
				log.Printf("publish, err: %v", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()

	consumeCh, err := conn.Channel()
	if err != nil {
		log.Panic(err)
	}

	go func() {
		d, err := consumeCh.Consume(*queueName, "", false, false, false, false, nil)
		if err != nil {
			log.Panic(err)
		}

		for msg := range d {
			if string(msg.Body) == "ping" {
				//log.Printf("PING RECEIVED")
				msg.Ack(true)
				continue
			}
			log.Printf("recv msg: %s", string(msg.Body))
			comment := Comment{}
			err := json.Unmarshal(msg.Body, &comment)
			if err != nil {
					log.Printf("json decode error %s", err)
					msg.Ack(true)
					continue
			}

			comment.JobID = uint64( time.Now().UnixNano() )
			fmt.Printf("Time now %d\n", comment.JobID)
			err = DoProcess(&comment)
			if err != nil {
				panic(err)
			}
			msg.Ack(true)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
