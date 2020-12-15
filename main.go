package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/jeffalyanak/johnny-hackerseed/icon"
	"github.com/mmcdole/gofeed"
)

var (
	window fyne.Window
	status *widget.Label
	logo   fyne.Resource

	count   int = 0
	episode string
	url     string
)

func main() {
	// Create new window and add the title and icon.
	a := app.New()
	ico, _ := icon.Asset("logo.png")
	var iconImage = &fyne.StaticResource{
		StaticName:    "Logo.png",
		StaticContent: ico}
	window = a.NewWindow("Twinnovation - Johnny Hackerseed")
	window.SetIcon(iconImage)

	// Set up the main label with some initial text
	status = widget.NewLabel("Hack the planet!")
	status.TextStyle = fyne.TextStyle{Monospace: true}
	window.SetContent(widget.NewVBox(status))

	// Begin the downloading loop and show the window.
	go beginTheMagic()
	window.ShowAndRun()
}

func beginTheMagic() {
	// Parse the Twinnovation RSS feed
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://rss.art19.com/twinnovation")

	// Inifinite loop, grab a random episode and begin to
	// download.
	for {
		total := len(feed.Items)
		i := rand.Intn(total)

		episode = feed.Items[i].Title
		url = feed.Items[i].Enclosures[0].URL

		downloadEpisode(url)
		count++
	}
}

// downloadEpisode takes a URL and downloads whatever it finds there,
// although it does not write this file to disk and essentially throws
// it away.
func downloadEpisode(u string) {
	client := http.Client{}

	// Put content on file
	resp, err := client.Get(u)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	readerpt := &PassThru{Reader: resp.Body, length: resp.ContentLength}
	_, err = ioutil.ReadAll(readerpt)
}

// PassThru struct is used to show the progress of
// the data being downloaded.
type PassThru struct {
	io.Reader
	total    int64
	length   int64
	progress float64
}

// Read function operates on the PassThre struct and does the hard
// work of calculating the percentage and outputting text to the
// terminal.
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)

	if n > 0 {
		pt.total += int64(n)
		percentage := float64(pt.total) / float64(pt.length) * float64(100)
		percentageRounded := fmt.Sprintf("%v", math.Round(percentage))

		if percentage-pt.progress > 2 {
			str := "Episode Downloaded so far:  " + fmt.Sprint(count) + "\n"
			str += "Episode Downloading now:    " + episode + "\n"
			str += "Episode Link:               " + url + "\n"
			str += "Current Download Progress:  " + fmt.Sprint(percentageRounded) + "%"

			status.SetText(str)
			pt.progress = percentage
		}
	}
	return n, err
}
