package torrentweb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

// Load public torrents in user profile view or elsewhere, set flag to true to
func LoadUserTorrents(user_ID int, access_type []string) []webdata.TorrentWeb {

	var resultTorrents []webdata.TorrentWeb

	if len(access_type) == 0 {
		return []webdata.TorrentWeb{}
	}

	q := `SELECT torrents.torrent_ID, torrents.created, torrents.name, torrents.upvotes, torrents.downvotes, torrents.access_type, torrents.size
	      FROM torrents
	      WHERE torrents.user_ID = ? AND torrents.access_type IN (` + strings.Repeat("?,", len(access_type)-1) + `?)`

	args := make([]interface{}, 0, len(access_type)+1)
	args = append(args, user_ID)
	for _, at := range access_type {
		args = append(args, at)
	}

	rows, err := db.DB.Query(q, args...)

	fmt.Println(err)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return []webdata.TorrentWeb{}
	}
	defer rows.Close()

	for rows.Next() {
		var torrent webdata.TorrentWeb
		err := rows.Scan(&torrent.TorrentID, &torrent.Created, &torrent.Name, &torrent.UpVotes, &torrent.DownVotes, &torrent.AccessType, &torrent.Size)
		if err != nil {
			return []webdata.TorrentWeb{}
		} else {
			resultTorrents = append(resultTorrents, torrent)
		}
	}

	return resultTorrents
}
