package models

import "strings"

func IsExcluded(pkgName string, structName string, fieldName string) bool {
	var aMap map[string][]string
	switch pkgName {
	case "das_analysis":
		aMap= dasAnalysisFeStructs
	case "das_dota2":
		aMap = dasAgDota2Structs
	}

	arr, ok := 	aMap[structName]
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