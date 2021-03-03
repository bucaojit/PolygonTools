package main

import (
	"time"
	"fmt"
	"os"
	"log"
	"testing"
	"github.com/bucaojit/PolygonTools/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func setup() (*kafka.Producer, error){
	configFile := "../conf/polygon.yaml"
	envConfigFile := os.Getenv("CONFIGFILE")
	if envConfigFile != "" {
		configFile = envConfigFile
	}
	conf, err := config.ReadConf(configFile)
	if err != nil {
		log.Fatal(err)
	}
	hosts := config.KafkaHosts(conf)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": hosts})
	return p, err
}

/*
*   Testing Producer function to Kafka topic: testTopic
*   Start a consumer prior to running test to receive the test message
*/
func TestWriteToKafka(t *testing.T) {
	fmt.Println("Testing Write to Kafka")

	if producer, err := setup(); err != nil {
		t.Errorf("An error occurred during setup")
	} else {
		currentTime := time.Now()
		chanMessages := make(chan interface{}, 100)
		go processMessage(chanMessages, producer, "testTopic")
		chanMessages <- currentTime.Format("2006-01-01 03:04:05 PM") + ": TestWriteToKafka Success"
		time.Sleep(100 * time.Millisecond)
		close(chanMessages)
	}
}
