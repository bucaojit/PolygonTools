package main

import (
	//"bufio"
	"fmt"
	"time"

	//"io"

	"encoding/json"
	"log"
	"os"

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
	conf, err := readConf(configFile)
	if err != nil {
		log.Fatal(err)
	}

	channels := tickerToChannels(conf.Tickers)
	kafkaTopic := kafkaTopic(conf)
	kafkaHosts := kafkaHosts(conf)

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
		logrus.Fatal("Failed to create producer: %s\n", err)
	}

	fmt.Printf("Created Producer %v\n", p)

	_ = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", conf.Apikey)))
	_ = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", channels)))

	chanMessages := make(chan interface{}, 10000)
	//go processMessage(chanMessages, p, kafkaTopic)

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

	deliveryChan := make(chan kafka.Event)
	e := <-deliveryChan
	for msgBytes := range chanMessages {
		logrus.Info("Message Bytes: ", msgBytes)
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
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			//fmt.Printf()
			logrus.Fatal("Delivery failed: %v\n", m.TopicPartition.Error)
		}
		if err != nil {
			logrus.Fatal("Producer failure: %v\n", err)
		}
	}
}
