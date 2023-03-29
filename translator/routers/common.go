package routers

import (
	"os"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context){
	fmt.Println("ping pong")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Download(c *gin.Context){
	fileName := c.Query("file")
	var filePath string
	typ := c.Query("type")
	if typ == "upload"{
		filePath = "uploads"
	}else{
		filePath = "downloads"
	}
	c.File(fmt.Sprintf("%s/%s", filePath, fileName))
}

func UploadFile(c *gin.Context){
	file, err := c.FormFile("file")
	ftyp := c.PostForm("type")
	filePath := "uploads/"
	if ftyp == "2"{
		filePath = "keywords/"
	}else if ftyp == "3"{
		filePath = "characters/"
	}
	if err != nil{
		fmt.Println("err1", err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	err = c.SaveUploadedFile(file, fmt.Sprintf("%s%s", filePath, file.Filename))
	if err != nil{
		fmt.Println("err2", err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, "File %s uploaded successfully!", file.Filename)
}

func DeleteFile(c *gin.Context){
	filename := c.Query("file")
	err := os.Remove(fmt.Sprintf("downloads/%s", filename))
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}

func BrowseFile(c *gin.Context){
	
}

