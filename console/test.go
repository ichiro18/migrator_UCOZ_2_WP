package console

import (
	"strings"

	"fmt"

	"path"

	"time"

	"github.com/ichiro18/migrator_UCOZ_2_WP/console/ucoz"
	"github.com/schollz/progressbar"
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
		count := len(news)
		bar := progressbar.New(count)
		for _, val := range news {
			bar.Add(1)
			updateMediaPath(&val)
			time.Sleep(100 * time.Millisecond)
		}
	},
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
	var filePath string
	filePath = path.Join(ucozPath["path"], "backup", mediaPath)

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
	if err != nil {
		panic(err)
	}
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
		Env.FileSystem.Mkdir(path, 0755)
	}
	return path
}
func init() {
	RootCmd.AddCommand(testCmd)
}
