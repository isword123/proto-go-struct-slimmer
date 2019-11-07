package models

import "strings"

var (
	dasAnalysisFeStructs = map[string][]string{
		"FeDota2Hero": {
			"Logo",
		},
		"MatchStatsValPoint": {
			"Value",
		},
		"FeDota2Team": {
			"Logo",
		},
		"FeDota2Player": {
			"Logo",
		},
		"LeagueSimpleInfo": {
			"Logo",
		},
		"ScheduleSimpleInfo_FeDota2TeamWithScore": {
			"Logo",
		},
		"FeLOLTeam": {
			"Logo",
		},
		"FeCSGOTeam": {
			"Logo",
		},
		"TeamVsTeamStats_VsStats": {
			"Matches",
			"FirstBlood",
			"TenKills",
		},
	}
)

func IsExcludeInDasAnalysis(structName string, fieldName string) bool {
	arr, ok := 	dasAgDota2Structs[structName]
	if !ok {
		return false
	}


	for _, name := range arr {
		if strings.EqualFold(name, fieldName) {
			return true
		}
	}

	return false
}
