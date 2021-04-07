# Dog Downloader
This is the first project I've made by myself in Go. I'm trying to learn how to use Git as well, so I posted it here.
### What does it do?
It downloads dog pictures from the dog.ceo API asynchronously using goroutines. There's also a simple GUI interface that you can use.
### How do I use it?
I cannot provide any binaries for this, but you can compile it yourself with `go build .` and then `fyne package`.
#### Dependencies
(Provided in `go.mod`)
- [grab](github.com/cavaliercoder/grab) for downloading files. 
- [resty](github.com/go-resty/resty/v2) for HTTP requests and decoding JSON. 
- [fyne](fyne.io) for the GUI.
