package wordpress

import (
	"fmt"

	"strconv"

	"encoding/json"
	"path/filepath"

	"time"

	"path"

	"strings"

	"github.com/fatih/color"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
	"github.com/ichiro18/migrator_UCOZ_2_WP/console/ucoz"
	"github.com/jinzhu/gorm"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type Post struct {
	ID              int       `gorm:"column:ID;type:bigint(20);AUTO_INCREMENT"`
	Author          int       `gorm:"column:post_author;type:bigint(20)"`
	Date            time.Time `gorm:"column:post_date;datetime"`
	DateGmt         time.Time `gorm:"column:post_date_gmt;datetime"`
	Content         string    `gorm:"column:post_content;type:longtext"`
	Title           string    `gorm:"column:post_title;type:text"`
	Excerpt         string    `gorm:"column:post_excerpt;type:text"`
	Status          string    `gorm:"column:post_status;type:varchar(20)"`
	Comment         string    `gorm:"column:comment_status;type:varchar(20)"`
	Ping            string    `gorm:"column:ping_status;type:varchar(20)"`
	Password        string    `gorm:"column:post_password;type:varchar(255)"`
	Name            string    `gorm:"column:post_name;type:varchar(200)"`
	ToPing          string    `gorm:"column:to_ping;type:text"`
	Pinged          string    `gorm:"column:pinged;type:text"`
	Modified        time.Time `gorm:"column:post_modified;type:datetime"`
	ModifiedGMT     time.Time `gorm:"column:post_modified_gmt;type:datetime"`
	ContentFiltered string    `gorm:"column:post_content_filtered;type:longtext"`
	Parent          string    `gorm:"column:post_parent;type:bigint(20)"`
	Guid            string    `gorm:"column:guid;type:varchar(255)"`
	MenuOrder       int       `gorm:"column:menu_order;type:int(11)"`
	Type            string    `gorm:"column:post_type;type:varchar(20);"`
	MimeType        string    `gorm:"column:post_mime_type;type:varchar(100);"`
	CommentCount    uint64    `gorm:"column:comment_count;type:bigint(20)"`
}

var (
	list  bool
	clear bool
)

var PostCmd = &cobra.Command{
	Use:   "post [options]",
	Short: "Post interface",
	Long:  `Post interface`,
	Run: func(cmd *cobra.Command, args []string) {
		if list {
			getPostList()
		}
		if clear {
			clearPost()
		}
		if !list && !clear {
			uploadNews()
		}
	},
}

func init() {
	PostCmd.Flags().BoolVarP(&list, "list", "l", false, "See post list")
	PostCmd.Flags().BoolVarP(&clear, "clear", "c", false, "Clear post list")
}

func getPostList() {
	db := checkDB()
	postList := []Post{}
	err := db.Where("post_type = ?", "post").Find(&postList)
	if err.GetErrors() != nil {
		fmt.Errorf("DB: %v", err.GetErrors())
	}
	if err.RecordNotFound() {
		fmt.Errorf("record not found")
	}

	for _, val := range postList {
		color.Yellow("%s - %v", strconv.Itoa(val.ID), val.Title)
	}

	defer db.Close()
}

func getLastPostID() int {
	db := checkDB()
	lastPost := Post{}

	err := db.Order("ID desc").First(&lastPost)
	if err.GetErrors() != nil {
		fmt.Errorf("DB: %v", err.GetErrors())
	}
	if err.RecordNotFound() {
		fmt.Errorf("record not found")
	}

	return lastPost.ID
}

func checkDB() *gorm.DB {
	if Env.Database == nil {
		color.Red("DB connect not found")
		Env.Database = services.NewConnectORM(Env.Config.GetStringMapString("wordpress"))
	}
	db := Env.Database
	if !db.HasTable("wp_posts") {
		fmt.Errorf("Table 'wp_posts' not exist. ")
	}

	return db
}

func uploadNews() {
	str := ucoz.Site{}
	// Get Ucoz News
	news := ucoz.GetNews(&str)
	postList := []Post{}
	lastID := getLastPostID()
	var postIDs []int

	db := checkDB()
	tx := db.Begin()
	for _, val := range news {
		lastID = lastID + 1
		post := convertUcozNewToWordpressPost(lastID, &val)
		postList = append(postList, post)
		postIDs = append(postIDs, post.ID)
		color.Yellow("create post ID=%v", post.ID)
		if err := tx.Create(&post).Error; err != nil {
			tx.Rollback()
			fmt.Errorf("Can't create post: %v", err.Error())
		}
	}

	err := tx.Commit().Error
	if err != nil {
		fmt.Errorf("Can't create posts: %v", err.Error())
	}
	// SaveFile
	config := Env.Config.GetStringMapString("wordpress")
	wpPath := checkFolder(config["path"])
	workPath := checkFolder(path.Join(wpPath, "work"))
	savePosts(&postList, workPath)
	saveIDs(&postIDs, workPath)
}

func convertUcozNewToWordpressPost(startID int, newItem *ucoz.NewItem) Post {
	var month string
	if len(newItem.URL_MONTH) == 1 {
		month = "0" + newItem.URL_MONTH
	} else {
		month = newItem.URL_MONTH
	}
	var day string
	if len(newItem.URL_DAY) == 1 {
		day = "0" + newItem.URL_DAY
	} else {
		day = newItem.URL_DAY
	}
	str := newItem.URL_YEAR + "-" + month + "-" + day
	date, err := time.Parse("2006-01-02", str)
	if err != nil {
		fmt.Errorf("Can't parse date")
	}
	post := Post{
		ID:           startID,
		Author:       1,
		Date:         date,
		DateGmt:      date,
		Content:      newItem.MESSAGE,
		Title:        newItem.TITLE,
		Status:       "publish",
		Type:         "post",
		Comment:      "open",
		Ping:         "open",
		Name:         newItem.TITLE,
		Modified:     date,
		ModifiedGMT:  date,
		Guid:         "http://u0500614.isp.regruhosting.ru/?p=" + strconv.Itoa(startID),
		CommentCount: 0,
	}

	return post
}

func savePosts(posts *[]Post, path string) {
	dataJson, err := json.MarshalIndent(posts, " ", "  ")
	if err != nil {
		fmt.Errorf("Can't convert data to json. ")
	}
	err = afero.WriteFile(Env.FileSystem, filepath.Join(path, "posts.json"), dataJson, 0755)
	if err != nil {
		fmt.Errorf("Can't create file. ")
	}
	color.Green("File '%v' created", filepath.Join(path, "posts.json"))
}

func saveIDs(posts *[]int, path string) {
	dataJson, err := json.MarshalIndent(posts, " ", "  ")
	if err != nil {
		fmt.Errorf("Can't convert data to json. ")
	}
	err = afero.WriteFile(Env.FileSystem, filepath.Join(path, "ids.json"), dataJson, 0755)
	if err != nil {
		fmt.Errorf("Can't create file. ")
	}
	color.Green("File '%v' created", filepath.Join(path, "ids.json"))
}

func clearPost() {
	// Open file
	config := Env.Config.GetStringMapString("wordpress")
	wpPath := checkFolder(config["path"])
	workPath := checkFolder(path.Join(wpPath, "work"))
	filePath := path.Join(workPath, "ids.json")
	file, err := afero.ReadFile(Env.FileSystem, filePath)
	if err != nil {
		fmt.Errorf("Can't read file IDs. ")
	}
	data := string(file)
	data = strings.Trim(data, "[")
	data = strings.Trim(data, "]")
	data = strings.Replace(data, "\n", " ", -1)
	data = strings.Trim(data, " ")
	dataArr := strings.Split(data, ",")
	db := checkDB()
	tx := db.Begin()
	for _, val := range dataArr {
		val = strings.Trim(val, " ")
		in, err := strconv.Atoi(val)
		if err != nil {
			fmt.Errorf("Can't convert string to int. ")
		}
		post := Post{
			ID: int(in),
		}
		color.Yellow("delete post ID=%v", post.ID)
	}
	err = tx.Commit().Error
	if err != nil {
		fmt.Errorf("Can't create posts: %v", err.Error())
	}
}

func (p Post) TableName() string {
	return "wp_posts"
}
