package data

import "github.com/oskar13/mini-tracker/pkg/web/webdata"

// Data object to hold torrent data when processing templates
type Torrent struct {
	Announce    string
	Name        string
	Description string //comment added on web page
	Comment     string //file metadata comment
	Type        string //category it was posted in
	PieceLength int64
	Pieces      []byte
	Private     bool
	GroupID     *int
	GroupName   string
	User        webdata.User //Uploader
	Date        int64        //Uploaded
	InfoHash    string
	Encoding    string
	Path        []string
	Seeders     int
	Leechers    int
	FilesCount  int
	Discussion  []TorrentComment
}

type TorrentComment struct {
	CommentID int
	User      webdata.User
	Date      string
	Content   string
}

type Peer struct {
	PeerID string
	IP     string
	Port   int
}
