package main

import (
	"encoding/json"
	"os"
)

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type Config struct {
   Username string
   Year int
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetConfig(pFilePath string) Config {
   file, _ := os.Open(pFilePath)
   defer file.Close()

   var config Config
   decoder := json.NewDecoder(file)

   err := decoder.Decode(&config)
   check(err)

   return config
}
