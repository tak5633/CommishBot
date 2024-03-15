package main

import "encoding/json"

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type LeagueUser struct {
   User_id string
   Display_name string
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetLeagueUsersData(pLeagueId string) string {
   return GetHttpResponse("https://api.sleeper.app/v1/league/" + pLeagueId + "/users")
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetLeagueUsers(pLeagueId string) []LeagueUser {

   leagueUsersData := GetLeagueUsersData(pLeagueId)

   var leagueUsers []LeagueUser
   err := json.Unmarshal([]byte(leagueUsersData), &leagueUsers)
   check(err)

   return leagueUsers
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func MakeDisplayNamesMap(pLeagueUsers []LeagueUser) map[string]string {

   displayNames := make(map[string]string)

   for _, leagueUser := range pLeagueUsers {
      displayNames[leagueUser.User_id] = leagueUser.Display_name
   }

   return displayNames
}
