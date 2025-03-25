package main

import (
	"github.com/oskar13/mini-tracker/pkg/publictracker"
	"github.com/oskar13/mini-tracker/pkg/web"
)

func main() {

	//torrenttools.DecodeUploadedTorrent(os.Args[1])

	go func() {
		publictracker.StartTracker()
	}()
	web.StartWebsite()

}
