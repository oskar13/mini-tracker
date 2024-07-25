package data

type Peer struct {
	TorrentID int
	InfoHash  string
	PeerID    string
	IP        string
	Port      int
	Left      int64
}
