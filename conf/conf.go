package conf

import (
	"os"
	"time"
	"sync"
	"io/ioutil"
	"encoding/json"
	"github.com/olling/logger"
)

type Configuration struct {
	PathInstall string
	PathWork string
	PathConfiguration string
	PathLog string
	TlsCert string
	TlsKey string
	TlsPort int
	Debug bool
	Timeout int
}

type Action struct {
	FriendlyName string
	Name string
	Cron string
	Description string
	channels []Channel
	actions []Action
	Command string
	Timeout int
	LastRun time.Time
	Enabled bool
}

type Channel struct {
	FriendlyName string
	Name string
	Cron string
	Description string
	Channels []Channel
	Actions []Action
}

var (
        outputMutex sync.Mutex
	CurrentConfiguration Configuration
)

func InitializeConfiguration(path string) {
	logger.Debug("Initializing configuration")
	logger.Debug("Configuration Path: " + path)
	err := ReadJsonFile(path,&CurrentConfiguration)

        if err != nil {
                logger.Error("Error while reading the configuration file - Exiting")
                logger.Error(err)
                os.Exit(1)
        }

	logger.Debug("Done initializing configuration")
}

func LoadChannelsFromPath(path string) (outputChannel Channel,err error) {
	items, err := ioutil.ReadDir(path)

	if err != nil {
		logger.Debug(err)
		return outputChannel,err
	}

	for _, item := range items {
		logger.Debug("Reading path:",item.Name())

		if item.Name() == "channel_info" {
			continue
		}

		if item.IsDir() {
			var c Channel
			err := ReadJsonFile(path + "/" + item.Name() + "/channel_info",&c)
			if err != nil {
				logger.Error("Could not get Channel info from:", path + "/" + item.Name() + "/channel_info", err)
			}

			child,loaderr := LoadChannelsFromPath(path + "/" + item.Name())
			if loaderr != nil {
				logger.Error("Could not get children from:",path + "/" + item.Name(),loaderr)
			}
			c.Channels = child.Channels
			c.Actions = child.Actions

			outputChannel.Channels = append(outputChannel.Channels,c)
		} else {
			var a Action
			err := ReadJsonFile(path + "/" + item.Name(),&a)
                        if err != nil {
                                logger.Error("Could not get Action info from:", path + "/" + item.Name(), err)
                        }

                        outputChannel.Actions = append(outputChannel.Actions,a)
		}
	}
	return outputChannel,nil
}

func LoadActionFromPath(path string) (action Action,err error) {
	logger.Debug("Reading action from", path)
	err = ReadJsonFile(path,action)
	return action, err
}

func ReadJsonFile (path string, output interface{}) (error) {
	logger.Debug("Reading JSON file:", path)
        file,fileerr := os.Open(path)
	if fileerr != nil {
		logger.Debug("Error opening JSON file for reading: ", fileerr)
		return fileerr
	}

        decoder := json.NewDecoder(file)
        decodererr := decoder.Decode(&output)
	if decodererr != nil {
		logger.Debug("Error decoding JSON file", decodererr)
		return decodererr
	}

	logger.Debug("Done reading JSON file:", path, output)
	return nil
}


func WriteJsonFile(s interface{}, path string) (err error){
	logger.Debug("Writing to path: " + path, "content:", s) 
        outputMutex.Lock()
        defer outputMutex.Unlock()

        bytes, marshalErr := json.MarshalIndent(s,"","\t")
        if marshalErr != nil {
                logger.Error("Could not convert struct to bytes", marshalErr)
                return marshalErr
        }

	err = ioutil.WriteFile(path, bytes, 0644)
	if err != nil {
		logger.Error("Could not write file: " + path + ". Error:", err)
		return err
	}
	return nil
}
