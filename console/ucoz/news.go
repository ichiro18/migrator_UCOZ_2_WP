package ucoz

import (
	"path/filepath"

	"fmt"

	"strings"

	"encoding/json"

	"github.com/fatih/color"
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
		config := Env.Config.GetStringMapString("ucoz")
		ucozPath := checkFolder(config["path"])
		site.SaveNews(UcozFileStruct, filepath.Join(ucozPath, "work"))
	},
}

type Site struct {
	News []NewItem
}

type NewItem struct {
	ID        string
	CAT_ID    string
	URL_YEAR  string
	URL_MONTH string
	URL_DAY   string
	TITLE     string
	MESSAGE   string
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
			MESSAGE:   stripHtmlComment(item[13]),
		}
		self.News = append(self.News, newsItem)
	}
}

func (self *Site) SaveNews(structure *ucozFileStruct, path string) {
	if self.News == nil {
		self.Load(structure.Site.Path)
	}
	data := self.News
	dataJson, err := json.MarshalIndent(data, " ", "  ")
	if err != nil {
		fmt.Errorf("Can't convert data to json. ")
	}
	afero.WriteFile(Env.FileSystem, filepath.Join(path, "news.json"), dataJson, 0755)
	color.Green("File '%v' created", filepath.Join(path, "news.json"))
}

func stripHtmlComment(html string) string {
	count := strings.Count(html, "<!--")
	if count != 0 {
		for i := 1; i <= count; i++ {
			startIndex := strings.Index(html, "<!--")
			endIndex := strings.Index(html, "-->") + 3
			html = strings.Replace(html, html[startIndex:endIndex], "", 1)
		}
	}
	return html
}

func GetNews(self *Site) []NewItem {
	if self.News == nil {
		self.Load(UcozFileStruct.Site.Path)
	}

	return self.News
}
