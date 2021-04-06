package main

import (
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
	"time"
)

var receiver = make(chan string)

type DogJSON struct { Message string } // For decoding JSON.

func ReceiveUntilDone(count int, writer chan int) {
	totalDownloaded := 0
	for i := 1; i <= count; i++ { // Will loop until all goroutines send a channel signal.
		out := <-receiver
		setLog("Download Received >"+out, strconv.Itoa(i), strconv.Itoa(count))
		if out == "OK" {
			totalDownloaded++
		}
	}
	writer <- totalDownloaded
}

func GetAndDecode(url string) (DogJSON, error) {
	var decoded DogJSON

	client := resty.New()
	_, err := client.R().
		SetResult(&decoded). // Sets JSON output variable
		Get(url)

	if err != nil {
		return DogJSON{}, err
	}
	// TODO: poggers
	return decoded, nil
}

func processDownload() {
	decoded, err := GetAndDecode("http://www.dog.ceo/api/breeds/image/random")
	if err != nil {
		receiver <- "BAD"
		return
	}

	_ = os.Mkdir("Dogs", 0755)

	receiver <- DownloadFile(decoded.Message, "Dogs")
}

func beginBack() { // Called from Begin button.
	count, _ := strconv.Atoi(dogNumber.Text)
	waitChannel := make(chan int)

	go ReceiveUntilDone(count, waitChannel)

	for i := 1; i <= count; i++ {
		setLog("Starting Threads...", strconv.Itoa(i), strconv.Itoa(count))
		go processDownload()
		time.Sleep(time.Duration(waitTime))
	}

	total := <-waitChannel // Waits until all threads are finished.
	setLog("Downloaded "+strconv.Itoa(total)+" out of "+strconv.Itoa(count)+".", strconv.Itoa(count), strconv.Itoa(count))
}
