package jobs

import (
	"log"
	"os"
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/errors"
)

func getRootPath () (string, error){
	// get current path
	root, err := os.Getwd()
	return root, err
}

func getFileNames (folderPath string ) ( []string, string, error ){
	var files []string
	root := folderPath
	// set the root path if path parameter empty
	if folderPath == "" {
		path, err := getRootPath(); if err != nil {
			return nil,path,err
		}
		root = path + "/musics"
	}
	// read file names
	err  	  := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files  = append(files, path)
		return err
	}); if   err !=  nil {
		return nil,root, err
	}
	if len(files) <= 1 {
		return nil,root, errors.New(".mp3 files couldnt't found at " + root)
	}
	return files,root,nil
}

func fixAndChangeFileNames(files *[]string) {
	for _, filePath := range *files {
		if strings.HasSuffix(filePath,".mp3") {
			fmt.Println(filePath)
		} else {
			// do not puts ignored musics root folder
			if  !strings.HasSuffix(filePath,"/musics") {
				log.Println("Ignored file: ", filePath)
			}
		}
	}
}

func FixFileNames (folderPath string,save bool) (string, error ){

	files , path , err := getFileNames(folderPath); if err != nil {
		log.Println("Couldn't get music file names: ",err)
		return path, err
	}
	if save == true {
		fixAndChangeFileNames(&files)
	}

	return path, nil
}

