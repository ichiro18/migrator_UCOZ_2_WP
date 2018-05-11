package wordpress

import (
	"fmt"

	"strconv"
	"strings"

	"time"

	"path"

	"errors"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var ImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Image interface",
	Long:  `Image interface`,
	Run: func(cmd *cobra.Command, args []string) {
		getImagesList()
		//clearPosts()
	},
}

func getImagesList() {
	imageListFromDB := []Post{}
	db := checkDB()
	err := db.Where("post_type = ?", "attachment").Find(&imageListFromDB)
	if err.GetErrors() != nil {
		fmt.Errorf("DB: %v", err.GetErrors())
	}
	if err.RecordNotFound() {
		fmt.Errorf("record not found")
	}
	color.Green("Len: %v", len(imageListFromDB))

	for _, image := range imageListFromDB {
		color.Red("id = 3339")
		color.Red("-----------")
		item := image
		color.Red("ID: %v", item.ID)
		date := checkDate(item)
		color.White("date: %v", date.String())

		postID, errno := getPostParentID(item.Guid)
		if errno != nil {
			fmt.Errorf("Can't get postID. ")
		}
		color.Red("Post: %v", postID)
		updateReferences(item.ID, date, postID)
		updateThumbnail(item.ID, postID)
	}
}

func checkDate(item Post) time.Time {
	// Date
	year := strconv.Itoa(item.Date.Year())
	var month string
	monthT := strconv.Itoa(int(item.Date.Month()))
	if len(monthT) == 1 {
		month = "0" + monthT
	} else {
		month = monthT
	}
	var day string
	day = strconv.Itoa(item.Date.Day())
	if len(day) == 1 {
		day = "0" + day
	}

	date := year + "-" + month + "-" + day

	pathFile := strings.TrimLeft(item.Guid, "http://u0500614.isp.regruhosting.ru/wp-content/uploads/")
	pathsFolders := strings.Split(pathFile, "/")
	currentYear := pathsFolders[0]
	currentMonth := pathsFolders[1]
	currentDate := currentYear + "-" + currentMonth + "-" + day

	if date != currentDate {
		color.Yellow("Date = %v, Current = %s", date, currentDate)
		color.Yellow("Updating...")
		date = currentDate
	}

	var hour string
	hour = strconv.Itoa(item.Date.Hour())
	if len(hour) == 1 {
		hour = "0" + hour
	}
	var minute string
	minute = strconv.Itoa(item.Date.Minute())
	if len(minute) == 1 {
		minute = "0" + minute
	}
	var second string
	second = strconv.Itoa(item.Date.Second())
	if len(second) == 1 {
		second = "0" + second
	}
	date = date + "T" + hour + ":" + minute + ":" + second + ".371Z"
	color.White("date: %v", date)
	const layout = "2006-01-02T15:04:05.000Z"
	successDate, err := time.Parse(layout, date)
	if err != nil {
		fmt.Errorf("Can't parse date. ")
	}

	return successDate
}

func getPostParentID(src string) (int, error) {
	config := Env.Config.GetStringMapString("wordpress")
	wpPath := checkFolder(config["path"])
	workPath := checkFolder(path.Join(wpPath, "work"))
	thumbPath := path.Join(workPath, "thumbnails.json")
	thumbListOld := []Thumbnail{}
	readThumbnails(thumbPath, &thumbListOld)

	srcRel := "/" + strings.TrimLeft(src, "http://u0500614.isp.regruhosting.ru")
	for _, thumb := range thumbListOld {
		if srcRel == thumb.ImagePath {
			return thumb.PostID, nil
		}
	}
	err := errors.New("not found")
	return 0, err
}

func updateReferences(id int, date time.Time, parentID int) {
	db := checkDB()

	post := Post{}
	err := db.Where("ID = ?", id).First(&post)
	if err.GetErrors() != nil {
		fmt.Errorf("DB: %v", err.GetErrors())
	}
	if err.RecordNotFound() {
		fmt.Errorf("record not found")
	}
	postNew := Post{}
	postNew.Date = date
	postNew.DateGmt = date
	postNew.Modified = date
	postNew.ModifiedGMT = date
	postNew.Parent = strconv.Itoa(parentID)
	tx := db.Begin()
	if err := tx.Model(&post).Where("ID = ?", id).UpdateColumns(postNew).Error; err != nil {
		tx.Rollback()
		color.Red("Cant update")
		fmt.Errorf("%v", err)
	}

	if err := tx.Commit().Error; err != nil {
		color.Red("Transaction fail")
	}
}

func updateThumbnail(imageID int, postID int) {
	db := checkDB()

	query := Thumbnail{
		PostID:    postID,
		MetaKey:   "_thumbnail_id",
		MetaValue: strconv.Itoa(imageID),
	}

	tx := db.Begin()
	if err := tx.Table("wp_postmeta").Create(&query).Error; err != nil {
		tx.Rollback()
		color.Red("Cant create")
		fmt.Errorf("%v", err)
	}

	if err := tx.Commit().Error; err != nil {
		color.Red("Transaction fail")
	}
}

// CLEAR ALL POSTS
func clearPosts() {
	imageListFromDB := []Post{}
	db := checkDB()
	err := db.Where("post_type = ?", "post").Delete(&imageListFromDB)
	if err.GetErrors() != nil {
		fmt.Errorf("DB: %v", err.GetErrors())
	}
	if err.RecordNotFound() {
		fmt.Errorf("record not found")
	}
	defer db.Close()
	color.Green("Len: %v", len(imageListFromDB))

}
