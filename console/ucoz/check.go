package ucoz

import (
	"fmt"
	"strings"

	"path/filepath"

	"encoding/json"

	table "github.com/crackcomm/go-clitable"
	"github.com/fatih/color"
	"github.com/fatih/structs"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type ucozFileStruct struct {
	Site        *StructItem `default:"folder:_s1;description:Базы данных категорий, материалов, пользователей"`
	Board       *StructItem `default:"folder:_bd;description:Файлы модуля «Доска объявлений» (board)"`
	Blog        *StructItem `default:"folder:_bl;description:Файлы модуля «Блог» (blog)"`
	SiteCatalog *StructItem `default:"folder:_dr;description:Файлы модуля «Каталог сайтов» (dir)"`
	FAQ         *StructItem `default:"folder:_fq;description:Файлы модуля «FAQ» (faq)"`
	Forum       *StructItem `default:"folder:_fr;description:Файлы модуля «Форум» (forum)"`
	FileCatalog *StructItem `default:"folder:_ld;description:Файлы модуля «Каталог файлов» (load)"`
	News        *StructItem `default:"folder:_nw;description:Файлы модуля «Новости» (news)"`
	Photo       *StructItem `default:"folder:_ph;description:Файлы модуля «Фотоальбом» (photo)"`
	Article     *StructItem `default:"folder:_pu;description:Файлы модуля «Каталог статей» (publ)"`
	Games       *StructItem `default:"folder:_sf;description:Файлы модуля «Онлайн-игры» (stuff)"`
	Shop        *StructItem `default:"folder:_sh;description:Файлы модуля «Интернет-магазин» (shop)"`
	Pages       *StructItem `default:"folder:_si;description:Файлы модуля «Страницы» (index)"`
	Styles      *StructItem `default:"folder:_st;description:Файлы стилей (my.css)"`
	Video       *StructItem `default:"folder:_vi;description:Файлы модуля «Видео» (video)"`
}

type StructItem struct {
	Name        string
	Status      bool
	Description string
	Folder      string
	Path        string
}

var Env *services.Env
var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate ucoz files",
	Long:  `Validate ucoz files`,
	Run: func(cmd *cobra.Command, args []string) {
		config := Env.Config.GetStringMapString("ucoz")
		ucozPath := checkFolder(config["path"])
		structure := createDefaultStruct()
		UcozFileStruct = checkBackup(filepath.Join(ucozPath, "backup"), structure, true)
		saveFileStruct(UcozFileStruct, filepath.Join(ucozPath, "work"))
	},
}
var UcozFileStruct *ucozFileStruct

func NewUcozStructure() *ucozFileStruct {
	config := Env.Config.GetStringMapString("ucoz")
	ucozPath := checkFolder(config["path"])
	structure := createDefaultStruct()
	UcozFileStruct = checkBackup(filepath.Join(ucozPath, "backup"), structure, false)
	saveFileStruct(UcozFileStruct, filepath.Join(ucozPath, "work"))
	return UcozFileStruct
}

func checkFolder(path string) string {
	isExistFolder, err := afero.DirExists(Env.FileSystem, path)
	if err != nil {
		fmt.Errorf("Can't check config path. ")
	}
	if !isExistFolder {
		color.Yellow("Ucoz folder is exist. Creating...")
		Env.FileSystem.Mkdir(path, 0755)
	}
	return path
}

func checkBackup(path string, structure *ucozFileStruct, print bool) *ucozFileStruct {
	isExistFolder, err := afero.DirExists(Env.FileSystem, path)
	if err != nil {
		fmt.Errorf("Can't check config path. ")
	}
	if !isExistFolder {
		color.Yellow("Backup folder is exist. Creating...")
		Env.FileSystem.Mkdir(path, 0755)
	}

	files, err := afero.ReadDir(Env.FileSystem, path)
	if err != nil {
		fmt.Errorf("Can't read backup dir. ")
	}

	tableOut := table.New([]string{"Name", "Status", "Description", "Folder", "Path"})

	ucoz := structs.New(structure)
	for _, field := range ucoz.Fields() {
		s := structs.New(field.Value())
		s.Field("Path").Set("-")
		for _, val := range files {
			// Only folder
			if val.IsDir() {
				if val.Name() == s.Field("Folder").Value() {
					s.Field("Status").Set(true)
					dirPath := filepath.Join(path, val.Name())
					isEmptyChild, err := afero.IsEmpty(Env.FileSystem, dirPath)
					if err != nil {
						fmt.Errorf("Can't read %v dir. ", dirPath)
					}
					if !isEmptyChild {
						s.Field("Path").Set(dirPath)
					}
				}
			}
		}

		m := s.Map()
		tableOut.AddRow(m)
	}

	if print {
		tableOut.Print()
	}
	return structure
}

func createDefaultStruct() *ucozFileStruct {
	ucozS := new(ucozFileStruct)
	s := structs.New(ucozS)

	m := s.Map()
	for index, _ := range m {
		field := s.Field(index)
		item := new(StructItem)
		defaultValuesString := field.Tag("default")
		type defaults struct {
			Folder      string
			Description string
		}
		defaultValuesArr := strings.Split(defaultValuesString, ";")
		dv := defaults{}
		for _, val := range defaultValuesArr {
			t := strings.Split(val, ":")
			name := strings.Title(t[0])
			value := strings.Title(t[1])
			if name == "Folder" {
				dv.Folder = value
			}
			if name == "Description" {
				dv.Description = value
			}
		}
		// Name
		item.Name = field.Name()
		// Description
		item.Description = dv.Description
		// Folder
		item.Folder = dv.Folder
		// Status
		item.Status = false

		switch field.Name() {
		case "Site":
			ucozS.Site = item
			break
		case "Board":
			ucozS.Board = item
			break
		case "Blog":
			ucozS.Blog = item
			break
		case "SiteCatalog":
			ucozS.SiteCatalog = item
			break
		case "FAQ":
			ucozS.FAQ = item
			break
		case "Forum":
			ucozS.Forum = item
			break
		case "FileCatalog":
			ucozS.FileCatalog = item
			break
		case "News":
			ucozS.News = item
			break
		case "Photo":
			ucozS.Photo = item
			break
		case "Article":
			ucozS.Article = item
			break
		case "Games":
			ucozS.Games = item
			break
		case "Shop":
			ucozS.Shop = item
			break
		case "Pages":
			ucozS.Pages = item
			break
		case "Styles":
			ucozS.Styles = item
			break
		case "Video":
			ucozS.Video = item
			break
		}
	}

	return ucozS
}

func saveFileStruct(structure *ucozFileStruct, path string) {
	data := structs.New(structure)
	dataMap := data.Map()
	dataJson, err := json.MarshalIndent(dataMap, " ", "  ")
	if err != nil {
		fmt.Errorf("Can't convert data to json. ")
	}
	afero.WriteFile(Env.FileSystem, filepath.Join(path, "struct.json"), dataJson, 0755)
}
