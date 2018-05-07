package wordpress

import (
	"fmt"

	"strconv"

	"github.com/fatih/color"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

type Post struct {
	ID              int      `gorm:"column:ID;type:bigint(20);AUTO_INCREMENT"`
	Author          int      `gorm:"column:post_author;type:bigint(20)"`
	Date            *[]uint8 `gorm:"column:post_date;datetime"`
	DateGmt         *[]uint8 `gorm:"column:post_date_gmt;datetime"`
	Content         string   `gorm:"column:post_content;type:longtext"`
	Title           string   `gorm:"column:post_title;type:text"`
	Excerpt         string   `gorm:"column:post_excerpt;type:text"`
	Status          string   `gorm:"column:post_status;type:varchar(20)"`
	Comment         string   `gorm:"column:comment_status;type:varchar(20)"`
	Ping            string   `gorm:"column:ping_status;type:varchar(20)"`
	Password        string   `gorm:"column:post_password;type:varchar(255)"`
	Name            string   `gorm:"column:post_name;type:varchar(200)"`
	ToPing          string   `gorm:"column:to_ping;type:text"`
	Pinged          string   `gorm:"column:pinged;type:text"`
	Modified        *[]uint8 `gorm:"column:post_modified;type:datetime"`
	ModifiedGMT     *[]uint8 `gorm:"column:post_modified_gmt;type:datetime"`
	ContentFiltered string   `gorm:"column:post_content_filtered;type:longtext"`
	Parent          string   `gorm:"column:post_parent;type:bigint(20)"`
	Guid            string   `gorm:"column:guid;type:varchar(255)"`
	MenuOrder       int      `gorm:"column:menu_order;type:int(11)"`
	Type            string   `gorm:"column:post_type;type:varchar(20);"`
	MimeType        string   `gorm:"column:post_mime_type;type:varchar(100);"`
	CommentCount    uint64   `gorm:"column:comment_count;type:bigint(20)"`
}

var (
	list bool
)

var PostCmd = &cobra.Command{
	Use:   "post [options]",
	Short: "Post interface",
	Long:  `Post interface`,
	Run: func(cmd *cobra.Command, args []string) {
		if list {
			getPostList()
		}
		getLastPostID()
	},
}

func init() {
	PostCmd.Flags().BoolVarP(&list, "list", "l", false, "See post list")
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

func (p Post) TableName() string {
	return "wp_posts"
}
