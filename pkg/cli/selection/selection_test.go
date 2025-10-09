package selection

import (
	"testing"
)

func TestFilterArray1(t *testing.T) {
	filter := "abc"
	items := []SelectItem{
		{
			Value: "ab",
			Label: "ab",
		},
		{
			Value: "abc",
			Label: "abc",
		},
		{
			Value: "abcd",
			Label: "abcd",
		},
	}

	result := filterArray(items, filter)

	if len(result) != 3 {
		t.Errorf("filterArray() = %d; wanted length %v", len(result), 3)
	}
}

func TestFilterArray2(t *testing.T) {
	filter := "haus"
	items := []SelectItem{
		{
			Value: "baum",
			Label: "baum",
		},
		{
			Value: "raus",
			Label: "raus",
		},
		{
			Value: "maus",
			Label: "maus",
		},
	}

	result := filterArray(items, filter)

	if len(result) != 3 {
		t.Errorf("filterArray() = %d; wanted length %v", len(result), 3)
	}
	if result[0].Value != "raus" {
		t.Errorf("filterArray()[0] = %s; wanted %s", result[0].Value, "raus")
	}
	if result[1].Value != "maus" {
		t.Errorf("filterArray()[0] = %s; wanted %s", result[0].Value, "raus")
	}
}

func TestFilterArray3(t *testing.T) {
	filter := "abc"
	items := []SelectItem{
		{
			Value: "acd",
			Label: "acd",
		},
	}

	result := filterArray(items, filter)

	if len(result) != 0 {
		t.Errorf("filterArray() = %d; wanted length %v", len(result), 1)
	}
}

func TestFilterArraySorting(t *testing.T) {
	filter := "comd"
	items := []SelectItem{
		{
			Value: "test.com",
			Label: "test.com",
		},
		{
			Value: "comdirect.de",
			Label: "comdirect.de",
		},
	}

	result := filterArray(items, filter)

	if len(result) != 2 {
		t.Errorf("filterArray() = %d; wanted length %v", len(result), 1)
	}
	if result[0].Value != "comdirect.de" {
		t.Errorf("filterArray()[0] = %s; wanted %s", result[0].Value, "comdirect")
	}
}
