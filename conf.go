package main

import (
	"os"
	"flag"
	"sync"
	"time"
	"strings"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"github.com/olling/logger"
)

type Configuration struct {
	PathInstall string
	PathWork string
	PathConfiguration string
	PathLog string
	HttpPort int
	HttpTlsCert string
	HttpTlsKey string
	HttpTlsPort int
	Debug bool
	Timeout int
}

type Action struct {
	FriendlyName string
	Name string
	Cron string
	Description string
	children []Action
	Command string
	Timeout int
	LastRun time.Time
	Enabled bool
	fileName string
}

var (
	ConfigurationPath = flag.String("configurationpath","/etc/repod/repod.conf","(Optional) The path to the configuration file")
        outputMutex sync.Mutex
	CurrentConfiguration Configuration
)

func InitializeConfiguration() {
	logger.Debug("Initializing configuration")
	flag.Parse()
	logger.Debug("Configuration Path: " + *ConfigurationPath)
	err := ReadJsonFile(*ConfigurationPath,&CurrentConfiguration)

        if err != nil {
                logger.Error("Error while reading the configuration file - Exiting")
                logger.Error(err)
                os.Exit(1)
        }

	logger.Debug("Done initializing configuration")
}

//func LoadChannelWithoutChildrenFromPath(path string) (outputChannel Channel,err error) {
//	err = ReadJsonFile(path + "/channel_info",&outputChannel)
//	return outputChannel,err
//}


func (action *Action) GetDirName () (dirName string) {
	dirName = strings.TrimSuffix(action.fileName,filepath.Ext(action.fileName))
	return dirName
}


func LoadActionsFromPath(path string) (actions []Action, err error) {
	logger.Debug("LoadActionsFromPath(" + path + ")")
	items, err := ioutil.ReadDir(path)

	if err != nil {
		return actions,err
	}

	for _, item := range items {
		logger.Debug("Reading path:",item.Name())

		if item.IsDir() == false {
			fullPath :=path + "/" + item.Name()
			logger.Debug("Trying to get action from:", item.Name())
			var a Action
			a.fileName = item.Name()
			err := ReadJsonFile(fullPath,&a)
                        if err != nil {
                                logger.Error("Could not get info from:", fullPath, err)
                        }

			pathChildren := strings.TrimSuffix(fullPath,filepath.Ext(fullPath))
			_, err = os.Stat(fullPath)
			if err != nil && os.IsNotExist(err) {
				logger.Debug("Could not find the action children directory:",pathChildren)
			} else {
				logger.Debug("Reading children from: ", pathChildren)
				children,_ := LoadActionsFromPath(pathChildren)
				a.children = append(a.children, children...)
			}
			actions = append(actions, a)
		}
	}
	return actions,nil
}



func LoadActionFromPath(path string) (action Action,err error) {
	logger.Debug("LoadActionFromPath(" + path + ")")

	err = ReadJsonFile(path,&action)

	if err != nil {
		logger.Error("Could not load json from:", path, err)
		return action,err
	}

	pathChildren := strings.TrimSuffix(path,filepath.Ext(path))
	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		logger.Debug("Could not find the action children directory:",pathChildren)
	} else {
		logger.Debug("Reading children from: ", pathChildren)
		children,_ := LoadActionsFromPath(pathChildren)
		action.children = children
	}

	return action,nil
}

func LoadActionWithoutChildrenFromPath(path string) (action Action,err error) {
	logger.Debug("Reading action from", path)
	err = ReadJsonFile(path,&action)
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
