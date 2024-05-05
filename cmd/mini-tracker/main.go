package main

import (
	"os"

	torrenttools "github.com/oskar13/mini-tracker/pkg/torrent-tools"
)

func main() {

	torrenttools.DecodeTorrent(os.Args[1])

}
