package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
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
func GetHttpResponse(pRequest string) string {
   resp, err := http.Get(pRequest)
   check(err)

   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   check(err)

   return string(body)
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type User struct {
   Username string
   User_id string
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetUserData(pUsername string) string {
   return GetHttpResponse("https://api.sleeper.app/v1/user/" + pUsername)
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetUser(pUsername string) User {
   userData := GetUserData(pUsername)

   var user User
   err := json.Unmarshal([]byte(userData), &user)
   check(err)

   return user
}

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
func GetTotalStarterPoints(pMatchup Matchup) float64 {
   totalStarterPoints := 0.0

   for _, starterPoints := range pMatchup.Starters_points {
      totalStarterPoints += starterPoints
   }

   return totalStarterPoints
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetBenchPlayers(pMatchup Matchup) []string {
   starterMap := make(map[string]int)

   for _, starter := range pMatchup.Starters {
      starterMap[starter] = 0
   }

   var benchPlayers []string

   for _, player := range pMatchup.Players {
      if _, isStarter := starterMap[player] ; !isStarter {
         benchPlayers = append(benchPlayers, player)
      }
   }

   return benchPlayers
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func GetBenchPlayerPoints(pMatchup Matchup) []float64 {
   var benchPlayerPoints []float64
   benchPlayers := GetBenchPlayers(pMatchup)

   for _, benchPlayer := range benchPlayers {
      benchPlayerPoints = append(benchPlayerPoints, pMatchup.Players_points[benchPlayer])
   }

   return benchPlayerPoints
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type Player struct {
   First_name string
   Last_name string
   Full_name string
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

   // playerStatsData := GetPlayerStatsData(pYear, pWeek)

   playerStatsDataBytes, err := os.ReadFile("./Nfl.2023.Stats.Week12.json")
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

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
type PrizeEntry struct {
   Score float64
   Owner string
}

type PrizeEntries []PrizeEntry

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (prizeEntries PrizeEntries) Len() int {
   return len(prizeEntries)
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (prizeEntries PrizeEntries) Less(i, j int) bool {
   return prizeEntries[i].Score < prizeEntries[j].Score
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (prizeEntries PrizeEntries) Swap(i, j int) {
   prizeEntries[i], prizeEntries[j] = prizeEntries[j], prizeEntries[i]
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (prizeEntries PrizeEntries) Reverse() {
   for i := len(prizeEntries)/2-1; i >= 0; i-- {
      opp := len(prizeEntries)-1-i
      prizeEntries[i], prizeEntries[opp] = prizeEntries[opp], prizeEntries[i]
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week1Summary(pLeagueInfo LeagueInfo) {

   week := 1
   criteria := "Hot Start - Highest Starting Team Score"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = GetTotalStarterPoints(matchupRoster)

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)
   prizeEntries.Reverse()

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Starter Points: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week2Summary(pLeagueInfo LeagueInfo) {

   week := 2
   criteria := "Dead Weight - Lowest Starting Player Score, Wins Matchup"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      matchupOpponentRoster, err :=  GetMatchupOpponentRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(1)

      totalStarterPoints := GetTotalStarterPoints(matchupRoster)
      totalOpponentStarterPoints := GetTotalStarterPoints(matchupOpponentRoster)

      if totalStarterPoints > totalOpponentStarterPoints {
         for _, starterPoints := range matchupRoster.Starters_points {
            prizeEntry.Score = math.Min(prizeEntry.Score, starterPoints)
         }
      }

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Starter Points: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week3Summary(pLeagueInfo LeagueInfo) {

   week := 3
   criteria := "MVP - Highest Starting Player Score"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(-1)

      for _, starterPoints := range matchupRoster.Starters_points {
         prizeEntry.Score = math.Max(prizeEntry.Score, starterPoints)
      }

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)
   prizeEntries.Reverse()

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Highest Starter Points: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week4Summary(pLeagueInfo LeagueInfo) {

   week := 4
   criteria := "Bench Warmers - Highest Team Bench Score"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = 0.0

      benchPlayerPoints := GetBenchPlayerPoints(matchupRoster)

      for _, benchPlayerPoints := range benchPlayerPoints {
         prizeEntry.Score += benchPlayerPoints
      }

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)
   prizeEntries.Reverse()

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Team Bench Score: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week5Summary(pLeagueInfo LeagueInfo) {

   week := 5
   criteria := "Biggest Loser - Highest Starting Team Score, Loses Matchup"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      matchupOpponentRoster, err :=  GetMatchupOpponentRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(-1)

      totalStarterPoints := GetTotalStarterPoints(matchupRoster)
      totalOpponentStarterPoints := GetTotalStarterPoints(matchupOpponentRoster)

      if totalStarterPoints < totalOpponentStarterPoints {
         prizeEntry.Score = totalStarterPoints
      }

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)
   prizeEntries.Reverse()

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Starter Points: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week6Summary(pLeagueInfo LeagueInfo) {

   week := 6
   criteria := "Photo Finish - Team With Closest Margin Of Victory"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      matchupOpponentRoster, err :=  GetMatchupOpponentRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(1)

      totalStarterPoints := GetTotalStarterPoints(matchupRoster)
      totalOpponentStarterPoints := GetTotalStarterPoints(matchupOpponentRoster)

      if totalStarterPoints > totalOpponentStarterPoints {
         prizeEntry.Score = totalStarterPoints - totalOpponentStarterPoints
      }

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Margin of Victory: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week7Summary(pLeagueInfo LeagueInfo) {

   week := 7
   criteria := "Biggest Blowout - Team With The Largest Margin of Victory"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      matchupOpponentRoster, err :=  GetMatchupOpponentRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(-1)

      totalStarterPoints := GetTotalStarterPoints(matchupRoster)
      totalOpponentStarterPoints := GetTotalStarterPoints(matchupOpponentRoster)

      if totalStarterPoints > totalOpponentStarterPoints {
         prizeEntry.Score = totalStarterPoints - totalOpponentStarterPoints
      }

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)
   prizeEntries.Reverse()

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Margin of Victory: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week12Summary(pLeagueInfo LeagueInfo, pYear int) {

   week := 12
   criteria := "Butterfingers - Most Starting Team Fumbles"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)
   playerStats := GetPlayerStats(pYear, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = 0.0

      for _, starter := range matchupRoster.Starters {
         numFumbles := GetNumFumbles(playerStats, starter)
         prizeEntry.Score += numFumbles
      }

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)
   prizeEntries.Reverse()

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Total Fumbles: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week13Summary(pLeagueInfo LeagueInfo) {

   week := 13
   criteria := "Blackjack - Staring Player Score Closest to 21 Without Going Over"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(-1)

      for _, starterPoints := range matchupRoster.Starters_points {
         if starterPoints <= 21.0 && prizeEntry.Score < starterPoints {
            prizeEntry.Score = starterPoints
         }
      }

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)
   prizeEntries.Reverse()

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Blackjack Points: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week14Summary(pLeagueInfo LeagueInfo, pYear int) {

   week := 14
   criteria := "Touchdown Dance - Team With The Most Touchdowns (Excludes QB Passing Touchdowns)"

   var prizeEntries PrizeEntries

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, week)
   playerStats := GetPlayerStats(pYear, week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         log.Printf("Week %d Summary: %s", week, err.Error())
         return
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = 0.0

      for _, starter := range matchupRoster.Starters {
         numNonPassingTds := GetNumNonPassingTds(playerStats, starter)
         prizeEntry.Score += numNonPassingTds
      }

      prizeEntries = append(prizeEntries, prizeEntry)
   }

   sort.Sort(prizeEntries)
   prizeEntries.Reverse()

   log.Printf("Week %d Criteria: %s", week, criteria)

   for _, prizeEntry := range prizeEntries {
      log.Printf("   Owner: %s, Total Touchdowns: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func main() {
   config := GetConfig("Config.json")
   log.Printf("%+v", config)

   user := GetUser(config.Username)
   userLeagues := GetUserLeagues(user.User_id, config.Year)

   if len(userLeagues) > 0 {

      leagueInfo := GetLeagueInfo(userLeagues[0].League_id)
      // players := GetPlayers()

      Week1Summary(leagueInfo)
      Week2Summary(leagueInfo)
      Week3Summary(leagueInfo)
      Week4Summary(leagueInfo)
      Week5Summary(leagueInfo)
      Week6Summary(leagueInfo)
      Week7Summary(leagueInfo)
      Week12Summary(leagueInfo, config.Year)
      Week13Summary(leagueInfo)
      Week14Summary(leagueInfo, config.Year)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func check(pE error) {
   if pE != nil {
      panic(pE)
   }
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
