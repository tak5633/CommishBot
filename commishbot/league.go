package main

import (
	"encoding/json"
	"errors"
	"strconv"
)

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type League struct {
   Name string
   Sport string
   Season string
   League_id string

   Total_rosters int
   Roster_positions []string

   Scoring_settings map[string]json.RawMessage
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetUserLeaguesData(pUserId string, pYear int) string {
   return GetHttpResponse("https://api.sleeper.app/v1/user/" + pUserId + "/leagues/nfl/" + strconv.Itoa(pYear))
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetUserLeagues(pUserId string, pYear int) []League {
   userLeagueData := GetUserLeaguesData(pUserId, pYear)

   var leagues []League
   err := json.Unmarshal([]byte(userLeagueData), &leagues)
   check(err)

   return leagues
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetLeagueData(pLeagueId string) string {
   return GetHttpResponse("https://api.sleeper.app/v1/league/" + pLeagueId)
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetLeague(pLeagueId string) League {

   leagueData := GetLeagueData(pLeagueId)

   var league League
   err := json.Unmarshal([]byte(leagueData), &league)
   check(err)

   return league
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetScoringValue(pScoringSettings map[string]json.RawMessage, pScoringKey string) (float64, error) {

   keyExceptions := make(map[string]float64)
   keyExceptions["adp_dd_ppr"] = 0.0
   keyExceptions["cmp_pct"] = 0.0
   keyExceptions["def_fum_td"] = 0.0
   keyExceptions["def_kr_yd"] = 0.0
   keyExceptions["def_pr_td"] = 0.0
   keyExceptions["def_pr_yd"] = 0.0
   keyExceptions["fga"] = 0.0
   keyExceptions["gp"] = 0.0
   keyExceptions["pos_adp_dd_ppr"] = 0.0
   keyExceptions["pr"] = 0.0
   keyExceptions["pr_td"] = 0.0
   keyExceptions["pr_yd"] = 0.0
   keyExceptions["pts_half_ppr"] = 0.0
   keyExceptions["pts_ppr"] = 0.0
   keyExceptions["pts_std"] = 0.0
   keyExceptions["rec_tgt"] = 0.0
   keyExceptions["xpa"] = 0.0

   scoringValueData, hasKey := pScoringSettings[pScoringKey]

   if !hasKey {
      if exceptionValue, isException := keyExceptions[pScoringKey] ; isException {
         return exceptionValue, nil
      }

      return 0.0, errors.New("Failed to retrieve " + pScoringKey + " score setting")
   }

   var scoringValue float64
   err := json.Unmarshal([]byte(scoringValueData), &scoringValue)

   if err != nil {
      return 0.0, errors.New("Failed to unmarshal " + pScoringKey + " score setting")
   }

   return scoringValue, nil
}

