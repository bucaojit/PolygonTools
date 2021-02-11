# PolygonTools
Tools to stream and access financial market data from Polygon.io to Kafka

### Run
```
$ cd Streamer/  
$ go get -d ./...  
$ go build  
$ ./Streamer [optional path to config]  
```

### Configuration
default: conf/polygon.yaml

```
apikey: <API Key from Polygon>
database:
  host: mongodb
  port: port
  username: user
  password: password
zookeeper: zk:port
kafka:
  topic: polygon
  host: hostname
  port: port
tickers:
  - SPY
  - AAPL
  - IWM
```
