package models

import (
	"fmt"
	"strings"
)

var (
	dasAgDota2Structs = map[string][]string{
		"FeDota2Team": {
			"Logo",
		},
		"FeDota2Player": {
			"Logo",
		},
		"MatchBriefInfo_Team": {
			"TeamLogo",
		},
		"MatchBriefInfo_League": {
			"TeamLogo",
		},
		"FePlayerSpecialTitle_Title": {
			"PlayerLogo",
		},
		"FePlayerSpecialTitle_MVP": {
			"PlayerLogo",
		},
		"FightAnalysis_Team": {
			"TeamLogo",
		},
		"FightAnalysis_Player": {
			"PlayerLogo",
		},
		"MatchPlayerData_PlayerInfo": {
			"PlayerLogo",
		},
		"FeMatchPlayerSnapshot_PlayerSnapshot": {
			"PlayerLogo",
		},
		"TeamGankAnalysis_TeamGankStats": {
			"TeamLogo",
		},
		"TeamLocationAnalysis_HeroLocationInfo": {
			"PlayerLogo",
		},
		"TeamVisionAnalysis_TeamInfo": {
			"TeamLogo",
		},
		"TeamLaneAnalysis_PlayerInfo": {
			"TeamLogo",
		},
		"GetFeMatchTextBroadcastReply_MatchTextBroadcast_AccountInfo": {
			"Logo",
		},
		"GetFeMatchTextBroadcastReply_MatchTextBroadcast_TeamInfo": {
			"Logo",
		},
		"GetFeMatchTextBroadcastReply_MatchTextBroadcast_ItemInfo": {
			"Logo",
		},
	}
)


func IsExcludeInDasAgDota2(structName string, fieldName string) bool {
	fmt.Println("to check", structName, fieldName)
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