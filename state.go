package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
)

type State struct {
	ArchiveFile    string
	TargetFile     string
	Charset        []rune
	StartLength    int
	EndLength      int
	CurrentLength  int
	IterationStart uint32
	IterationEnd   uint32
}

func NextState(previous_state State, iterations int) (State, error) {
	state := previous_state
	var stateError error = nil
	charset_length := len(state.Charset)
	total_combinations := uint32(math.Pow(float64(charset_length), float64(state.CurrentLength)))
	state.IterationStart = state.IterationEnd + 1
	state.IterationEnd = state.IterationStart + uint32(iterations) - 1
	if state.IterationStart >= total_combinations {
		if state.CurrentLength < state.EndLength {
			state.CurrentLength++
			total_combinations = uint32(math.Pow(float64(charset_length), float64(state.CurrentLength)))
			state.IterationStart = 0
			state.IterationEnd = uint32(iterations) - 1
		} else {
			stateError = fmt.Errorf("max length reached: %d", state.CurrentLength)
		}
	}
	if state.IterationEnd >= total_combinations {
		state.IterationEnd = total_combinations - 1
	}
	if stateError == nil {
		SaveState(state)
	}
	return state, stateError
}

func newState(archive_file string, charset []rune, start_length, end_length, iterations int) State {
	total_combinations := uint32(math.Pow(float64(len(charset)), float64(start_length)))
	var iteration_end uint32 = uint32(iterations) - 1
	if iteration_end >= total_combinations {
		iteration_end = total_combinations - 1
	}
	state := State{
		ArchiveFile:    archive_file,
		TargetFile:     GetSmallestFile(archive_file),
		Charset:        charset,
		StartLength:    start_length,
		EndLength:      end_length,
		CurrentLength:  start_length,
		IterationStart: 0,
		IterationEnd:   iteration_end,
	}
	return state
}

func LoadState(archive_file string, charset []rune, start_length, end_length, iterations int) State {
	state_file := archive_file + ".json"
	var state State
	content, err := os.ReadFile(state_file)
	new_state := newState(archive_file, charset, start_length, end_length, iterations)
	if err != nil {
		state = new_state
		SaveState(state)
	} else {
		json.Unmarshal(content, &state)
		if !reflect.DeepEqual(new_state.Charset, state.Charset) || new_state.StartLength > state.CurrentLength || new_state.EndLength != state.EndLength {
			state = new_state
			SaveState(state)
		}
	}
	return state
}

func SaveState(state State) {
	state_file := state.ArchiveFile + ".json"
	stateJson, _ := json.MarshalIndent(state, "", " ")
	log.Println("Saving state:", string(stateJson))
	os.WriteFile(state_file, stateJson, 0644)
}
