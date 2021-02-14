package main

import (
	"fmt"

	"github.com/bucaojit/PolygonTools/config"
	"github.com/sirupsen/logrus"
)

func main() {
	configFile := "../conf/polygon.yaml"

	/* smtp info
	 */

	conf, err := config.ReadConf(configFile)
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Println(conf)

}
