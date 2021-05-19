package main

import (
	"os/exec"
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
	for {
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
