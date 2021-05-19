package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/"+r.URL.Path[1:])
	})
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		logPrograms()
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
