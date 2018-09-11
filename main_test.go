package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestRandomValue(t *testing.T) {
	sum := 0
	//TODO: This should be a map so it can be broken out into a helper function
	hasBeen5 := false
	hasBeen6 := false
	hasBeen7 := false
	hasBeen8 := false
	hasBeen9 := false
	hasBeen10 := false
	hasBeen11 := false
	hasBeen12 := false
	hasBeen13 := false
	hasBeen14 := false
	hasBeen15 := false
	allNumHit := false

	for i := 0; i < 10000; i++ {
		num := randomValue(5, 15)
		switch num {
		case 5:
			if !hasBeen5 {
				fmt.Println("5")
				hasBeen5 = true
			}
		case 6:
			if !hasBeen6 {
				fmt.Println("6")
				hasBeen6 = true
			}
		case 7:
			if !hasBeen7 {
				fmt.Println("7")
				hasBeen7 = true
			}
		case 8:
			if !hasBeen8 {
				fmt.Println("8")
				hasBeen8 = true
			}
		case 9:
			if !hasBeen9 {
				fmt.Println("9")
				hasBeen9 = true
			}
		case 10:
			if !hasBeen10 {
				fmt.Println("10")
				hasBeen10 = true
			}
		case 11:
			if !hasBeen11 {
				fmt.Println("11")
				hasBeen11 = true
			}
		case 12:
			if !hasBeen12 {
				fmt.Println("12")
				hasBeen12 = true
			}
		case 13:
			if !hasBeen13 {
				fmt.Println("13")
				hasBeen13 = true
			}
		case 14:
			if !hasBeen14 {
				fmt.Println("14")
				hasBeen14 = true
			}
		case 15:
			if !hasBeen15 {
				fmt.Println("15")
				hasBeen15 = true
			}
		default:
			break
		}

		if hasBeen5 &&
			hasBeen6 &&
			hasBeen7 &&
			hasBeen8 &&
			hasBeen9 &&
			hasBeen10 &&
			hasBeen11 &&
			hasBeen12 &&
			hasBeen13 &&
			hasBeen14 &&
			hasBeen15 {
			allNumHit = true
			break
		}
		sum += i
	}

	if !allNumHit {
		t.Errorf("not all numbers hit in test")
	}
}
