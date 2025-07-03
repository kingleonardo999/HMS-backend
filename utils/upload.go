package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-management-system/global"
	"hotel-management-system/models"
	"net/http"
	"path/filepath"
)

// ImgUploadChan 图片上传队列
var ImgUploadChan chan *models.Img

var picGoUrl = "http://localhost:36677/upload"

// Upload2cloud 启动后台阻塞上传协程
func Upload2cloud() {
	ImgUploadChan = make(chan *models.Img, 100)
	go func() {
		for img := range ImgUploadChan {
			// 将相对路径转为绝对路径
			localPath, err := filepath.Abs(img.Url)
			if err != nil {
				fmt.Println("获取绝对路径失败:", err)
				continue
			}

			// 用 JSON 格式上传图片路径
			payload := map[string]interface{}{
				"list": []string{localPath},
			}
			data, err := json.Marshal(payload)
			if err != nil {
				fmt.Println("序列化JSON失败:", err)
				continue
			}

			req, err := http.NewRequest("POST", picGoUrl, bytes.NewReader(data))
			if err != nil {
				fmt.Println("创建请求失败:", err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("上传图片失败:", err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Println("上传图片返回非200:", resp.Status)
				continue
			}

			var result struct {
				Success bool     `json:"success"`
				Result  []string `json:"result"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				fmt.Println("解析响应失败:", err)
				continue
			}
			if !result.Success || len(result.Result) == 0 {
				fmt.Println("上传失败或未返回url")
				continue
			}

			// 更新数据库中的url
			if img.ID != 0 {
				global.Db.Model(&models.Img{}).Where("id = ?", img.ID).Update("url", result.Result[0])
			} else {
				fmt.Println("图片ID无效，无法更新url")
			}
		}
	}()
}
