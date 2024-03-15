package main

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type LeagueInfo struct {
   mLeague League
   mLeagueUsers []LeagueUser
   mDisplayNames map[string]string
   mRosters []Roster
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetLeagueInfo(pLeagueId string) LeagueInfo {
   var leagueInfo LeagueInfo

   leagueInfo.mLeague = GetLeague(pLeagueId)
   leagueInfo.mLeagueUsers = GetLeagueUsers(pLeagueId)
   leagueInfo.mRosters = GetRosters(pLeagueId)

   leagueInfo.mDisplayNames = MakeDisplayNamesMap(leagueInfo.mLeagueUsers)

   return leagueInfo
}
