package jobs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

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

func fixAndSaveFileName(files *[]string,save bool) {
	for i, filePath := range *files {
		if strings.HasSuffix(filePath,".mp3") {
			(*files)[i] = exposeMusicName(filePath,(*files)[0])
			if save {
				err := os.Rename(filePath, (*files)[0]+ "/" + (*files)[i])
				if err != nil {
					log.Println("File rename failed:\n", filePath+ "\to: "+ (*files)[i])
				}
			}
			fmt.Println((*files)[i])
		} else {
			// do not puts ignored musics root folder
			if  !strings.HasSuffix(filePath,"/musics") {
				log.Println("Ignored file: ", filePath)
			}
		}
	}
}

func clearIfSepWithLine ( name *string ) {
	// if ( - ) count more than 1, words can be seperated with -
	// so replace them with white spaces
	if strings.Count( *name, "-") > 1 {
		*name = strings.Replace(*name, "-", " ", -1)
	}
}

func clearBrackets( name string ) string {
	var fixedName string
	nameLength := len(name)
	for i := 0; i < nameLength; i++ {
		if name[i] == '[' || name[i] == '{' {
			expect := name[i] + 2 // checked from ascii table
			k := i
			for ; k < nameLength && name[k] != expect ; k++ {}
			if name[k] == expect {
				i = k
			}
		} else {
			fixedName += string(name[i])
		}
	}
	return fixedName
}

func trimRootPath( name,root *string ) {
	*name = strings.Replace(*name,*root+"/","",1)
}

func exposeMusicName (name, root string) string {
	var fixedName string
	trimRootPath(&name, &root)
	clearIfSepWithLine(&name)
	name = clearBrackets(name)
	for _, r := range name {
		if !unicode.IsLetter(r) {
			switch r {
			case '_', '^' , '+' , ';' , ':' , '"', ',' :
				fixedName += " "
			default:
				fixedName += string(r)
			}
		} else {
			fixedName += string(r)
		}
	}
	return fixedName
}

func FixFileNames (folderPath string,save bool) (string, error ){

	files , path , err := getFileNames(folderPath); if err != nil {
		log.Println("Couldn't get music file names: ",err)
		return path, err
	}
	fixAndSaveFileName(&files, save)
	return path, nil
}

