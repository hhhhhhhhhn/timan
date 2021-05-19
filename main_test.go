package main

import (
	"testing"
)

func TestIndexReturnsM1IfNotInString(t *testing.T) {
	i := index("Hello", 'a', 0)
	if i != -1 {
		t.Log("Return value should be -1 if not found, is", i)
		t.Fail()
	}
}

func TestIndexReturnsCorrectValue(t *testing.T) {
	i := index("Hello", 'l', 0)
	if i != 2 {
		t.Log("Index of l in Hello should be 2, is", i)
		t.Fail()
	}
}

func TestIndexCanSkip(t *testing.T) {
	i := index("motorcar", 'o', 2)
	if i != 3 {
		t.Log("Index should be 3, is", i)
		t.Fail()
	}
}

func TestContainsEmptyListIsFalse(t *testing.T) {
	isContained := contains([]string{}, "")
	if isContained {
		t.Log("Should be false, is true")
		t.Fail()
	}
}

func TestContainsCanFail(t *testing.T) {
	isContained := contains([]string{"Hello", "world"}, "there")
	if isContained {
		t.Log("Should be false, is true")
		t.Fail()
	}
}

func TestContainsCanSucceed(t *testing.T) {
	isContained := contains([]string{"Hello", "world"}, "world")
	if !isContained {
		t.Log("Should be true, is false")
		t.Fail()
	}
}

var testJson = `[{"name":"Test1","imageName":"","secsUsed":0,"secsLimit":0,"secsGoal":0,"goalResets":0,"image":""},
{"name":"","imageName":"","secsUsed":10,"secsLimit":0,"secsGoal":0,"goalResets":0,"image":""}]`

func TestProgramsToJson(t *testing.T) {
	programs = []program{
		{
			Name: "Test1",
		},
		{
			SecsUsed: 10,
		},
	}
	jsonOutput := programsToJson()
	if jsonOutput != testJson {
		t.Log("Programs to JSON is incorrect:", jsonOutput)
		t.Fail()
	}
}

func TestResetPrograms(t *testing.T) {
	programs = []program{
		{
			SecsUsed:   99,
			SecsLimit:  100,
			SecsGoal:   0,
			GoalResets: 2,
		},
	}
	resetPrograms()
	if programs[0].SecsUsed != 0 {
		t.Log("SecsUsed not reset")
		t.Fail()
	}
	if programs[0].SecsLimit != 50 {
		t.Log("Limit should be 50, is", programs[0].SecsLimit)
		t.Fail()
	}
	if programs[0].GoalResets != 1 {
		t.Log("GoalResets should be 1, is", programs[0].GoalResets)
		t.Fail()
	}
	resetPrograms()
	if programs[0].SecsLimit != 0 {
		t.Log("Limit should be 0, is", programs[0].SecsLimit)
		t.Fail()
	}
	if programs[0].GoalResets != 0 {
		t.Log("GoalResets should be 0, is", programs[0].GoalResets)
		t.Fail()
	}
}

func TestGetProgramIndexByImageName(t *testing.T) {
	programs = []program{
		{
			SecsUsed: 10,
		},
		{
			ImageName: "test",
		},
	}
	i := getProgramIndexByImageName("test")
	if i != 1 {
		t.Log("Index should be 1, is", i)
		t.Fail()
	}
	i = getProgramIndexByImageName("asdfas")
	if i != -1 {
		t.Log("Value should be -1, is", i)
		t.Fail()
	}
}

func TestAddProgram(t *testing.T) {
	programs = []program{
		{},
	}
	addProgram(program{Name: "test"})
	if programs[1].Name != "test" {
		t.Log("Cannot add program")
		t.Fail()
	}
}

func TestRemoveProgram(t *testing.T) {
	programs = []program{
		{Name: "test"},
		{},
	}
	removeProgram(1)
	if programs[len(programs)-1].Name != "test" {
		t.Log("Cannot remove program")
		t.Fail()
	}
}
