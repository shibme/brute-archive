package main

import (
	"log"
	"sync"
)

func passwordRunner(wg *sync.WaitGroup, state State, password string, matched_password *string) {
	defer wg.Done()
	if AttemptPassword(state.ArchiveFile, state.TargetFile, password) {
		*matched_password = password
	}
}

func groupExecutor(state State) (string, error) {
	matched_password := ""
	passwords, err := GetPasswordsForState(state)
	wg := new(sync.WaitGroup)
	if err != nil {
		return matched_password, err
	}
	for i := 0; i < len(passwords); i++ {
		wg.Add(1)
		go passwordRunner(wg, state, passwords[i], &matched_password)
	}
	wg.Wait()
	return matched_password, nil
}

func passwordFoundEvent(password string) {
	log.Println("Password Found:", password)
}

func crackRar(archive_file, charset string, min_length, max_length, batch_size int) {
	charset_arr := []rune(charset)
	state := LoadState(archive_file, charset_arr, min_length, max_length, batch_size)
	var err error = nil
	for err == nil {
		matched_password, executorError := groupExecutor(state)
		if executorError != nil {
			panic(executorError)
		}
		if matched_password != "" {
			passwordFoundEvent(matched_password)
			break
		}
		state, err = NextState(state, 100)
	}
}

func main() {
	archive_file := "test.rar"
	charset := "txyur"
	charset_arr := []rune(charset)
	state := LoadState(archive_file, charset_arr, 5, 5, 100)
	var err error = nil
	for err == nil {
		matched_password, executorError := groupExecutor(state)
		if executorError != nil {
			panic(executorError)
		}
		if matched_password != "" {
			passwordFoundEvent(matched_password)
			break
		}
		state, err = NextState(state, 100)
	}
}
