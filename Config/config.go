package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Conf struct {
	Apikey   string
	Database struct {
		Host     string
		Port     int
		Username string
		password string
	}
	Kafka struct {
		Hosts []string
		Topic string
	}
	Zookeeper string
	Tickers   []string
}

// func setupWebsocket()
func tickerToChannels(tickers []string) string {
	channels := ""
	for _, ticker := range tickers {
		channels += "T." + ticker + ","
		channels += "Q." + ticker + ","
	}
	channels = channels[:len(channels)-1]
	return channels
}

func kafkaHosts(config *Conf) string {
	retString := ""
	for _, entry := range config.Kafka.Hosts {
		retString += entry + ","
	}
	retString = retString[:len(retString)-1]
	return retString
}

func kafkaTopic(config *Conf) string {
	return config.Kafka.Topic
}

func readConf(filename string) (*Conf, error) {
	dat, err := ioutil.ReadFile(filename)
	check(err)

	conf := &Conf{}
	err = yaml.Unmarshal(dat, conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v")
	}
	// fmt.Print(string(dat))
	return conf, nil
}
