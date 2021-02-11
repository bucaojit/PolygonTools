package main

type Trade struct {
  Symbol	string		`json:"sym"`
  Exchange	int		`json:"x"`
  Id		string		`json:"i"`
  Tape		int		`json:"z"`
  Price		float64		`json:"p"`
  Size		int		`json:"s"`
  Condition	int		`json:"c"`
  Timestamp	int		`json:"t"`
}

type Quote struct {
  Symbol	string		`json:"sym"`
  BidEx		int		`json:"bx"`
  BidPrice      float64         `json:"bp"`
  BidSize	int		`json:"bs"`
  AskEx		int		`json:"ax"`
  AskPrice	float64		`json:"ap"`
  AskSize	int		`json:"as"`
  Condition	int		`json:"c"`
  Timestamp	int		`json:"t"`
}

type MinAggregate struct {
  Symbol	string		`json:"sym"`
  Volume	int		`json:"v"`
  AccumVolume	int		`json:"av"`
  OpenDay	float64		`json:"op"`
  VWAP		float64		`json:"vw"`
  OpenTick	float64		`json:"o"`
  Closing	float64		`json:"c"`
  High		float64		`json:"h"`
  Low		float64		`json:"l"`
  TickVw	float64		`json:"a"`
  AvgSize	int		`json:"z"`
  StartTime	int		`json:"s"`
  EndTime	int		`json:"e"`
}
