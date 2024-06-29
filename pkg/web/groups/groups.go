package groups

import (
	"github.com/oskar13/mini-tracker/pkg/db"
	"github.com/oskar13/mini-tracker/pkg/web/webdata"
)

func GetUserGroupsList(userID int, visibility string) []webdata.UserGroupListObject {
	var groupList []webdata.UserGroupListObject
	q := " SELECT group_members.group_ID, groups.group_name, groups.group_icon, group_roles.role_type, groups.group_visibility FROM group_members LEFT JOIN groups ON group_members.group_ID = groups.group_ID LEFT JOIN group_roles on group_roles.group_ID = group_members.group_ID WHERE group_members.user_ID = ? AND groups.group_visibility = ?"

	rows, err := db.DB.Query(q, userID, visibility)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var listItem webdata.UserGroupListObject
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
