package console

import (
	"strings"

	"fmt"

	"path"

	"github.com/fatih/color"
	"github.com/ichiro18/migrator_UCOZ_2_WP/console/ucoz"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "This command for test",
	Long:  `This command for test`,
	Run: func(cmd *cobra.Command, args []string) {
		str := ucoz.Site{}
		// Get Ucoz News
		news := ucoz.GetNews(&str)
		item := news[0]
		item.MESSAGE = updateMediaPath(item)
	},
}

func updateMediaPath(data ucoz.NewItem) string {
	item := data.MESSAGE
	srcIndex := strings.Index(item, `src="`) + 5
	var src string
	var srcIndexStart int
	var srcIndexEnd int
	var res string
	if srcIndex != -1 {
		srcIndexEnd = strings.Index(item[srcIndex:], string('"'))
		srcPart := item[srcIndex:]
		srcIndexStart = srcIndex
		src = srcPart[:srcIndexEnd]
	}
	if src != "" {
		newSrc := copyMedia(src, &data)

		// Replace String
		before := item[:srcIndexStart]
		strStart := item[srcIndexStart:]
		after := strStart[srcIndexEnd:]
		res = before + newSrc + after
	}
	if res == "" {
		fmt.Errorf("Can't update media path. ")
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
	postFolder := checkFolder(path.Join(yearFolder, "posts"))

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
func checkFolder(path string) string {
	isExistFolder, err := afero.DirExists(Env.FileSystem, path)
	if err != nil {
		fmt.Errorf("Can't check config path. ")
	}
	if !isExistFolder {
		color.Yellow("WP folder is exist. Creating...")
		Env.FileSystem.Mkdir(path, 0755)
	}
	return path
}
func init() {
	RootCmd.AddCommand(testCmd)
}
