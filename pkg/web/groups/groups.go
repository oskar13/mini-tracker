package groups

import (
	"github.com/oskar13/mini-tracker/pkg/db"
)

func GetUserGroupsList(userID int, visibility string) []GroupInfo {
	var groupList []GroupInfo
	q := " SELECT group_members.group_ID, groups.group_name, groups.group_icon, group_roles.role_type, groups.group_visibility FROM group_members LEFT JOIN groups ON group_members.group_ID = groups.group_ID LEFT JOIN group_roles on group_roles.group_ID = group_members.group_ID WHERE group_members.user_ID = ? AND groups.group_visibility = ?"

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
