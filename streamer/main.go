package main

import (
	"fmt"
	"time"

	"encoding/json"
	"log"
	"os"

	"github.com/bucaojit/PolygonTools/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func main() {
	configFile := "../conf/polygon.yaml"
	webSocketEndpoint := "wss://socket.polygon.io/stocks"
	// CHANNELS := "T.SPY,Q.SPY"

	fmt.Println("Starting streamer...")

	args := os.Args[1:]
	if len(args) == 1 {
		configFile = os.Args[1]
	}
	conf, err := config.ReadConf(configFile)
	if err != nil {
		log.Fatal(err)
	}

	channels := config.TickerToChannels(conf.Tickers)
	kafkaTopic := config.KafkaTopic(conf)
	kafkaHosts := config.KafkaHosts(conf)

	fmt.Println(channels)
	fmt.Println(kafkaTopic)
	fmt.Println(kafkaHosts)

	c, _, err := websocket.DefaultDialer.Dial(webSocketEndpoint, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaHosts})
	defer p.Close()

	if err != nil {
		logrus.Fatal("Failed to create producer: ", err , "\n")
	}

	fmt.Printf("Created Producer %v\n", p)

	_ = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", conf.Apikey)))
	_ = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", channels)))

	chanMessages := make(chan interface{}, 10000)
	go processMessage(chanMessages, p, kafkaTopic)

	// Read messages off the buffered queue:
	/*
	go func() {
		for msgBytes := range chanMessages {
			logrus.Info("Message Bytes: ", msgBytes)
		}
	}()
	*/

	for {
		var msg interface{}
		err := c.ReadJSON(&msg)
		// logrus.Info(msg)

		if err != nil {
			panic(err)
		}
		chanMessages <- msg
	}
}

func processMessage(chanMessages <-chan interface{}, producer *kafka.Producer, topic string) {
	logrus.Info("topic: ", topic)

	deliveryChan := make(chan kafka.Event)
	logrus.Info("Before processing the messages")
	for msgBytes := range chanMessages {
		// logrus.Info("Message Bytes: ", msgBytes)
		value, _ := json.Marshal(msgBytes)
		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value,
			Key:            []byte{},
			Timestamp:      time.Time{},
			TimestampType:  0,
			Opaque:         nil,
			Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
		}, deliveryChan)
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			//fmt.Printf()
			logrus.Fatal("Delivery failed: ",m.TopicPartition.Error, "\n" )
		}
		if err != nil {
			logrus.Fatal("Producer failure: ",err, "\n", )
		}
	}
}
