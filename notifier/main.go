package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"os"
	"strconv"

	"github.com/bucaojit/PolygonTools/config"
	"github.com/sirupsen/logrus"
)

/*
	Notifier service

	Use to send emails and text messages (email to phone number)
	Uses a yaml configuration file to specify:
	* smtp information
	* tls cert and key files
	* port to listen on

	Takes a JSON payload with the format:
	{
		"Message" : "Message to Send",
		"Recipients" : [
			"user@email.com"
		]
	}
*/

type mail struct {
	Message    string
	Recipients []string
}

func main() {
	configFile := "../conf/polygon.yaml"
	args := os.Args[1:]
	if len(args) == 1 {
		configFile = os.Args[1]
	}

	conf, err := config.ReadConf(configFile)
	if err != nil {
		logrus.Fatal(err)
	}

	createServer(conf)

}

func createServer(conf *config.Conf) {
	serverPort := strconv.Itoa(conf.Notifierport)
	handler := http.NewServeMux()
	handler.Handle("/email", sendMail(conf))
	server := &http.Server{
		Handler: handler,
		Addr:    ":" + serverPort,
	}

	logrus.Fatal(server.ListenAndServeTLS(conf.Tls.Certfile, conf.Tls.Keyfile))

}

func sendMail(conf *config.Conf) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		r.ParseForm()
		switch r.Method {
		case http.MethodPost:
			var mailToSend mail
			err := json.NewDecoder(r.Body).Decode(&mailToSend)
			if err != nil {
				logrus.Warn(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			auth := smtp.PlainAuth("",
				conf.Smtp.User,
				conf.Smtp.Password,
				conf.Smtp.Server)
			fmt.Println(auth)
			serverAndPort := conf.Smtp.Server + ":" + strconv.Itoa(conf.Smtp.Port)

			fmt.Println(mailToSend.Recipients, "--", []byte(mailToSend.Message))
			err = smtp.SendMail(serverAndPort,
				auth,
				conf.Smtp.User,
				mailToSend.Recipients,
				[]byte(mailToSend.Message))
			if err != nil {
				logrus.Warn(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		default:
			errMsg := "Unsupported method"
			logrus.Error(errMsg)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
	})
}
