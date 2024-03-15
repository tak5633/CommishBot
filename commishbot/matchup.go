package main

import (
	"encoding/json"
	"errors"
	"strconv"
)

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type Matchup struct {
   Matchup_id int
   Roster_id int
   Starters []string
   Starters_points []float64
   Players []string
   Players_points map[string]float64
   Points float64
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetMatchupsData(pLeagueId string, pWeek int) string {
   return GetHttpResponse("https://api.sleeper.app/v1/league/" + pLeagueId + "/matchups/" + strconv.Itoa(pWeek))
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetMatchups(pLeagueId string, pWeek int) []Matchup {

   matchupsData := GetMatchupsData(pLeagueId, pWeek)

   var matchups []Matchup
   err := json.Unmarshal([]byte(matchupsData), &matchups)
   check(err)

   return matchups
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetMatchupRoster(pMatchups []Matchup, pRosterId int) (Matchup, error) {

   for _, matchup := range pMatchups {
      if matchup.Roster_id == pRosterId {
         return matchup, nil
      }
   }

   return Matchup{}, errors.New("GetMatchupRoster: Failed to find roster (Id: " + strconv.Itoa(pRosterId) + ")")
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetMatchupOpponentRoster(pMatchups []Matchup, pRosterId int) (Matchup, error) {

   matchupRoster, err := GetMatchupRoster(pMatchups, pRosterId)

   if err != nil {
      return Matchup{}, err
   }

   for _, matchup := range pMatchups {
      if matchup.Matchup_id == matchupRoster.Matchup_id && matchup.Roster_id != matchupRoster.Roster_id {
         return matchup, nil
      }
   }

   return Matchup{}, errors.New("GetMatchupOpponentRoster: Failed to find opponent roster (Id: " + strconv.Itoa(pRosterId) + ")")
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (matchup Matchup) GetTotalStarterPoints() float64 {
   totalStarterPoints := 0.0

   for _, starterPoints := range matchup.Starters_points {
      totalStarterPoints += starterPoints
   }

   return totalStarterPoints
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (matchup Matchup) GetStarterPlayerPoints() map[string]float64 {
   starterMap := make(map[string]float64)

   for _, starter := range matchup.Starters {
      starterMap[starter] = matchup.Players_points[starter]
   }

   return starterMap
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (matchup Matchup) GetBenchPlayers() []string {
   starterMap := make(map[string]int)

   for _, starter := range matchup.Starters {
      starterMap[starter] = 0
   }

   var benchPlayers []string

   for _, player := range matchup.Players {
      if _, isStarter := starterMap[player] ; !isStarter {
         benchPlayers = append(benchPlayers, player)
      }
   }

   return benchPlayers
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (matchup Matchup) GetBenchPlayerPoints() []float64 {
   var benchPlayerPoints []float64
   benchPlayers := matchup.GetBenchPlayers()

   for _, benchPlayer := range benchPlayers {
      benchPlayerPoints = append(benchPlayerPoints, matchup.Players_points[benchPlayer])
   }

   return benchPlayerPoints
}
