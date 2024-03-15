package main

import (
	"encoding/json"
	"errors"
	"strconv"
)

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetProjectedPlayerStatsData(pPlayerId string, pYear int) string {
   return GetHttpResponse("https://api.sleeper.com/projections/nfl/player/" + pPlayerId + "?season_type=regular&season=" + strconv.Itoa(pYear) + "&grouping=week")
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetProjectedPlayerStats(pPlayerId string, pYear int) map[string]json.RawMessage {

   projectedPlayerStatsData := GetProjectedPlayerStatsData(pPlayerId, pYear)

   var projectedPlayerStats map[string]json.RawMessage
   err := json.Unmarshal([]byte(projectedPlayerStatsData), &projectedPlayerStats)
   check(err)

   return projectedPlayerStats
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetProjectedPlayerWeekStats(pPlayerId string, pYear int, pWeek int) (map[string]float64, error) {

   yearStr := strconv.Itoa(pYear)
   weekStr := strconv.Itoa(pWeek)

   projectedPlayerStats := GetProjectedPlayerStats(pPlayerId, pYear)
   projectedWeekData, hasKey := projectedPlayerStats[weekStr]

   if !hasKey {
      return nil, errors.New("Failed to retrieve " + yearStr + " week " + weekStr + " projections for player Id " + pPlayerId)
   }

   var projectedWeek map[string]json.RawMessage
   err := json.Unmarshal([]byte(projectedWeekData), &projectedWeek)

   if err != nil {
      return nil, errors.New("Failed to unmarshal " + yearStr + " week " + weekStr + " projections for player Id " + pPlayerId)
   }

   projectedWeekStatsData, hasKey := projectedWeek["stats"]

   if !hasKey {
      return nil, errors.New("Failed to retrieve " + yearStr + " week " + weekStr + " stat projections for player Id " + pPlayerId)
   }

   var projectedWeekStats map[string]float64
   err = json.Unmarshal([]byte(projectedWeekStatsData), &projectedWeekStats)

   if err != nil {
      return nil, errors.New("Failed to unmarshal " + yearStr + " week " + weekStr + " stat projections for player Id " + pPlayerId)
   }

   return projectedWeekStats, nil
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetProjectedPlayerWeekScore(pPlayerId string, pYear int, pWeek int, pScoringSettings map[string]json.RawMessage) (float64, error) {

   projectedWeekStats, err := GetProjectedPlayerWeekStats(pPlayerId, pYear, pWeek)

   if err != nil {
      return 0.0, err
   }

   starterProjection := 0.0

   for statKey, statValue := range projectedWeekStats {

      scoringValue, err := GetScoringValue(pScoringSettings, statKey)

      if err != nil {
         return 0.0, err
      }

      starterProjection += scoringValue * statValue
   }

   return starterProjection, nil
}
