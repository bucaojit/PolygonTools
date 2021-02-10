package main

import (
    //"bufio"
    "fmt"
    //"io"
    
    
    "log"
    "os"
    "github.com/sirupsen/logrus"
    "github.com/gorilla/websocket"
)


// func setupWebsocket() 
func tickerToChannels(tickers []string) string {
    channels := ""
    for _, ticker := range tickers {
        channels += "T." + ticker + ","
        channels += "Q." + ticker + ","
    }
    channels = channels[:len(channels) - 1]
    return channels
}

func main() {
    configFile := "../conf/polygon.yaml"
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

    CHANNELS := tickerToChannels(conf.Tickers)
    fmt.Println(CHANNELS)

    c, _, err := websocket.DefaultDialer.Dial("wss://socket.polygon.io/stocks", nil)
	if err != nil {
		panic(err)
    }
    defer c.Close()

	_ = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"action\":\"auth\",\"params\":\"%s\"}", conf.Apikey)))
	_ = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("{\"action\":\"subscribe\",\"params\":\"%s\"}", CHANNELS)))

	// Buffered channel to account for bursts or spikes in data:
	chanMessages := make(chan interface{}, 10000)

	// Read messages off the buffered queue:
	go func() {
		for msgBytes := range chanMessages {
			logrus.Info("Message Bytes: ", msgBytes)
		}
	}()

	// As little logic as possible in the reader loop:
	for {
		var msg interface{}
		err := c.ReadJSON(&msg)
		// Ideally use c.ReadMessage instead of ReadJSON so you can parse the JSON data in a
		// separate go routine. Any processing done in this loop increases the chances of disconnects
		// due to not consuming the data fast enough.
		if err != nil {
			panic(err)
		}
		chanMessages <- msg
	}

}
