package torrenttools

import (
	"fmt"
	"os"

	gotorrentparser "github.com/oskar13/go-torrent-parser"
)

func DecodeUploadedTorrent(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	torrent, err := gotorrentparser.Parse(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(torrent.Announce)
	fmt.Println(torrent.InfoHash)
	fmt.Println(torrent.Files[0].Path)

}
