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

const MAX_RETRY_TIME = 20
const MAX_GOROUTINE = 20

func startTranslate(input, keyword, character, temperature string, progress *Progress)error {
	result, readErr := ParseFullFile(input)	
	if readErr != nil{
		//fmt.Println(readErr)
		return readErr	
	}
	totalLine := result.CalcTotalLine()
	progress.CreateNewPool(input, MAX_GOROUTINE, totalLine)
	conf := initialization.LoadConfig("./config.yaml")
	gpt := NewChatGPT(*conf)
	keywordDict, _ := GetBaseNoun(keyword)
	characterDict, _ := GetBaseCharacter(character)
	workerPool, _ := progress.GetWorkerPool(input)
	for _, content := range result.Contents{
		workerPool.AddTask(func(con *Content){
			for i:=0;i<=MAX_RETRY_TIME;i++{
				line, transferErr := con.TranslateMultiLines(gpt, &keywordDict, &characterDict, temperature)
				if transferErr != nil{
					fmt.Println(transferErr)
				}else{
					progress.UpdateProgress(input, line)
					break
				}
				if i==MAX_RETRY_TIME{
					fmt.Println("translate error, retry times max, line", con.Id)
					progress.UpdateProgress(input, len(con.Cells))

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
		temperature := c.Query("temperature")

		fmt.Println("start translate file", file, characterFile, temperature, keywordFile)
		if _, exist := progress.GetWorkerPool(file);exist{
			c.JSON(http.StatusTooManyRequests, struct{}{})
			return
		}
		go startTranslate(file, keywordFile, characterFile, temperature, progress)
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
	r.GET("/conversation", func(c *gin.Context){
		conf := initialization.LoadConfig("../config.yaml")
		gpt := NewChatGPT(*conf)
		feifeisheding := "请扮演菲菲，和我讨论一只从天而降的走路菇的故事。 菲菲是一只有礼貌的、普普通通的小猫娘，生活在异世界。 你和我经营着一家普普通通的旅店，在一天晚上，这只从天而降的走路菇砸坏了我们的旅店。 我们之后的对话将会涵盖走路菇的故事、旅店的修复和我们之间的交流和互动。请注意：你现在是菲菲，之后说话的口吻都要和菲菲的身份匹配。"
		beiersheding := "请扮演贝尔，和我讨论一只从天而降的走路菇的故事。 贝尔是一个傲娇、容易着急、直率、坦诚、直截了当的乡村人。 你和我经营着一家普普通通的旅店，在一天晚上，这只从天而降的走路菇砸坏了我们的旅店。 我们之后的对话将会涵盖走路菇的故事、旅店的修复和我们之间的交流和互动。请注意：你现在是贝尔，之后说话的口吻都要和贝尔的身份匹配。"
		feifei := NewCharacter("feifei", feifeisheding, gpt)
		beier := NewCharacter("beier", beiersheding, gpt)
		
		go func(){
			res := SimulateConversation(feifei, beier, "菲菲，不好了不好了，我们的旅店...", 50, 0.7)
			fmt.Println("simulation success", res)
		}()
		c.JSON(http.StatusOK, struct{}{})
	})
	r.Run(":9003")
}

