package main

import (
	"fmt"
	"github.com/cavaliercoder/grab"
	"os"
	"os/exec"
	"runtime"
)

func setLog(log string, current string, max string){
	/* Sets log labels
	current - current progress
	max - max progress
	 */
	logLabel.Text = "[LOG] "+log; progressLabel.Text = "("+current+"/"+max+")"
	logLabel.Refresh(); progressLabel.Refresh()
}

func removeImages(){ // Called from Remove Images button.
	_ = os.RemoveAll("Dogs")
	setLog("Images removed.", "1", "1")
}

func openGithub(){ // Called from Github button.
	OpenBrowser("https://github.com/NotRoland/DogDownloader")
	setLog("Opened Github.", "1", "1")
}

func OpenBrowser(url string) { // Opens a browser depending on platform.
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	default:
		fmt.Errorf("unsupported platform")
	}

}

func DownloadFile(url string, filepath string) string {
	_, err := grab.Get(filepath, url)
	if err != nil {
		fmt.Println(err)
		return "BAD"
	}
	return "OK"
}
