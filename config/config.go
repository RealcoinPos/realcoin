package config

import (
    "encoding/json"
    "os"
    "fmt"
)

type Configuration struct {
	Masterserver	string	`json:"masterServer"`
	Netname			string	`json:"netname"`
    Netset			string	`json:"netset"`
	Nethost			string	`json:"nethost"`
    Netport			int		`json:"netport"`
	Blockperiod		int		`json:"blockperiod"`
    Maxnumber		int		`json:"maxnumber"`
}

func SetConfiguration(File string) Configuration {
    var Config Configuration
    configFile, err := os.Open(File)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&Config)
	if Config.Netname=="realcoin" {
 		Config.Netset="pos"
		Config.Blockperiod=300
 		Config.Maxnumber=1000
	} else {
		if Config.vaildConfiguration()==false {
	        panic("Please confirm configuration file.")
		}
	}
    return Config
}

func (c Configuration) vaildConfiguration() bool {
	if c.Netname=="" || c.Netset=="" || c.Nethost=="" || c.Netport<0 {
        fmt.Println("Please confirm network infomation.")
		return false
	}
	if c.Blockperiod<10 {
        fmt.Println("Please confirm blockperiod in configuration file.")
		return false
	}
	return true
}
