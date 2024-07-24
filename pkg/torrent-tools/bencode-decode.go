package torrenttools

import (
	"fmt"
	"io"

	gotorrentparser "github.com/oskar13/go-torrent-parser"
)

func DecodeUploadedTorrent(file io.Reader) (*gotorrentparser.Torrent, error) {

	torrent, err := gotorrentparser.Parse(file)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return torrent, nil
}
