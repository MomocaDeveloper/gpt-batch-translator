package utils
import (

	"os"
	"fmt"
)

func GetFileNames(typ string)[]string{
	var dirPath string
	if typ == "upload"{
		dirPath = "uploads"
	}else{
		dirPath = "downloads"
	}   
	fileNames := []string{}
	dir, err := os.Open(dirPath)
	if err != nil {
		fmt.Println("get file name fail", err)
		return nil 
	}   
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("get file name fail 2", err)
		return nil 
	}   

	for _, fileInfo := range fileInfos{
		if fileInfo.Mode().IsRegular(){
			fileNames = append(fileNames, fileInfo.Name())
		}   
	}   
	return fileNames
}

