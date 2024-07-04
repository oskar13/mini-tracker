package groups

import (
	"database/sql"
	"fmt"

	"github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func GetUserGroupsList(userID int, visibility string) []GroupInfo {
	var groupList []GroupInfo
	q := " SELECT group_members.group_ID, groups.group_name, groups.group_icon, group_members.group_role, groups.group_visibility FROM group_members LEFT JOIN groups ON group_members.group_ID = groups.group_ID WHERE group_members.user_ID = ? AND groups.group_visibility = ?"

	rows, err := db.DB.Query(q, userID, visibility)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var listItem GroupInfo
		err = rows.Scan(&listItem.GroupID, &listItem.GroupName, &listItem.GroupIcon, &listItem.GroupRole, &listItem.GroupVisibility)
		if err != nil {
			// handle this error
			panic(err)
		}

		groupList = append(groupList, listItem)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return groupList

}

func ListPublicGroups() []GroupInfo {
	var groupList []GroupInfo
	q := " SELECT groups.group_ID, groups.group_name, groups.group_icon, groups.tagline, groups.join_type FROM groups WHERE groups.group_visibility ='Public'"

	rows, err := db.DB.Query(q)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var listItem GroupInfo
		err = rows.Scan(&listItem.GroupID, &listItem.GroupName, &listItem.GroupIcon, &listItem.GroupTagline, &listItem.GroupJoinType)
		if err != nil {
			// handle this error
			panic(err)
		}

		groupList = append(groupList, listItem)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return groupList

}

// List last x posts from public groups and user groups
func GetCommunityUpdates(userID int, count int) []GroupPost {
	var postList []GroupPost
	q := "SELECT group_posts.group_name, group_posts.group_ID,  group_posts.post_ID, group_posts.title, group_posts.date FROM group_posts INNER JOIN group_members ON group_members.group_ID = group_posts.group_ID INNER JOIN groups ON groups.group_ID = group_posts.group_ID WHERE group_members.user_ID = ? ORDER BY group_posts.date DESC LIMIT ?"

	rows, err := db.DB.Query(q, userID, count)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var group GroupInfo
		var post GroupPost

		err = rows.Scan(&group.GroupName, &group.GroupID, &post.PostID, &post.Title, &post.Date)
		if err != nil {
			// handle this error
			panic(err)
		}

		post.Group = group
		postList = append(postList, post)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return postList

}

// Load user group role or return "" if they have no right to browse the group
// If group is public then add the role as Guest
func LoadGroupAccess(userID int, groupID int) string {

	var role string

	q := `SELECT group_members.group_role FROM group_members WHERE user_ID = ? and group_ID = ?`
	err := db.DB.QueryRow(q, userID, groupID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			//User is not a member of the group check if the group is public and give guest role
			if CheckIfGroupPublic(groupID) {
				return "Guest"
			} else {
				return ""
			}
		}
		fmt.Println(err)
		return ""
	}

	return role

}

func CheckIfGroupPublic(groupID int) bool {
	var group int

	q := `SELECT group_ID FROM groups WHERE group_ID = ? AND group_visibility = "Public"`
	err := db.DB.QueryRow(q, groupID).Scan(&group)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		fmt.Println(err)
		return false
	}

	return true
}

func LoadGroupInfo(groupID int) GroupInfo {
	var group GroupInfo

	q := `SELECT group_ID, group_name, group_icon, group_visibility, tagline, tagline FROM groups WHERE group_ID = ?`
	err := db.DB.QueryRow(q, groupID).Scan(&group.GroupID, &group.GroupName, &group.GroupIcon, &group.GroupVisibility, &group.GroupTagline, &group.GroupDescription)
	if err != nil {
		if err == sql.ErrNoRows {
			return GroupInfo{}
		}
		fmt.Println(err)
		return GroupInfo{}
	}

	return group
}

func LoadGroupPostsList(groupID int) []GroupPost {
	var postList []GroupPost

	q := "SELECT  group_posts.post_ID, group_posts.title, group_posts.content, group_posts.date, group_posts.sticky, group_posts.user_ID, group_posts.username, group_posts.profile_pic FROM group_posts WHERE group_posts.group_ID = ? ORDER BY group_posts.sticky DESC ,group_posts.date DESC"

	rows, err := db.DB.Query(q, groupID)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {

		var post GroupPost
		var user webdata.User

		err = rows.Scan(&post.PostID, &post.Title, &post.Content, &post.Date, &post.Sticky, &user.UserID, &user.Username, &user.Cover)
		if err != nil {
			// handle this error
			panic(err)
		}
		post.User = user

		postList = append(postList, post)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return postList
}

func LoadGroupPost(groupID int, postID int) GroupPost {
	var post GroupPost
	var user webdata.User

	q := "SELECT group_posts.post_ID, group_posts.title, group_posts.content, group_posts.date , group_posts.updated , group_posts.sticky, group_posts.user_ID, group_posts.username, group_posts.profile_pic FROM group_posts WHERE group_posts.group_ID = ? AND group_posts.post_ID = ?"

	err := db.DB.QueryRow(q, groupID, postID).Scan(&post.PostID, &post.Title, &post.Content, &post.Date, &post.Updated, &post.Sticky, &user.UserID, &user.Username, &user.Cover)
	if err != nil {
		if err == sql.ErrNoRows {
			return GroupPost{}
		}
		fmt.Println(err)
		return GroupPost{}
	} else {
		//load replies
		post.User = user
		post.Replies = LoadGroupPostsReplies(postID)
	}

	return post
}

func LoadGroupPostsReplies(postID int) []GroupPostReply {
	var replyList []GroupPostReply

	q := "SELECT group_post_replies.reply_ID, group_post_replies.user_ID, users.username, users.profile_pic, group_post_replies.date, group_post_replies.updated, group_post_replies.content FROM group_post_replies LEFT JOIN users ON group_post_replies.user_ID = users.user_ID WHERE  group_post_replies.post_ID = ?"

	rows, err := db.DB.Query(q, postID)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {

		var reply GroupPostReply
		var user webdata.User

		err = rows.Scan(&reply.ReplyID, &user.UserID, &user.Username, &user.Cover, &reply.Date, &reply.Updated, &reply.Content)
		if err != nil {
			// handle this error
			panic(err)
		}
		reply.User = user

		replyList = append(replyList, reply)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return replyList
}
