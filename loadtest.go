package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

//time.Sleep(100 * time.Millisecond)
/*GET requests are for retrieving information,
POST requests are for creating data,
PUT requests are for updating existing records. for HTTP*/

func client() {
	t0 := time.Now()
	//count_sent++
	con, err := net.Dial("tcp", "0.0.0.0:9090")

	if err != nil {
		log.Fatalln(err)
		count_fail++
	}
	defer con.Close()

	// clientReader := bufio.NewReader(os.Stdin)
	// serverReader := bufio.NewReader(con)
	var request = make([]byte, 100)

	for {
		_, err = con.Read(request)

		if err != nil {
			log.Println("failed to read request contents")
			count_fail++
			return
		}
		fmt.Printf(" Latency Time:   %v ", time.Since(t0))
		log.Println(&con, string(request))
		count_res++
		return
		// request = make([]byte, 100)
	}
}

//var count_sent = 0

var count_res = 0
var count_fail = 0

//https://gist.github.com/AntoineAugusti/80e99edfe205baf7a094
func main() {
	//var start int = time.Now()
	Maxroutine := flag.Int("maxNbConcurrentGoroutines", 15000, "the number of goroutines that are allowed to run concurrently")
	nbclients := flag.Int("nbJobs", 15000, "the number of jobs that we need to do")
	flag.Parse()
	//concurrentGoroutines
	ch := make(chan struct{}, *Maxroutine)

	for i := 0; i < *Maxroutine; i++ {
		ch <- struct{}{}
	}
	done := make(chan bool)
	waitForAllclients := make(chan bool)

	// Collect all the jobs, and since the job is finished, we can release another spot for a goroutine.
	go func() {
		for i := 0; i < *nbclients; i++ {
			<-done
			// Say that another goroutine can now start.
			ch <- struct{}{}
		}
		waitForAllclients <- true
	}()

	// Try to start nbclients jobs
	start := time.Now()
	for i := 1; i <= *nbclients; i++ {
		// time := time.Now()
		fmt.Printf("ID: %v: waiting to launch!\n", i)
		// Try to receive from the channel when one finish and start new routine.
		// If not, it will block the execution until receive something.
		<-ch
		//fmt.Printf("ID: %v: it's my turn!\n", i)
		go func(id int) {
			client()
			//fmt.Printf("ID: %v: all done!\n", id)
			done <- true
		}(i)
		// end = time.Now()
	}
	<-waitForAllclients // Wait for all clients to finish/ terminate

	//time sine (total from all requests)
	fmt.Printf(" Total Timeeeeeeeeeeeeeeeeee:   %v \n", time.Since(start))
	fmt.Printf("Number Sent: %d\n", *nbclients)
	fmt.Printf("Number Response: %d\n", count_res)
	fmt.Printf("Number fail: %d\n", count_fail)
	fmt.Printf(" done \n")

}
