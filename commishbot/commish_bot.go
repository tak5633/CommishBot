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
func Week10Summary(pLeagueInfo LeagueInfo, pYear int) WeekSummary {

   var summary WeekSummary
   summary.Week = 10
   summary.Criteria = "Overachiver - Team With The Most Points Over Their Weekly Projection"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      starterPlayerPoints := matchupRoster.GetStarterPlayerPoints()

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = 0.0

      for _, starter := range matchupRoster.Starters {

         starterProjection, err := GetProjectedPlayerWeekScore(starter, pYear, summary.Week, pLeagueInfo.mLeague.Scoring_settings)

         if err != nil {
            summary.Err = err
            return summary
         }

         starterPoints, hasStarterPoints := starterPlayerPoints[starter]

         if !hasStarterPoints {
            summary.Err = errors.New("Failed to retrieve player " + starter + " points")
            return summary
         }

         prizeEntry.Score += (starterPoints - starterProjection)
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
func Week11Summary(pLeagueInfo LeagueInfo, pYear int) WeekSummary {

   var summary WeekSummary
   summary.Week = 11
   summary.Criteria = "Underperformer - Team With The Most Points Under Their Weekly Projection"

   matchups := GetMatchups(pLeagueInfo.mLeague.League_id, summary.Week)

   for _, roster := range pLeagueInfo.mRosters {

      matchupRoster, err := GetMatchupRoster(matchups, roster.Roster_id)

      if err != nil {
         summary.Err = err
         return summary
      }

      starterPlayerPoints := matchupRoster.GetStarterPlayerPoints()

      var prizeEntry PrizeEntry
      prizeEntry.Owner = pLeagueInfo.mDisplayNames[roster.Owner_id]
      prizeEntry.Score = 0.0

      for _, starter := range matchupRoster.Starters {

         starterProjection, err := GetProjectedPlayerWeekScore(starter, pYear, summary.Week, pLeagueInfo.mLeague.Scoring_settings)

         if err != nil {
            summary.Err = err
            return summary
         }

         starterPoints, hasStarterPoints := starterPlayerPoints[starter]

         if !hasStarterPoints {
            summary.Err = errors.New("Failed to retrieve player " + starter + " points")
            return summary
         }

         prizeEntry.Score += (starterPoints - starterProjection)
      }

      summary.PrizeEntries = append(summary.PrizeEntries, prizeEntry)
   }

   sort.Sort(summary.PrizeEntries)

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
      Week10Summary(leagueInfo, config.Year).Print()
      Week11Summary(leagueInfo, config.Year).Print()
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
