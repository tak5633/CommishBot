package main

import "log"

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

   for idx, prizeEntry := range summary.PrizeEntries {
      log.Printf("   %d. Owner: %s, Score: %f", idx+1, prizeEntry.Owner, prizeEntry.Score)
   }
}
