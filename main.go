package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	// Blocking
	//blockingDemo()

	// Goroutine
	goroutineDemo()

	// With channel
	channel := make(chan int)
	go sendToChannel(channel)
	go receiveFromChannel(channel)

	// Channel is Blocking, error, do not run !!
	//blockingChannelDemo()

	// Channel with buffer
	bufferedChannelDemo()

	// With channel then receive in range
	channel2 := make(chan int)
	go sendToChannel(channel2)
	go receiveWithRange(channel2)

	// With channel then close
	go closeChannelDemo()

	// Select
	go selectDemo()

	// sync.WaitGroups
	go waitGroupsDemo()

	otherTasks()
	time.Sleep(22 * time.Second)
}

func selectDemo() {
	channel3 := make(chan int)
	go calcTripleToChannel(3, channel3)
	channel4 := make(chan string)
	go aVeryLongTimeProcess(channel4)

	select {
	case int1 := <-channel3:
		fmt.Println("Received from channel3, value: " + strconv.Itoa(int1))
	case text := <-channel4:
		fmt.Println(text)
	}
}

func blockingChannelDemo() {
	channel := make(chan int)
	calcTripleToChannel(3, channel)
	calcTripleToChannel(30, channel)

	num1, num2 := <-channel, <-channel
	fmt.Println(strconv.Itoa(num1) + " " + strconv.Itoa(num2))
}

func bufferedChannelDemo() {
	channel := make(chan int, 2)
	calcTripleToChannel(3, channel)  // Take 1 buffer, 1 buffer remaining
	calcTripleToChannel(30, channel) // Take another 1 buffer, no buffer remaining

	// Error, blocking again, do not uncomment this
	//calcTripleToChannel(50, channel)

	num1, num2 := <-channel, <-channel // Receiving from Channel, 2 buffers remaining
	fmt.Println(strconv.Itoa(num1) + " " + strconv.Itoa(num2))
}

func aVeryLongTimeProcess(channel chan string) {
	time.Sleep(8 * time.Second)
	channel <- "After 8 seconds"
}

func receiveWithRange(channel chan int) {
	for tripledNum := range channel {
		fmt.Println(strconv.Itoa(tripledNum))
	}
}

func closeChannelDemo() {
	channel := make(chan int)

	go calcTripleToChannel(3, channel)
	// int3 is 3, ok should be true
	int3, ok := <-channel
	fmt.Println(int3)

	go calcTripleToChannel(5, channel)
	// int5 is 5, ok should be true
	int5, ok := <-channel
	fmt.Println(int5)

	close(channel)

	// Error, channel already closed
	//go calcTripleToChannel(10, channel)

	// Channel already closed, ok is false
	_, ok = <-channel
	fmt.Println(strconv.FormatBool(ok))
}

func sendToChannel(channel chan int) {
	for i := 1; i <= 10; i++ {
		fmt.Println("Emitting " + strconv.Itoa(i))
		go calcTripleToChannel(i, channel)
	}
}

func calcTripleToChannel(num int, channel chan int) {
	tripledNum := triple(num)
	channel <- tripledNum
}

func receiveFromChannel(channel chan int) {
	for i := 1; i <= 10; i++ {
		tripledNum, ok := <-channel
		if ok {
			fmt.Println("Receiving from channel: " + strconv.Itoa(tripledNum))
		} else {
			fmt.Println("Channel Closed")
		}
	}
}

func triple(num int) int {
	// Simulate long time processing, example: sending email, calculating some numbers,
	// inserting into database, etc anything that can be processed in parallel.
	time.Sleep(2 * time.Second)

	return num * 3
}
func calcTripleAndPrint(num int) {
	tripledNum := triple(num)
	fmt.Println(strconv.Itoa(num) + " tripled: " + strconv.Itoa(tripledNum))
}

func goroutineDemo() {
	for i := 1; i <= 20; i++ {
		go calcTripleAndPrint(i)
	}
}

func blockingDemo() {
	for i := 1; i <= 20; i++ {
		calcTripleAndPrint(i)
	}
}

func calcTripleAndPrintWithWaitGroup(num int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	calcTripleAndPrint(num)
}

func waitGroupsDemo() {
	var waitGroups sync.WaitGroup

	waitGroups.Add(1) // counter is 1
	waitGroups.Add(1) // counter is 2
	waitGroups.Add(1) // counter is 3
	waitGroups.Done() // counter is 2
	waitGroups.Done() // counter is 1
	time.Sleep(2 * time.Second)
	waitGroups.Done() // counter is back to 0

	waitGroups.Wait() // blocking for 2 seconds, will wait until counter is 0
	fmt.Println("Marco")

	for i := 1; i <= 10; i++ {
		waitGroups.Add(1)
		go calcTripleAndPrintWithWaitGroup(i, &waitGroups)
	}

	waitGroups.Wait()
	fmt.Println("Polo")
}

func otherTasks() {
	fmt.Println("Another important tasks")
}
