package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

const interval = 5

type program struct {
	Name       string `json:"name"`
	ImageName  string `json:"imagename"`
	SecsUsed   int    `json:"secsUsed"`
	SecsLimit  int    `json:"secsLimit"`
	SecsGoal   int    `json:"secsGoal"`
	GoalResets int    `json:"goalRestes"` // time resets to reach goal
	Image      string `json:"image"`
}

var (
	mutex sync.Mutex
)

var programs []program = make([]program, 0, 10)
var resetTime int64 = 86400 // = day
var lastReset int64 = lastMonday()

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

// handles errors
func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

func programsToJson() string {
	output := "["
	for i, program := range programs {
		json, err := json.Marshal(program)
		handle(err)
		output += string(json)
		if i != len(programs)-1 {
			output += ",\n"
		}
	}
	return output + "]"
}

func save() {
	handle(ioutil.WriteFile("save/programs.json", []byte(programsToJson()), 0777))
	handle(ioutil.WriteFile("save/resetTime", []byte(fmt.Sprint(resetTime)), 0777))
	handle(ioutil.WriteFile("save/lastReset", []byte(fmt.Sprint(lastReset)), 0777))
}

func load() error {
	programsJson, err := ioutil.ReadFile("save/programs.json")
	if err != nil {
		return err
	}
	programJsonList := strings.Split(string(programsJson)[1:], "\n")

	mutex.Lock()
	for i, programJson := range programJsonList {
		var parsedProgram program
		err := json.Unmarshal([]byte(programJson[:len(programJson)-1]), &parsedProgram)
		if err != nil {
			return err
		}
		programs[i] = parsedProgram
	}

	resetTimeString, err := ioutil.ReadFile("save/resetTime")
	if err != nil {
		return err
	}
	resetTime, err = strconv.ParseInt(string(resetTimeString), 10, 64)
	if err != nil {
		return err
	}

	lastResetString, err := ioutil.ReadFile("save/lastReset")
	if err != nil {
		return err
	}
	lastReset, err = strconv.ParseInt(string(lastResetString), 10, 64)
	if err != nil {
		return err
	}

	mutex.Unlock()
	return nil
}

// does the actual program monitoring and stuff
func background() {
	for true {
		output, err := exec.Command("tasklist", "/fo", "csv").Output()
		handle(err)

		running := strings.Split(string(output), "\n")
		for i := 0; i < len(running); i++ {
			if running[i] == "" {
				continue
			}
			running[i] = running[i][1:index(running[i], '"', 1)]
		}

		mutex.Lock()
		for i := 0; i < len(programs); i++ {
			if contains(running, programs[i].ImageName) {
				programs[i].SecsUsed += interval
				if programs[i].SecsUsed >= programs[i].SecsLimit {
					exec.Command("taskkill", "/f", "/im", programs[i].ImageName).Run()
				}
			}
		}

		for time.Now().Unix() >= lastReset+resetTime {
			for i := 0; i < len(programs); i++ {
				programs[i].SecsUsed = 0
				if programs[i].GoalResets != 0 {
					programs[i].SecsLimit +=
						(programs[i].SecsGoal - programs[i].SecsLimit) / programs[i].GoalResets
					programs[i].GoalResets--
				} else {
					programs[i].SecsLimit = programs[i].SecsGoal
				}
			}
			lastReset += resetTime
		}
		mutex.Unlock()
		logPrograms()
		time.Sleep(interval * time.Second)
	}
}

func server() {
}

// entry
func main() {
	err := load()
	if err != nil {
		programs[0] = program{
			ImageName:  "brave.exe",
			Name:       "Brave browser",
			SecsLimit:  10000,
			SecsUsed:   9990,
			SecsGoal:   9000,
			GoalResets: 5,
		}
		programs[1] = program{
			ImageName:  "chrome.exe",
			Name:       "Chrome browser",
			SecsLimit:  1000,
			SecsUsed:   990,
			SecsGoal:   900,
			GoalResets: 10,
		}
	}
	go background()
}
