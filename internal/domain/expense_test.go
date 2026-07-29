package domain

import (
	"testing"
)

func TestIsValidCategory(t *testing.T) {
	validCategories := []string{
		CategoryFood,
		CategoryHousing,
		CategoryTransport,
		CategoryEntertainment,
		CategoryUtilities,
		CategoryPet,
		CategoryTravel,
		CategoryHealth,
		CategoryOther,
		CategoryCasa,
		CategoryCarro,
	}

	for _, cat := range validCategories {
		if !IsValidCategory(cat) {
			t.Errorf("Expected category %s to be valid", cat)
		}
	}

	invalidCategories := []string{
		"",
		"Invalid",
		"UnknownCategory",
		" Random ",
	}

	for _, cat := range invalidCategories {
		if IsValidCategory(cat) {
			t.Errorf("Expected category %s to be invalid", cat)
		}
	}
}
