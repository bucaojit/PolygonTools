package main

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type conf struct {
    Apikey string
    Database struct {
        Host string
        Port int
        Username string
        password string
    }
    Zookeeper string
    Tickers []string
}

func readConf(filename string) (*conf, error) {
    dat, err := ioutil.ReadFile(filename)
    check(err)

    conf := &conf{}
    err  = yaml.Unmarshal(dat, conf)
    if err != nil {
        return nil, fmt.Errorf("in file %q: %v", )
    }
    // fmt.Print(string(dat))
    return conf, nil
}