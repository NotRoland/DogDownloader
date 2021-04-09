package main

import (
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
	"sync"
	"time"
)

type DogJSON struct{ Message string } // For decoding JSON.

func GetAndDecode(url string) (DogJSON, error) {
	var decoded DogJSON

	client := resty.New()
	_, err := client.R().
		SetResult(&decoded). // Sets JSON output variable
		Get(url)

	if err != nil {
		return DogJSON{}, err
	}
	return decoded, nil
}

func processDownload(fileStatus sync.WaitGroup) {
	fileStatus.Add(1); defer fileStatus.Done()

	decoded, err := GetAndDecode("http://www.dog.ceo/api/breeds/image/random")
	if err != nil {
		return
	}

	_ = os.Mkdir("Dogs", 0755)

	DownloadFile(decoded.Message, "Dogs")
}

func beginBack() { // Called from Begin button.
	count, _ := strconv.Atoi(dogNumber.Text)
	var fileStatus sync.WaitGroup

	for i := 1; i <= count; i++ {
		setLog("Starting Threads...", strconv.Itoa(i), strconv.Itoa(count))

		go processDownload(fileStatus)
		time.Sleep(time.Duration(waitTime))
	}

	fileStatus.Wait() // Waits until all threads are done.
	setLog("Downloaded dogs successfully!", strconv.Itoa(count), strconv.Itoa(count))
}
