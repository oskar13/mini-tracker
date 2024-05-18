package main

import (
	"os"

	torrenttools "github.com/oskar13/mini-tracker/pkg/torrent-tools"
	web "github.com/oskar13/mini-tracker/pkg/web"
)

func main() {

	torrenttools.DecodeUploadedTorrent(os.Args[1])

	web.StartWebsite()

}
