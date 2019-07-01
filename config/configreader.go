package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Contains the configurable options of the microservice
var Configuration Config

// Encapsulates the fields of the configuration file
type Config struct {
	ListeningAddress  string
	CoursesTableName  string
	MessageQueueName  string
	PollingWaitTime   int64
	AwsSesRegion      string
	AwsSqsRegion      string
	AwsDynamoDbRegion string
	MailTemplate      string
	MailAddress       string
}

func SetConfiguration(configFile string) error {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &Configuration)
	if err != nil {
		return err
	}
	return nil
}
