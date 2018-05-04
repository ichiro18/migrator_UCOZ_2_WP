package ucoz

import (
	"path/filepath"

	"fmt"

	"strings"

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
	},
}

type Site struct {
	News []NewItem
}

type NewItem struct {
	id       string
	catID    string
	year     string
	month    string
	day      string
	pending  string
	ontop    string
	com_may  string
	addtime  string
	num_com  string
	author   string
	title    string
	brief    string
	message  string
	attach   string
	files    string
	reads    string
	rating   string
	rate_num string
	rate_sum string
	rate_ip  string
	other1   string
	other2   string
	other3   string
	other4   string
	other5   string
	uid      string
	sbscr    string
	lastmod  string
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
	elems := strings.Split(string(file), "|")

	color.Red("SUM: %v", len(elems))
	ND(len(elems))
}

func ND(sum int) {
	for i := 1; i < sum; i++ {
		del := sum / i
		ost := sum % i
		if (ost == 0) && (del != 1) && (del != sum) {
			color.Red("DEL: %v", del)
		}
	}
}
