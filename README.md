# PolygonTools
Tools to stream and access financial market data from Polygon.io to Kafka

## Streamer
### Run
```
$ cd streamer/  
$ go get -d ./...  
$ go build  
$ ./streamer [optional path to config]  
```

## Configuration
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
