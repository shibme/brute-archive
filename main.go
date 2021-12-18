package main

import (
	"flag"
	"log"
	"os"
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
		state, err = NextState(state, batch_size)
	}
}

func main() {
	var archive_file string
	var charset string
	var min_length int
	var max_length int
	var batch_size int
	flag.StringVar(&archive_file, "file", "", "the archive file to target")
	flag.StringVar(&charset, "charset", "abcdefghijklmnopqrstuvwxyz", "the character set to use for bruteforce")
	flag.IntVar(&min_length, "min", 1, "the minimum number of characters to brute-force the password with")
	flag.IntVar(&max_length, "max", 6, "the minimum number of characters to brute-force the password with")
	flag.IntVar(&batch_size, "threads", 1000, "the number of executor threads/routines")
	flag.Parse()
	if archive_file == "" {
		flag.Usage()
	} else if _, err := os.Stat(archive_file); os.IsNotExist(err) {
		flag.Usage()
	} else {
		crackRar(archive_file, charset, min_length, max_length, batch_size)
	}
}
