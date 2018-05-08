package wordpress

import (
	"fmt"
	"regexp"

	"strconv"

	"encoding/json"
	"path/filepath"

	"time"

	"path"

	"strings"

	"github.com/fatih/color"
	"github.com/fiam/gounidecode/unidecode"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
	"github.com/ichiro18/migrator_UCOZ_2_WP/console/ucoz"
	"github.com/jinzhu/gorm"
	"github.com/schollz/progressbar"
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

	count := len(news)
	bar := progressbar.New(count)
	for _, val := range news {
		bar.Add(1)
		lastID = lastID + 1
		post := convertUcozNewToWordpressPost(lastID, &val)
		postList = append(postList, post)
		postIDs = append(postIDs, post.ID)
		color.Yellow("create post ID=%v \n", post.ID)
		if err := tx.Create(&post).Error; err != nil {
			tx.Rollback()
			fmt.Errorf("Can't create post: %v", err.Error())
		}
		time.Sleep(100 * time.Millisecond)
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
		Content:      updateMediaPath(newItem),
		Title:        newItem.TITLE,
		Status:       "publish",
		Type:         "post",
		Comment:      "open",
		Ping:         "open",
		Name:         translite(newItem.TITLE),
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

func translite(str string) string {
	// Русские буквы в латиницу
	res := unidecode.Unidecode(str)
	// в нижний регистр
	res = strings.ToLower(res)
	res = strings.Replace(res, " - ", " ", -1)
	res = strings.Replace(res, " ", "-", -1)
	// убираем ненужные символы
	reg, err := regexp.Compile("[^a-z,A-Z,0-9,-]+")
	if err != nil {
		fmt.Errorf("Can't remove symbols. ")
	}
	res = reg.ReplaceAllString(res, "")
	return res
}

func updateMediaPath(data *ucoz.NewItem) string {
	var res string
	item := data.MESSAGE
	indexMedia := strings.Count(item, `src="`)
	if indexMedia != -1 {
		for i := 0; i < indexMedia; i++ {
			index := strings.Index(item, `src="`)
			srcIndex := index + 5
			var src string
			var srcIndexStart int
			var srcIndexEnd int
			if srcIndex != -1 {
				srcIndexEnd = strings.Index(item[srcIndex:], string('"'))
				srcPart := item[srcIndex:]
				srcIndexStart = srcIndex
				src = srcPart[:srcIndexEnd]
			}
			if src != "" {
				hasUrl := strings.Index(src, "http")
				hasLocal := strings.Index(src, "c:/TMP")
				isWpPath := strings.Index(src, "wp-content")
				if hasUrl == -1 && hasLocal == -1 && isWpPath == -1 {
					newSrc := copyMedia(src, data)

					// Replace String
					before := item[:srcIndexStart]
					strStart := item[srcIndexStart:]
					after := strStart[srcIndexEnd:]
					res = before + newSrc + after
				}
			}
			if res == "" {
				fmt.Errorf("Can't update media path. ")
			}
			item = res
		}

	} else {
		res = data.MESSAGE
	}

	return res
}
func copyMedia(mediaPath string, item *ucoz.NewItem) string {
	ucozPath := Env.Config.GetStringMapString("ucoz")
	filePath := path.Join(ucozPath["path"], "backup", mediaPath)
	_, fileName := path.Split(filePath)
	isEmptyFile, err := afero.IsEmpty(Env.FileSystem, filePath)
	if err != nil {
		fmt.Errorf("Can't check file. ")
	}

	if isEmptyFile {
		fmt.Errorf("File is not exist. ")
	}
	wpConfig := Env.Config.GetStringMapString("wordpress")
	wpPath := wpConfig["path"]
	contentFolder := checkFolder(path.Join(wpPath, "wp-content"))
	uploadFolder := checkFolder(path.Join(contentFolder, "uploads"))
	yearFolder := checkFolder(path.Join(uploadFolder, item.URL_YEAR))
	var month string
	if len(item.URL_MONTH) == 1 {
		month = "0" + item.URL_MONTH
	} else {
		month = item.URL_MONTH
	}
	postFolder := checkFolder(path.Join(yearFolder, month))

	// Copy mediaFile
	image, err := afero.ReadFile(Env.FileSystem, filePath)
	resultFilePath := path.Join(postFolder, fileName)
	err = afero.WriteFile(Env.FileSystem, resultFilePath, image, 0777)
	if err != nil {
		fmt.Errorf("Can't copy file. ")
	}

	// Save relative Path for image
	urlPath := strings.Replace(resultFilePath, "wordpress", "", 1)

	return urlPath
}
