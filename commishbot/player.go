package main

import (
	"encoding/json"
	"os"
)

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type Player struct {
   First_name string
   Last_name string
   Full_name string
   Position string
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetPlayerData() string {
   return GetHttpResponse("https://api.sleeper.app/v1/players/nfl")
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetPlayers() map[string]Player {

   // TODO (tknack): Temporarily disable player data retrieval to avoid stressing the server
   // playerData := GetPlayerData()
   playerDataBytes, err := os.ReadFile("./Nfl.2023.Players.json")
   check(err)

   playerData := string(playerDataBytes)

   playerMap := make(map[string]Player)
   err = json.Unmarshal([]byte(playerData), &playerMap)
   check(err)

   return playerMap
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetPlayerName(pPlayers map[string]Player, pPlayerId string) string {

   if player, hasKey := pPlayers[pPlayerId] ; hasKey {
      if player.Full_name != "" {
         return player.Full_name
      }

      return player.First_name + " " + player.Last_name
   }

   return pPlayerId
}
