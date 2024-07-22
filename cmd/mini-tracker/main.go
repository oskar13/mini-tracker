package main

import (
	"os"

	torrenttools "github.com/oskar13/mini-tracker/pkg/torrent-tools"
)

func main() {

	torrenttools.DecodeUploadedTorrent(os.Args[1])
	//tracker.StartTracker()

	//web.StartWebsite()

}
