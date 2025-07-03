package controllers

import (
	"github.com/gin-gonic/gin"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"hotel-management-system/utils"
	"net/http"
	"strconv"
	"time"
)

func GetImg(c *gin.Context) {
	filename := c.Param("filename")
	c.File("./uploads/" + filename)
}

type ImgResponse struct {
	Id  uint   `json:"id"`
	Url string `json:"url"` // 图片的访问 URL
}

func UploadImg(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "上传文件失败"})
		return
	}

	// 检查文件类型
	if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "只允许上传 JPEG 或 PNG 格式的图片"})
		return
	}

	// 保存文件到服务器，文件名前加当天的秒数
	now := time.Now()
	secondsOfDay := now.Hour()*3600 + now.Minute()*60 + now.Second()
	filePath := "uploads/" + strconv.Itoa(secondsOfDay) + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存文件失败"})
		return
	}

	// 将文件路径存储到数据库
	img := models.Img{
		Url: filePath,
	}
	if err := global.Db.Create(&img).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "保存图片信息失败"})
		return
	}

	// 将图片对象加入上传队列
	if utils.ImgUploadChan != nil {
		utils.ImgUploadChan <- &img
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "图片上传成功",
		"img": ImgResponse{
			Id:  img.ID,
			Url: img.Url,
		},
	})
}
