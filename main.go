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

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

const interval = 5

type program struct {
	Name       string `json:"name"`
	ImageName  string `json:"imageName"`
	SecsUsed   int    `json:"secsUsed"`
	SecsLimit  int    `json:"secsLimit"`
	SecsGoal   int    `json:"secsGoal"`
	GoalResets int    `json:"goalResets"` // time resets to reach goal
	Image      string `json:"image"`
}

var (
	mutex    sync.Mutex
	programs []program
)

var resetTime int64 = 86400 // = day
var lastReset int64 = lastMonday()
var unit int64 = 60 // seconds

// kills program by imagename
func (prog program) kill() {
	exec.Command("taskkill", "/f", "/im", prog.ImageName).Run()
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

// prints error and exits
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

// returns list of programs as JSON
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

// loads programs from JSON
func jsonToPrograms(programsJson string) ([]program, error) {
	programJsonList := strings.Split(string(programsJson)[1:], "\n")

	programList := make([]program, len(programJsonList))

	for i, programJson := range programJsonList {
		var parsedProgram program
		err := json.Unmarshal([]byte(programJson[:len(programJson)-1]), &parsedProgram)
		if err != nil {
			return nil, err
		}
		programList[i] = parsedProgram
	}
	return programList, nil
}

// reads int64 from file contents
func readInt64FromFile(filename string) (value int64, err error) {
	valueString, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	value, err = strconv.ParseInt(string(valueString), 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// saves main data
func save() {
	handle(ioutil.WriteFile("save/programs.json", []byte(programsToJson()), 0777))
	handle(ioutil.WriteFile("save/resetTime", []byte(fmt.Sprint(resetTime)), 0777))
	handle(ioutil.WriteFile("save/lastReset", []byte(fmt.Sprint(lastReset)), 0777))
	handle(ioutil.WriteFile("save/unit", []byte(fmt.Sprint(unit)), 0777))
}

// loads main data
// NOTE: no mutex lock
func load() error {
	programsJson, err := ioutil.ReadFile("save/programs.json")
	if err != nil {
		return err
	}
	programs, err = jsonToPrograms(string(programsJson))
	if err != nil {
		return err
	}
	resetTime, err = readInt64FromFile("save/resetTime")
	if err != nil {
		return err
	}
	lastReset, err = readInt64FromFile("save/lastReset")
	if err != nil {
		return err
	}
	unit, err = readInt64FromFile("save/unitString")
	if err != nil {
		return err
	}
	return nil
}

// lists all imagenames of all programs
func listRunningPrograms() (list []string, err error) {
	commandOutput, err := exec.Command("tasklist", "/fo", "csv").Output()
	if err != nil {
		return nil, nil
	}

	running := strings.Split(string(commandOutput), "\n")
	for i := 0; i < len(running); i++ {
		if running[i] == "" {
			continue
		}
		running[i] = running[i][1:index(running[i], '"', 1)]
	}
	return running, nil
}

// resets time used and changes limits/goals of each program
// NOTE: No mutex lock/unlock
func resetPrograms() {
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
}

// does the actual program monitoring and stuff
func background() {
	for true {
		running, err := listRunningPrograms()
		handle(err)

		mutex.Lock()
		for i := 0; i < len(programs); i++ {
			if contains(running, programs[i].ImageName) {
				programs[i].SecsUsed += interval
				if programs[i].SecsUsed >= programs[i].SecsLimit {
					programs[i].kill()
				}
			}
		}

		for time.Now().Unix() >= lastReset+resetTime {
			resetPrograms()
			lastReset += resetTime
		}
		mutex.Unlock()
		//logPrograms()
		time.Sleep(interval * time.Second)
	}
}

// gets the index of the first program with the given imageName, or -1
func getProgramIndexByImageName(imageName string) int {
	for i := range programs {
		if programs[i].ImageName == imageName {
			return i
		}
	}
	return -1
}

func addProgram(prog program) {
	programs = append(programs, prog)
}

func removeProgram(index int) {
	if index == -1 || index >= len(programs) {
		return
	}
	programs = append(programs[:index], programs[index+1:]...)
}

// the webserver
func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/"+r.URL.Path[1:])
	})
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			{
				"programs": %s,
				"resetTime": %v,
				"lastReset": %v,
				"unit": %v
			}
		`, programsToJson(), resetTime, lastReset, unit)
	})
	http.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {
		args := strings.Split(r.URL.RawQuery, "+")
		if len(args) == 0 {
			return
		}
		mutex.Lock()

		switch args[0] {
		case "add":
			addProgram(program{ImageName: "program.exe", Image: "./image.svg"})

		case "remove":
			if len(args) != 2 {
				break
			}
			index := getProgramIndexByImageName(args[1])
			removeProgram(index)

		case "set":
			if len(args) != 3 {
				break
			}
			value, err := strconv.ParseInt(string(args[2]), 10, 64)
			if err != nil {
				break
			}
			switch args[1] {
			case "unit":
				unit = value
			case "resetTime":
				resetTime = value
			case "lastReset":
				lastReset = value
			}

		case "change":
			if len(args) != 4 {
				break
			}
			index := getProgramIndexByImageName(args[1])
			if index == -1 {
				break
			}
			switch args[2] { // these ones are all strings
			case "name":
				programs[index].Name = args[3]
			case "imageName":
				programs[index].ImageName = args[3]
			case "image":
				programs[index].Image = args[3]

			default: // these ones are all ints
				value, err := strconv.Atoi(args[3])
				if err == nil {
					switch args[2] {
					case "secsUsed":
						programs[index].SecsUsed = value
					case "secsLimit":
						programs[index].SecsLimit = value
					case "secsGoal":
						programs[index].SecsGoal = value
					case "goalResets":
						programs[index].GoalResets = value
					}
				}
			}
		}
		mutex.Unlock()
	})
	log.Fatal(http.ListenAndServe(":5019", nil))
}

func setDefault() {
	resetTime = 86400 // = day in seconds
	lastReset = lastMonday()
	unit = 60 // = minute in seconds
	programs = []program{
		{
			ImageName:  "brave.exe",
			Name:       "Brave_browser",
			SecsLimit:  10000,
			SecsUsed:   990,
			SecsGoal:   9000,
			GoalResets: 5,
			Image:      "./image.svg",
		}, {
			ImageName:  "chrome.exe",
			Name:       "Chrome_browser",
			SecsLimit:  1000,
			SecsUsed:   990,
			SecsGoal:   900,
			GoalResets: 10,
			Image:      "./image.svg",
		},
	}
}

func setupTray() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("timan")
	systray.SetTooltip("timan")
	mUrl := systray.AddMenuItem("Open", "Open in browser")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit timan")

	go func() {
		for {
			select {
			case <-mQuitOrig.ClickedCh:
				systray.Quit()
			case <-mUrl.ClickedCh:
				exec.Command("cmd", "/c", "start", "http://localhost:5019").Start()
			}
		}
	}()
}

// entry
func main() {
	mutex.Lock()
	err := load()
	if err != nil {
		setDefault()
	}
	mutex.Unlock()

	go background()
	go server()
	systray.Run(setupTray, save)
}
