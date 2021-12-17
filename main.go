package main

import (
	"sync"
)

func passwordRunner(wg *sync.WaitGroup, state State, password string, matched_password *string) {
	defer wg.Done()
	if password == "shibly" {
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
		go passwordRunner(wg, state, passwords[i], &matched_password)
	}
	wg.Add(len(passwords))
	wg.Wait()
	return matched_password, nil
}

/* func crack_pass(archive_file string, charset string, start_length, end_length int) string {
	var save_state_per_iterations uint32 = 4
	charset_arr := []rune(charset)
	state := LoadState(NewState(archive_file, charset_arr, start_length, end_length))

	charset_length := uint32(len(charset_arr))
	routines_count := 100
	target_file := GetSmallestFile(archive_file)
	queued_passwords := make(chan string)
	processed_passwords := make(chan string)
	state_queue := make(chan State)
	//go SaveStateQueue(archive_file, state_queue)
	matched_password := ""
	for i := 0; i < routines_count; i++ {
		go runForPassword(archive_file, target_file, queued_passwords, processed_passwords, &matched_password, i)
	}
	for i := start_length; i <= end_length; i++ {
		current_length := i
		log.Println("Starting iteration for length: ", i)
		total_combinations := uint32(math.Pow(float64(charset_length), float64(current_length)))
		var j uint32
		log.Println(total_combinations)
		for j = 0; j < total_combinations && matched_password == ""; j++ {
			password, err := GetBruteString(charset_arr, current_length, j)
			queued_passwords <- password
			if err != nil {
				break
			}
			if (j > 0 && j%save_state_per_iterations == 0) || j == (total_combinations-1) {
				for {
					processed_password := <-processed_passwords
					if processed_password == password {
						state_queue <- State{
							File:             archive_file,
							EndLength:        end_length,
							CurrentLength:    current_length,
							CurrentIteration: j,
						}
						break
					}
				}
			}
			if matched_password != "" {
				return matched_password
			}
			//time.Sleep(1 * time.Second)
		}
	}
	return ""
} */

/* func main() {
	archive_file := "test.rar"
	charset := "shilyb"
	crack_pass(archive_file, charset, 6, 6)
} */

func main() {
	archive_file := "test.rar"
	charset := "hilybs"
	charset_arr := []rune(charset)
	state := LoadState(archive_file, charset_arr, 2, 4, 100)
	var err error = nil
	for i := 0; i < 1; i++ {
		state, err = NextState(state, 100)
		if err != nil {
			break
		}
	}
}
