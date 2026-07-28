package domain

import (
	"testing"
)

func TestIsValidCategory(t *testing.T) {
	validCategories := []string{
		CategoryCasa,
		CategoryCarro,
		CategoryDog,
		CategoryMercado,
		CategoryViagem,
		CategoryLazer,
		CategorySaude,
		CategoryOutros,
	}

	for _, cat := range validCategories {
		if !IsValidCategory(cat) {
			t.Errorf("Expected category %s to be valid", cat)
		}
	}

	invalidCategories := []string{
		"",
		"Invalid",
		"casa",
		"MERCADO",
		" Random ",
	}

	for _, cat := range invalidCategories {
		if IsValidCategory(cat) {
			t.Errorf("Expected category %s to be invalid", cat)
		}
	}
}
