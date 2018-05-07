package ucoz

import (
	"path/filepath"

	"fmt"

	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var NewsCmd = &cobra.Command{
	Use:   "news",
	Short: "News interface",
	Long:  `News interface`,
	Run: func(cmd *cobra.Command, args []string) {
		site := new(Site)
		site.parseNews(UcozFileStruct.Site.Path)
	},
}

type Site struct {
	News []NewItem
}

type NewItem struct {
	ID          string
	CAT_ID      string
	URL_YEAR    string
	URL_MONTH   string
	URL_DAY     string
	TITLE       string
	SNIPPET     string
	MESSAGE     string
	ATTACHMENTS string
	VIEWS_COUNT string
	URL_ALIAS   string
}

func (self *Site) parseNews(path string) {
	self.Load(path)
}

func (self *Site) Load(path string) {
	newsPath := filepath.Join(path, "news.txt")
	isEmptyNews, err := afero.IsEmpty(Env.FileSystem, newsPath)
	if err != nil {
		fmt.Errorf("Can't read news.txt file. ")
	}
	if isEmptyNews {
		fmt.Errorf("News is empty. ")
	}

	file, err := afero.ReadFile(Env.FileSystem, newsPath)
	if err != nil {
		fmt.Errorf("Can't read news.txt file. ")
	}
	// Удаляем пробелы с конца и начала файла
	data := strings.TrimSpace(string(file))
	//// Экранированные переносы строк
	data = strings.Replace(data, "\\\n", " ", -1)
	// Получаем новости
	elems := strings.Split(data, "\n")

	for _, val := range elems {
		dataItem := strings.Replace(val, "\\\t", "", -1)
		item := strings.Split(dataItem, "|")
		newsItem := NewItem{
			ID:        item[0],
			CAT_ID:    item[1],
			URL_YEAR:  item[2],
			URL_MONTH: item[3],
			URL_DAY:   item[4],
			TITLE:     item[11],
			MESSAGE:   item[13],
		}
		self.News = append(self.News, newsItem)
	}
}

func GetNews() {
	site := new(Site)
	site.parseNews(UcozFileStruct.Site.Path)
}
