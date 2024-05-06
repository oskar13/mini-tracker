package data

// Data to hold when user uploads a torrent from site
type TorrentUpload struct {
	Announce     string
	Name         string
	Comment      string
	PieceLength  int64
	Pieces       []byte
	Private      bool
	CreatedBy    string
	CreationDate int64
	Encoding     string
	Path         []string
}
