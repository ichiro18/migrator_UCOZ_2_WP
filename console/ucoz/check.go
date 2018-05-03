package ucoz

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/ichiro18/migrator_UCOZ_2_WP/common/services"
	"github.com/spf13/afero"
	"github.com/fatih/color"
)

type ucozFileStruct struct {
	Site    		StructItem `default: Name>_s1`	//Базы данных категорий, материалов, пользователей 		"_s1"
	Board 			StructItem `default: Name>_bd`  //Файлы модуля «Доска объявлений» (board) 				"_bd"
	Blog			StructItem `default: Name>_bl`  //Файлы модуля «Блог» (blog) 							"_bl"
	SiteCatalog		StructItem `default: Name>_dr`  //Файлы модуля «Каталог сайтов» (dir)					"_dr"
	FAQ				StructItem `default: Name>_fq`	//Файлы модуля «FAQ» (faq)								"_fq"
	Forum			StructItem `default: Name>_fr`	//Файлы модуля «Форум» (forum)							"_fr"
	FileCatalog		StructItem `default: Name>_ld`	//Файлы модуля «Каталог файлов» (load)					"_ld"
	News			StructItem `default: Name>_nw`	//Файлы модуля «Новости» (news)							"_nw"
	Photo			StructItem `default: Name>_ph`	//Файлы модуля «Фотоальбом» (photo)						"_ph"
	Article			StructItem `default: Name>_pu`	//Файлы модуля «Каталог статей» (publ)					"_pu"
	Games			StructItem `default: Name>_sf`	//Файлы модуля «Онлайн-игры» (stuff)					"_sf"
	Shop			StructItem `default: Name>_sh`	//Файлы модуля «Интернет-магазин» (shop)				"_sh"
	Pages			StructItem `default: Name>_si`	//Файлы модуля «Страницы» (index)						"_si"
	Styles			StructItem `default: Name>_st`	//Файлы стилей (my.css)									"_st"
	Video			StructItem `default: Name>_vi`	//Файлы модуля «Видео» (video)							"_vi"
}

type StructItem struct {
	Name 	string
	Status 	bool
	Folder	string
	File	string
}

var Env *services.Env

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate ucoz files",
	Long:  `Validate ucoz files`,
	Run: func(cmd *cobra.Command, args []string) {
		config := Env.Config.GetStringMapString("ucoz")
		ucozPath := checkFolder(config["path"])
		createDefaultStruct()
		//checkBackup(ucozPath+"backup/")
		fmt.Printf("Inside subCmd Run with args: %v\n", ucozPath+"backup/")
	},
}

func checkFolder(path string) string{
	isExistFolder, err := afero.DirExists(Env.FileSystem, path)
	if err != nil{
		fmt.Errorf("Can't check config path. ")
	}
	if !isExistFolder{
		color.Yellow("Ucoz folder is exist. Creating...")
		Env.FileSystem.Mkdir(path, 0755)
	}
	return path
}

func checkBackup(path string) bool {
	isExistFolder, err := afero.DirExists(Env.FileSystem, path)
	if err != nil{
		fmt.Errorf("Can't check config path. ")
		return false
	}
	if !isExistFolder{
		color.Yellow("Backup folder is exist. Creating...")
		Env.FileSystem.Mkdir(path, 0755)
	}

	files, err := afero.ReadDir(Env.FileSystem, path)
	if err != nil {
		fmt.Errorf("Can't read backup dir. ")
		return false
	}
	for _, val := range files{
		if val.IsDir() {
			color.Red("Files: %v", val.Name())
		}
	}
	return true
}

func createDefaultStruct(){
	ucoz := new(ucozFileStruct)

	fmt.Println(ucoz)
}