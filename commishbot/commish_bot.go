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
type WeekSummary struct {
   Week int
   Criteria string
   PrizeEntries PrizeEntries
   Err error
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func (summary WeekSummary) Print() {

   if summary.Err != nil {
      log.Printf("Week %d Summary: %s", summary.Week, summary.Err.Error())
      return
   }

   log.Printf("Week %d Criteria: %s", summary.Week, summary.Criteria)

   for _, prizeEntry := range summary.PrizeEntries {
      log.Printf("   Owner: %s, Starter Points: %f", prizeEntry.Owner, prizeEntry.Score)
   }
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week1Summary(pLeagueInfo LeagueInfo) WeekSummary {

   var summary WeekSummary
   summary.Week = 1
   summary.Criteria = "Hot Start - Highest Starting Team Score"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = matchupRoster.GetTotalStarterPoints()

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)
   summary.PrizeEntries.Reverse()

   return summary
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week2Summary(pLeagueInfo LeagueInfo) WeekSummary {

   var summary WeekSummary
   summary.Week = 2
   summary.Criteria = "Dead Weight - Lowest Starting Player Score, Wins Matchup"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      matchupOpponentRoster, err :=  GetMatchupOpponentRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(1)

      totalStarterPoints := matchupRoster.GetTotalStarterPoints()
      totalOpponentStarterPoints := matchupOpponentRoster.GetTotalStarterPoints()

      if totalStarterPoints > totalOpponentStarterPoints {
         for _, starterPoints := range matchupRoster.Starters_points {
            prizeEntry.Score = math.Min(prizeEntry.Score, starterPoints)
         }
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)

   return summary
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week3Summary(pLeagueInfo LeagueInfo) WeekSummary {

   var summary WeekSummary
   summary.Week = 3
   summary.Criteria = "MVP - Highest Starting Player Score"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(-1)

      for _, starterPoints := range matchupRoster.Starters_points {
         prizeEntry.Score = math.Max(prizeEntry.Score, starterPoints)
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)
   summary.PrizeEntries.Reverse()

   return summary
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week4Summary(pLeagueInfo LeagueInfo) WeekSummary {

   var summary WeekSummary
   summary.Week = 4
   summary.Criteria = "Bench Warmers - Highest Team Bench Score"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = 0.0

      benchPlayerPoints := matchupRoster.GetBenchPlayerPoints()

      for _, benchPlayerPoints := range benchPlayerPoints {
         prizeEntry.Score += benchPlayerPoints
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)
   summary.PrizeEntries.Reverse()

   return summary
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week5Summary(pLeagueInfo LeagueInfo) WeekSummary {

   var summary WeekSummary
   summary.Week = 5
   summary.Criteria = "Biggest Loser - Highest Starting Team Score, Loses Matchup"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      matchupOpponentRoster, err :=  GetMatchupOpponentRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(-1)

      totalStarterPoints := matchupRoster.GetTotalStarterPoints()
      totalOpponentStarterPoints := matchupOpponentRoster.GetTotalStarterPoints()

      if totalStarterPoints < totalOpponentStarterPoints {
         prizeEntry.Score = totalStarterPoints
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)
   summary.PrizeEntries.Reverse()

   return summary
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week6Summary(pLeagueInfo LeagueInfo) WeekSummary {

   var summary WeekSummary
   summary.Week = 6
   summary.Criteria = "Photo Finish - Team With Closest Margin Of Victory"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      matchupOpponentRoster, err :=  GetMatchupOpponentRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(1)

      totalStarterPoints := matchupRoster.GetTotalStarterPoints()
      totalOpponentStarterPoints := matchupOpponentRoster.GetTotalStarterPoints()

      if totalStarterPoints > totalOpponentStarterPoints {
         prizeEntry.Score = totalStarterPoints - totalOpponentStarterPoints
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)

   return summary
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week7Summary(pLeagueInfo LeagueInfo) WeekSummary {

   var summary WeekSummary
   summary.Week = 7
   summary.Criteria = "Biggest Blowout - Team With The Largest Margin of Victory"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      matchupOpponentRoster, err :=  GetMatchupOpponentRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(-1)

      totalStarterPoints := matchupRoster.GetTotalStarterPoints()
      totalOpponentStarterPoints := matchupOpponentRoster.GetTotalStarterPoints()

      if totalStarterPoints > totalOpponentStarterPoints {
         prizeEntry.Score = totalStarterPoints - totalOpponentStarterPoints
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)
   summary.PrizeEntries.Reverse()

   return summary
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week12Summary(pLeagueInfo LeagueInfo, pYear int) WeekSummary {

   var summary WeekSummary
   summary.Week = 12
   summary.Criteria = "Butterfingers - Most Starting Team Fumbles"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)
   playerStats := GetPlayerStats(pYear, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = 0.0

      for _, starter := range matchupRoster.Starters {
         numFumbles := GetNumFumbles(playerStats, starter)
         prizeEntry.Score += numFumbles
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)
   summary.PrizeEntries.Reverse()

   return summary
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week13Summary(pLeagueInfo LeagueInfo) WeekSummary {

   var summary WeekSummary
   summary.Week = 13
   summary.Criteria = "Blackjack - Staring Player Score Closest to 21 Without Going Over"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = math.Inf(-1)

      for _, starterPoints := range matchupRoster.Starters_points {
         if starterPoints <= 21.0 && prizeEntry.Score < starterPoints {
            prizeEntry.Score = starterPoints
         }
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)
   summary.PrizeEntries.Reverse()

   return summary
}

//--------------------------------------------------------------------------------------------------
//
//--------------------------------------------------------------------------------------------------
func Week14Summary(pLeagueInfo LeagueInfo, pYear int) WeekSummary {

   var summary WeekSummary
   summary.Week = 14
   summary.Criteria = "Touchdown Dance - Team With The Most Touchdowns (Excludes QB Passing Touchdowns)"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)
   playerStats := GetPlayerStats(pYear, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = 0.0

      for _, starter := range matchupRoster.Starters {
         numNonPassingTds := GetNumNonPassingTds(playerStats, starter)
         prizeEntry.Score += numNonPassingTds
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)
   summary.PrizeEntries.Reverse()

   return summary
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

      Week1Summary(leagueInfo).Print()
      Week2Summary(leagueInfo).Print()
      Week3Summary(leagueInfo).Print()
      Week4Summary(leagueInfo).Print()
      Week5Summary(leagueInfo).Print()
      Week6Summary(leagueInfo).Print()
      Week7Summary(leagueInfo).Print()
      Week12Summary(leagueInfo, config.Year).Print()
      Week13Summary(leagueInfo).Print()
      Week14Summary(leagueInfo, config.Year).Print()
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
