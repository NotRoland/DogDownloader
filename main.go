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

type DogJSON struct {
	Message string
}

var totalDownloaded int

func downloadFile(url string, filepath string, finished chan float64){
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		finished <- 1
		return
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		finished <- 1
		return
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	finished <- 1
	totalDownloaded += 1
	fmt.Println("Downloaded "+filepath)
	return
}

func getDogs(finished chan float64, iter int) {
	//resp, err := http.Get("https://api.thedogapi.com/v1/images/search")
	resp, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil{
		fmt.Println(err)
		finished <- 2
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)

	var decoded DogJSON
	json.Unmarshal(bytes, &decoded)
	/*if len(decoded) == 0 {
		fmt.Println("API Error [IGNORING]")
		finished <- 2
		return
	}*/

	var url = decoded.Message
	fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
	fmt.Println("Downloading "+fileName+" (Thread "+strconv.Itoa(iter)+")")
	finished <- 1
	go downloadFile(url, "Dogs/"+fileName, finished)
	return
}

func loopUntilDone(finished chan float64, fullFin chan int, iters float64) {
	thrF := 0.0
	for thrF < iters*2 {
		thrF += <-finished
	}
	fullFin <- 1
}

func main() {
	fmt.Println("====\nDog Downloader\n====\nPress CTRL+C to interrupt.\nDogs will be saved in [RUNPATH]/Dogs.\n")
	fmt.Print("Enter how many dogs you want: ")
	var iters float64
	fmt.Scanln(&iters)

	finished := make(chan float64)
	fullFin := make(chan int)

	os.Mkdir("Dogs", 0755)

	go loopUntilDone(finished, fullFin, iters)

	var smCheck float64
	if iters<200 { smCheck = iters } else { smCheck = math.Ceil(iters/math.Ceil(iters/200)) }
	fmt.Println("Starting threads...")
	for j := 0.0; j<math.Ceil(iters/200); j++{
		for i := 0.0; i < smCheck; i++ {
			go getDogs(finished, int(i))
		}
		time.Sleep(1*time.Second)
	}

	fmt.Println("All threads online!")
	<- fullFin
	fmt.Println(strconv.Itoa(totalDownloaded)+" dogs downloaded! A few might be missing due to API errors.")
}
