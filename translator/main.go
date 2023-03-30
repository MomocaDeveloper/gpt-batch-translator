package main

import (
		//"time"
		"net/http"
		"strings"
		"fmt"
		"github.com/gin-gonic/gin"

		. "translator/services"
		router "translator/routers"
		util "translator/utils"
		"translator/initialization"
	)

var MAX_RETRY_TIME int = 20

func startTranslate(input, keyword, character string, progress *Progress)error {
	result, readErr := ParseFullFile(input)	
	if readErr != nil{
		//fmt.Println(readErr)
		return readErr	
	}
	totalLine := result.CalcTotalLine()
	progress.CreateNewPool(input, 10, totalLine)
	conf := initialization.LoadConfig("./config.yaml")
	gpt := NewChatGPT(*conf)
	keywordDict, _ := GetBaseNoun(keyword)
	characterDict, _ := GetBaseCharacter(character)
	workerPool, _ := progress.GetWorkerPool(input)
	for _, content := range result.Contents{
		workerPool.AddTask(func(con *Content){
			for i:=0;i<=MAX_RETRY_TIME;i++{
				line, transferErr := con.TranslateMultiLines(gpt, &keywordDict, &characterDict)
				if transferErr != nil{
					fmt.Println(transferErr)
				}else{
					progress.UpdateProgress(input, line)
					break
				}
			}
		}, content)
	}

	workerPool.Wait()
	if progress.AlreadyFinishPool(input){
		progress.FinishPool(input)
	}else{
		fmt.Println("error: not finishi all translate, please check the output")
	}

	writeErr := WriteIntoFile(*result)
	if writeErr != nil {
		fmt.Println(writeErr)
		return writeErr
	}
	return nil
	//fmt.Println("write success", output)
}

func main(){
	//input := "test1.xlsx"
	//output := "result1.xlsx"
	//_ = startTranslate(input, output, basic)
	allPools := make(map[string]*WorkerPool)
	progress := &Progress{allPools,}
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context){
		data := gin.H{
			"removeSuffix": func(str string)string{
				return strings.TrimSuffix(str, "-bar")
			},
			"Progress": progress.CreateMsgData(),
			"UploadFile": util.GetFileNames("upload"), 
			"DownloadFile": util.GetFileNames("download"),
		}
		c.HTML(http.StatusOK, "index.html", data)
	})
	r.GET("/ping", router.Ping)
	r.POST("/upload", router.UploadFile)
	r.GET("/download", router.Download)
	r.GET("/delete", router.DeleteFile)
	r.GET("/translate", func(c *gin.Context){
		file := c.Query("file")
		//keywordFile := c.Query("keyword")
		//characterFile := c.Query("character")
		characterFile := "character.xlsx"
		keywordFile := "basic.xlsx"
		fmt.Println("start translate file", file, characterFile, keywordFile)
		if _, exist := progress.GetWorkerPool(file);exist{
			c.JSON(http.StatusTooManyRequests, struct{}{})
			return
		}
		go startTranslate(file, keywordFile, characterFile, progress)
		c.JSON(http.StatusOK, struct{}{})
	})
	r.GET("/progress", func(c *gin.Context){
		c.JSON(http.StatusOK, progress.CreateMsgData())
	})
	r.GET("/browse", func(c *gin.Context){
		rawFileName := c.Query("file")
		file := fmt.Sprintf("downloads/%s", rawFileName)
		fileData, err := ParseFullFile(file)
		if err != nil {
			fmt.Println("check file error", err)
		}
		c.HTML(http.StatusOK, "browse.html", gin.H{
			"data": fileData,
		})
	})
	r.Run(":9003")
}

