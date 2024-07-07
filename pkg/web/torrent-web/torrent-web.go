package torrentweb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/tracker"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func LoadTorrentData(torrentID int, userID int) (webdata.TorrentWeb, error) {
	var user webdata.User
	var torrent webdata.TorrentWeb

	q := "SELECT torrents.created, torrents.name, torrents.size, torrents.anonymous, torrents.access_type, torrents.group_ID, torrents.upvotes, torrents.downvotes, torrents.description, torrents.info_hash, torrents.pieces, torrents.piece_length, torrents.path, users.user_ID, users.username, users.profile_pic, groups.group_name FROM torrents LEFT JOIN users ON torrents.user_ID = users.user_ID LEFT JOIN groups ON groups.group_ID = torrents.group_ID WHERE torrents.torrent_ID = ?"

	err := db.DB.QueryRow(q, torrentID).Scan(&torrent.Created, &torrent.Name, &torrent.Size, &torrent.Anonymous, &torrent.AccessType, &torrent.GroupID, &torrent.UpVotes, &torrent.DownVotes, &torrent.Description, &torrent.InfoHash, &torrent.Pieces, &torrent.PieceLength, &torrent.PathJSON, &user.UserID, &user.Username, &user.Cover, &torrent.GroupName)
	if err != nil {
		if err == sql.ErrNoRows {

			return webdata.TorrentWeb{}, errors.New("no torrent found")
		}
		return webdata.TorrentWeb{}, err
	}

	if !torrent.Anonymous {
		torrent.User = user
	}

	torrent.Discussion, err = LoadTorrentComments(torrentID)

	if err != nil {
		fmt.Println("Error loading comments : ", err)
	}

	//Check torrent access based on access_type field that overrides all else
	if torrent.AccessType == "Public" || torrent.AccessType == "WWW" || torrent.AccessType == "Link Only" {
		return torrent, nil
	}

	//Check if user is on direct access list of torrent
	listAccess, err := tracker.CheckTorrentAccessList(torrentID, userID)
	if err != nil {
		return webdata.TorrentWeb{}, err
	} else if listAccess {
		return torrent, nil
	}

	if torrent.GroupID != nil {
		//Check if user can access torrent by group
		if groups.LoadGroupAccess(userID, *torrent.GroupID) != "" {
			return torrent, nil
		}
	}

	return webdata.TorrentWeb{}, errors.New("no access")
}

func LoadTorrentComments(torrentID int) ([]webdata.TorrentComment, error) {

	var comments []webdata.TorrentComment

	q := `SELECT torrent_comments.comment_ID, torrent_comments.content, torrent_comments.date, users.user_ID, users.username, users.profile_pic FROM torrent_comments LEFT JOIN users ON users.user_ID = torrent_comments.user_ID WHERE torrent_comments.torrent_ID = ?`

	rows, err := db.DB.Query(q, torrentID)

	fmt.Println(err)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return []webdata.TorrentComment{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment webdata.TorrentComment
		var user webdata.User
		err := rows.Scan(&comment.CommentID, &comment.Content, &comment.Date, &user.UserID, &user.Username, &user.Cover)
		if err != nil {
			return []webdata.TorrentComment{}, err
		} else {
			comment.User = user
			comments = append(comments, comment)
		}
	}

	return comments, nil
}
