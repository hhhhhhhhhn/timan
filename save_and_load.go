package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

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
