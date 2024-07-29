package torrentweb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/oskar13/mini-tracker/pkg/data"
	"github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/tracker"
	"github.com/oskar13/mini-tracker/pkg/web/groups"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func LoadTorrentData(torrentUuid string, userID int) (webdata.TorrentWeb, error) {
	var user webdata.User
	var torrent webdata.TorrentWeb

	q := "SELECT torrents.torrent_ID, torrents.uuid, torrents.uploaded, torrents.name, torrents.size, torrents.anonymous, torrents.access_type, torrents.group_ID, torrents.upvotes, torrents.downvotes, torrents.description, torrents.info_hash, torrents.pieces, torrents.piece_length, torrents.path, users.user_ID, users.username, users.profile_pic, groups.group_name, torrents.category_ID FROM torrents LEFT JOIN users ON torrents.user_ID = users.user_ID LEFT JOIN groups ON groups.group_ID = torrents.group_ID WHERE torrents.uuid = ?"

	err := db.DB.QueryRow(q, torrentUuid).Scan(&torrent.TorrentID, &torrent.Uuid, &torrent.Uploaded, &torrent.Name, &torrent.Size, &torrent.Anonymous, &torrent.AccessType, &torrent.GroupID, &torrent.UpVotes, &torrent.DownVotes, &torrent.Description, &torrent.InfoHash, &torrent.Pieces, &torrent.PieceLength, &torrent.PathJSON, &user.UserID, &user.Username, &user.Cover, &torrent.GroupName, &torrent.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {

			return webdata.TorrentWeb{}, errors.New("no torrent found")
		}
		return webdata.TorrentWeb{}, err
	}

	if !torrent.Anonymous {
		torrent.User = user
	}

	torrent.ParentCategoryID, torrent.ParentCategory, torrent.Category = data.GetCategoryNameAndID(torrent.CategoryID)

	torrent.Discussion, err = LoadTorrentComments(torrent.TorrentID)

	if err != nil {
		fmt.Println("Error loading comments : ", err)
	}

	//Check torrent access based on access_type field that overrides all else
	if torrent.AccessType == "Public Listed" || torrent.AccessType == "Public Unlisted" {
		return torrent, nil
	}

	//Check if user is on direct access list of torrent
	listAccess, err := tracker.CheckTorrentAccessList(torrent.TorrentID, userID)
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

func LoadTorrentInfoField(torrentID int) ([]byte, error) {
	var torrent webdata.TorrentWeb

	q := "SELECT info_field FROM torrents WHERE torrent_ID = ?"

	err := db.DB.QueryRow(q, torrentID).Scan(&torrent.InfoField)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, errors.New("no torrent found")
		}
		return nil, err
	}

	return torrent.InfoField, nil
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

	q := `SELECT torrents.torrent_ID, torrents.uploaded, torrents.name, torrents.upvotes, torrents.downvotes, torrents.access_type, torrents.size, torrents.uuid
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
		err := rows.Scan(&torrent.TorrentID, &torrent.Uploaded, &torrent.Name, &torrent.UpVotes, &torrent.DownVotes, &torrent.AccessType, &torrent.Size, &torrent.Uuid)
		if err != nil {
			return []webdata.TorrentWeb{}
		} else {
			resultTorrents = append(resultTorrents, torrent)
		}
	}

	return resultTorrents
}

// Create a database entry for a torrent that can be used to create the torrent file and return uuid
func CreateTorrentEntry(torrent webdata.TorrentWeb, userID int) (string, error) {

	if torrent.Announce == "" {
		return "", errors.New("announce empty")
	} else if torrent.InfoField == nil {
		return "", errors.New("info field required")
	}

	torrent.Uuid = uuid.New().String()

	stmt, err := db.DB.Prepare("INSERT INTO torrents (user_ID,name,size,access_type,description,info_hash,info_field,uuid,category_ID) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()
	fmt.Println(torrent.Uuid)

	_, err = stmt.Exec(userID, torrent.Name, torrent.Size, torrent.AccessType, torrent.Description, torrent.InfoHash, torrent.InfoField, torrent.Uuid, torrent.CategoryID)

	if err != nil {
		return "", fmt.Errorf("error preparing statement: %v", err)
	}
	return torrent.Uuid, err

}
