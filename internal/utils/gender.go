package utils

import "strings"


type Gender uint 

const (
	Unknown Gender = iota
	Male 
	Female
	Intersex
)



func (gender Gender) String()	string {
	terms := []string{"unknown","Male", "Female", "Intersex"}
	if gender < Male || gender > Intersex {
		return terms[Unknown]
}

return terms[gender]

}


func StringPtr(s string) *string{
	return &s
}


func NormalizeRelationship(relationship string) string {
	switch strings.ToLower(relationship){
    case "father", "dad", "daddy":
		return "Father"
	case "mother", "mom", "mommy":
		return "Mother"
	case "sibling", "brother", "sister":
		return "Sibling"
		case "guardian":
		return "Guardian"
		default:
			return relationship
	}
}