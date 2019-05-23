package diffrunner

import (
	"fmt"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

func getTestCaseSqlFile(dir string) []string {
	return loadItemsFromDir(dir, false)
}

func getAllTestCaseDir(dir string) []string {
	return loadItemsFromDir(dir, true)
}

// load all items from a directory, sub directory or file base on the loadDirectory parameter
func loadItemsFromDir(dir string, loadDirectory bool) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error("read dir failed", zap.String("directory", dir), zap.String("err", fmt.Sprintf("%v", err)))
		os.Exit(-1)
	}
	var filesPaths []string
	for _, f := range files {
		if !loadDirectory && !f.IsDir() {
			if strings.HasSuffix(f.Name(), ".sql") && f.Name() != "dml.sql" {
				filesPaths = append(filesPaths, path.Join(dir, f.Name()))
			} else {
				log.Info("ignore dml.sql file")
			}
		} else if loadDirectory && f.IsDir() {
			filesPaths = append(filesPaths, path.Join(dir, f.Name()))
		} else {
			log.Info("ignore items", zap.String("name", f.Name()))
		}
	}
	//sort it
	sort.Strings(filesPaths)
	return filesPaths
}
