package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
)

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type PlayerStats struct {
   Def_td float64
   Fum float64
   Fum_lost float64
   Ff float64
   Rec_td float64
   Rush_td float64
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetPlayerStatsData(pYear int, pWeek int) string {
   return GetHttpResponse("https://api.sleeper.app/v1/stats/nfl/regular/" + strconv.Itoa(pYear) + "/" + strconv.Itoa(pWeek))
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetPlayerStats(pYear int, pWeek int) map[string]PlayerStats {

   // TODO (tknack): Temporarily disable player stats retrieval to avoid stressing the server
   // playerStatsData := GetPlayerStatsData(pYear, pWeek)

   targetWeek := 12

   if pWeek != targetWeek {
      err := errors.New("Cannot get player stats for week " + strconv.Itoa(pWeek) + " since player stats are locked to week " + strconv.Itoa(targetWeek))
      panic(err)
   }

   playerStatsDataBytes, err := os.ReadFile("./Nfl.2023.Stats.Week" + strconv.Itoa(targetWeek) + ".json")
   check(err)
   playerStatsData := string(playerStatsDataBytes)

   playerStatsMap := make(map[string]PlayerStats)
   err = json.Unmarshal([]byte(playerStatsData), &playerStatsMap)
   check(err)

   return playerStatsMap
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetNumFumbles(pPlayerStats map[string]PlayerStats, pPlayerId string) float64 {

   starterStats := pPlayerStats[pPlayerId]

   return starterStats.Fum_lost
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetNumNonPassingTds(pPlayerStats map[string]PlayerStats, pPlayerId string) float64 {

   starterStats := pPlayerStats[pPlayerId]

   return starterStats.Rec_td + starterStats.Rush_td + starterStats.Def_td
}
