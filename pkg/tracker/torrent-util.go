package tracker

import (
	"database/sql"
	"net/http"

	"github.com/oskar13/mini-tracker/pkg/data"
	db "github.com/oskar13/mini-tracker/pkg/db"
	"github.com/zeebo/bencode"
)

// Check torrent access rights based on torrent_access_list table
func CheckTorrentAccessList(torrentID int, userID int) (bool, error) {
	var result int
	q := "SELECT torrent_access_list.user_ID FROM torrent_access_list WHERE torrent_access_list.torrent_ID = ? AND torrent_access_list.user_ID = ?"
	err := db.DB.QueryRow(q, torrentID, userID).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	} else if result == userID {
		return true, nil
	}
	return false, nil
}

func LoadPeers(torrentID int) ([]data.Peer, error) {

	var peerList []data.Peer
	q := "SELECT peer_id, ip, port FROM peers where torrent_ID = ?"

	rows, err := db.DB.Query(q, torrentID)
	if err != nil {
		return []data.Peer{}, err
	}
	defer rows.Close()
	for rows.Next() {

		var peer data.Peer
		err = rows.Scan(&peer.PeerID, &peer.IP, &peer.Port)
		if err != nil {
			// handle this error
			return []data.Peer{}, err

		}
		peerList = append(peerList, peer)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return []data.Peer{}, err
	}
	return peerList, nil

}

func EncodePeerListAndRespond(w http.ResponseWriter, interval int, peerList []data.Peer) error {

	var response = make(map[string]interface{})

	response["interval"] = interval

	var list []interface{}

	for _, v := range peerList {

		var peer = make(map[string]interface{})
		peer["ip"] = v.IP
		peer["port"] = v.Port
		peer["peer id"] = v.PeerID

		list = append(list, peer)

	}

	response["peers"] = list

	enc := bencode.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		return err
	}
	return nil
}
