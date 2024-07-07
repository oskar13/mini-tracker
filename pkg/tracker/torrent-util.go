package tracker

import (
	"database/sql"

	"github.com/oskar13/mini-tracker/pkg/db"
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
