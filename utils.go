package main

import (
	"fmt"
	"log"
	"time"
)

// prints error and exits
func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// gets timestamp of monday morning in seconds
func lastMonday() int64 {
	output := time.Now()
	if output.Weekday() == 0 {
		output = output.AddDate(0, 0, -6)
	} else {
		output = output.AddDate(0, 0, -int(output.Weekday()-1))
	}
	output = output.Add(-time.Duration(output.Hour()) * time.Hour)
	output = output.Add(-time.Duration(output.Minute()) * time.Minute)
	output = output.Add(-time.Duration(output.Second()) * time.Second)

	return output.Unix()
}

// index of character chr in string str, start searching from index start
func index(str string, chr byte, start int) int {
	for i := start; i < len(str); i++ {
		if str[i] == chr {
			return i
		}
	}
	return -1
}

// checks if array includes the string
func contains(array []string, str string) bool {
	for _, element := range array {
		if element == str {
			return true
		}
	}
	return false
}

// logs all programs
func logPrograms() {
	println()
	for _, program := range programs {
		fmt.Println("Name:", program.Name)
		fmt.Println("Image Name:", program.ImageName)
		fmt.Println("Used:", program.SecsUsed, "s")
		fmt.Println("Limit:", program.SecsLimit, "s")
		fmt.Println("Goal:", program.SecsGoal, "s")
		fmt.Println("Goal Resets:", program.GoalResets)
		println()
	}
}
