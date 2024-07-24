package main

import (
	"github.com/oskar13/mini-tracker/pkg/tracker"
	"github.com/oskar13/mini-tracker/pkg/web"
)

func main() {

	//torrenttools.DecodeUploadedTorrent(os.Args[1])

	go func() {
		tracker.StartTracker()
	}()
	web.StartWebsite()

}
