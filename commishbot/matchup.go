package main

import (
	"encoding/json"
	"errors"
	"math"
	"sort"
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

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (matchup Matchup) GetMaxRosterPoints(pPlayers map[string]Player, pStarterPositionCounts map[string]int) float64 {

   var maxStarterPoints []float64

   maxQbPoints, _ := matchup.getMaxPlayerPoints("QB", pPlayers, pStarterPositionCounts)
   maxStarterPoints = append(maxStarterPoints, maxQbPoints...)

   maxRbPoints, remainingRbPoints := matchup.getMaxPlayerPoints("RB", pPlayers, pStarterPositionCounts)
   maxStarterPoints = append(maxStarterPoints, maxRbPoints...)

   maxWrPoints, remainingWrPoints := matchup.getMaxPlayerPoints("WR", pPlayers, pStarterPositionCounts)
   maxStarterPoints = append(maxStarterPoints, maxWrPoints...)

   maxTePoints, remainingTePoints := matchup.getMaxPlayerPoints("TE", pPlayers, pStarterPositionCounts)
   maxStarterPoints = append(maxStarterPoints, maxTePoints...)

   maxKPoints, _ := matchup.getMaxPlayerPoints("K", pPlayers, pStarterPositionCounts)
   maxStarterPoints = append(maxStarterPoints, maxKPoints...)

   maxDefPoints, _ := matchup.getMaxPlayerPoints("DEF", pPlayers, pStarterPositionCounts)
   maxStarterPoints = append(maxStarterPoints, maxDefPoints...)

   var flexPoints []float64
   numFlexPositions := pStarterPositionCounts["FLEX"]

   flexRbPoints, _ := removeElements(remainingRbPoints, numFlexPositions)
   flexPoints = append(flexPoints, flexRbPoints...)

   flexWrPoints, _ := removeElements(remainingWrPoints, numFlexPositions)
   flexPoints = append(flexPoints, flexWrPoints...)

   flexTePoints, _ := removeElements(remainingTePoints, numFlexPositions)
   flexPoints = append(flexPoints, flexTePoints...)

   sort.Sort(sort.Reverse(sort.Float64Slice(flexPoints)))
   maxFlexPoints, _ := removeElements(flexPoints, numFlexPositions)
   maxStarterPoints = append(maxStarterPoints, maxFlexPoints...)

   maxPoints := 0.0

   for _, curMaxStarterPoints := range maxStarterPoints {
      maxPoints += curMaxStarterPoints
   }

   return maxPoints
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (matchup Matchup) getMaxPlayerPoints(pPosition string, pPlayers map[string]Player, pPositionCounts map[string]int) ([]float64, []float64) {

   playerPoints := matchup.getPlayerPoints(pPosition, pPlayers)
   sort.Sort(sort.Reverse(sort.Float64Slice(playerPoints)))

   maxPlayerPoints, remainingPlayerPoints := removeElements(playerPoints, pPositionCounts[pPosition])

   return maxPlayerPoints, remainingPlayerPoints
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (matchup Matchup) getPlayerPoints(pPosition string, pPlayers map[string]Player) []float64 {
   var playerPoints []float64
   players := matchup.getPlayers(pPosition, pPlayers)

   for _, player := range players {
      playerPoints = append(playerPoints, matchup.Players_points[player])
   }

   return playerPoints
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (matchup Matchup) getPlayers(pPosition string, pPlayers map[string]Player) []string {
   var playerIds []string

   for _, playerId := range matchup.Players {
      player := pPlayers[playerId]

      if player.Position == pPosition {
         playerIds = append(playerIds, playerId)
      }
   }

   return playerIds
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func removeElements(pSlice []float64, pNumElements int) ([]float64, []float64) {

   clampedNumElements := int(math.Min(float64(len(pSlice)), float64(pNumElements)))
   removedElements, remainingElements := pSlice[0:clampedNumElements], pSlice[clampedNumElements:]

   return removedElements, remainingElements
}
