package configure

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Bachelor-project-f20/eventToGo"
	etgNats "github.com/Bachelor-project-f20/eventToGo/nats"
	nats "github.com/nats-io/nats.go"
)

var (
	MessageBrokerTypeDefault       = eventToGo.NATS
	MessageBrokerDefaultConnection = "localhost:4222"
	ExchangeDefault                = "user"
	QueueTypeDefault               = "queue"
)

type ServiceConfig struct {
	MessageBrokerType       eventToGo.BrokerType
	MessageBrokerConnection string
	Exchange                string
	QueueType               string
	EventEmitter            eventToGo.EventEmitter
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	config := ServiceConfig{
		MessageBrokerTypeDefault,
		MessageBrokerDefaultConnection,
		ExchangeDefault,
		QueueTypeDefault,
		nil,
	}
	config.extractFromFile(filename)
	config.extractFromEnv()

	err := config.setupMessageBroker()
	if err != nil {
		return ServiceConfig{}, err
	}

	return config, nil
}

func (config *ServiceConfig) extractFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Config file not found")
	}
	json.NewDecoder(file).Decode(config)
}

func (config *ServiceConfig) extractFromEnv() {
	if v := os.Getenv("NATS_BROKER_URL"); v != "" {
		config.MessageBrokerType = eventToGo.NATS
		config.MessageBrokerConnection = v
	}
}

func (config *ServiceConfig) setupMessageBroker() error {
	var eventEmitter eventToGo.EventEmitter
	if config.MessageBrokerType == eventToGo.NATS {
		log.Println("Setting up nats")
		encodedConn, err := config.setupNatsConn()
		if err != nil {
			log.Fatalf("Error connecting to Nats: %v \n", err)
			return err
		}
		eventEmitter, err = etgNats.NewNatsEventEmitter(encodedConn, config.Exchange, config.QueueType)
		if err != nil {
			log.Fatalf("Error creating Emitter: %v \n", err)
		}
	}
	log.Println("Event emitter ready")
	config.EventEmitter = eventEmitter
	return nil
}

func (config *ServiceConfig) setupNatsConn() (*nats.EncodedConn, error) {

	natsConn, err := nats.Connect(config.MessageBrokerConnection)

	if err != nil {
		fmt.Println("Connection to Nats failed")
		return nil, err
	}

	encodedConn, err := nats.NewEncodedConn(natsConn, "json")

	if err != nil {
		fmt.Println("Creation of encoded connection failed")
		return nil, err
	}

	return encodedConn, nil

}
