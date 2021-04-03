package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type DogJSON struct { Message string }
var totalDownloaded int

func downloadFile(url string, filepath string) float64{
	resp, err := http.Get(url); defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	out, err := os.Create(filepath); defer out.Close()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	_, err = io.Copy(out, resp.Body)
	totalDownloaded += 1
	fmt.Println("Downloaded "+filepath)
	return 1
}

func processDogs(finished chan float64, iter int){ finished <- getDogs(finished, iter) }
func processDownload(url string, filepath string, finished chan float64){ finished <- downloadFile(url, filepath) }

func getDogs(finished chan float64, iter int) float64 {
	resp, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil{
		return 2
	}
	bytes, _ := ioutil.ReadAll(resp.Body)

	var decoded DogJSON
	json.Unmarshal(bytes, &decoded)

	var url = decoded.Message
	fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	fmt.Println("Downloading "+fileName+" (Thread "+strconv.Itoa(iter)+")")
	go processDownload(url, "Dogs/"+fileName, finished)
	return 1
}

func loopUntilDone(finished chan float64, fullFin chan int, iters float64) {
	thrF := 0.0
	for thrF < iters*2 {
		thrF += <-finished
	}
	fullFin <- 1
}

func main() {
	fmt.Print("====\nDog Downloader\n====\nPress CTRL+C to interrupt.\nDogs will be saved in [RUNPATH]/Dogs.\n\nEnter how many dogs you want: ")
	var iters float64
	fmt.Scanln(&iters)

	finished, fullFin := make(chan float64), make(chan int)

	os.Mkdir("Dogs", 0755)
	go loopUntilDone(finished, fullFin, iters)

	var smCheck float64
	if iters<200 { smCheck = iters } else { smCheck = math.Ceil(iters/math.Ceil(iters/200)) }

	fmt.Println("Starting threads...")
	for j := 0.0; j<math.Ceil(iters/200); j++{
		for i := 0.0; i < smCheck; i++ {
			go processDogs(finished, int(i))
		}
		time.Sleep(time.Second)
	}

	fmt.Println("All threads online!")
	<- fullFin
	fmt.Println(strconv.Itoa(totalDownloaded)+" dogs downloaded! A few might be missing due to API errors.")
}
