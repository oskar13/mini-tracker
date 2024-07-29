package data

var TorrentCategories = map[int]string{
	100: "Audio",
	101: "Music",
	199: "Other",
	200: "Video",
	201: "Movies",
	299: "Other",
	300: "Applications",
	301: "Windows",
	399: "Other",
	400: "Games",
	401: "PC",
	499: "Other",
	600: "Other",
	601: "E-books",
	699: "Other",
}

// Do the simple job of separating category and sub category as IDs
// and returning their string representations.
// Takes child category id and returns parent category id,
// parent category name and child category name
func GetCategoryNameAndID(id int) (int, string, string) {
	parentCatID := id / 100 * 100
	return parentCatID, TorrentCategories[parentCatID], TorrentCategories[id]
}
