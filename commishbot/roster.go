package main

import "encoding/json"

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type Roster struct {
   Owner_id string
   Roster_id int
   Players []string
   Starters []string
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetRostersData(pLeagueId string) string {
   return GetHttpResponse("https://api.sleeper.app/v1/league/" + pLeagueId + "/rosters")
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetRosters(pLeagueId string) []Roster {

   rostersData := GetRostersData(pLeagueId)

   var rosters []Roster
   err := json.Unmarshal([]byte(rostersData), &rosters)
   check(err)

   return rosters
}
