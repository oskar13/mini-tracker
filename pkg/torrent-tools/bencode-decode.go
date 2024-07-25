package torrenttools

import (
	"fmt"
	"io"

	gotorrentparser "github.com/oskar13/go-torrent-parser"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func DecodeUploadedTorrent(file io.Reader) (webdata.TorrentWeb, error) {

	decoded, err := gotorrentparser.Parse(file)
	var torrent webdata.TorrentWeb

	if err != nil {
		fmt.Println(err)
		return webdata.TorrentWeb{}, err
	}

	fmt.Println(decoded.Files)

	torrent.Name = string(decoded.Files[0].Path[0])

	torrent.InfoField = decoded.Metadata.Info
	torrent.InfoHash = decoded.InfoHash
	torrent.Announce = decoded.Announce
	return torrent, nil
}

func TorrentFromWebTorrent(torrent webdata.TorrentWeb) ([]byte, error) {
	newTorrentFile, err := gotorrentparser.NewTorrent(torrent.Announce, torrent.InfoField)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return newTorrentFile, nil
}
