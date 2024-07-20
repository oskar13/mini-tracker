package torrenttools

import (
	"fmt"
	"net/url"
	"os"

	gotorrentparser "github.com/j-muller/go-torrent-parser"
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
	fmt.Println(url.QueryEscape(torrent.InfoHash))
	fmt.Println(torrent.Files[0].Path)

}
