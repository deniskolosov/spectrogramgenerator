package main

import (
	_ "encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func UploadHandler(c *gin.Context) {

	reader, err := c.Request.MultipartReader()

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := JsonResponse{make(map[string][]string)}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() == "" {
			continue
		}
		curDir, err := os.Getwd()
		PrintError(err)

		dst, err := os.Create(curDir + "/upload/" + part.FileName())
		defer dst.Close()
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Spectrograms["files"] = append(resp.Spectrograms["files"], "http://"+c.Request.Host+"/png/"+part.FileName()+".png")
	}
	//js, err := json.Marshal(resp)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	BuildSpectrograms()
	c.JSON(200, resp)

}
func init() {
	if _, err := os.Stat("./processed"); os.IsNotExist(err) {
		os.Mkdir("processed", 0777)

	}
	if _, err := os.Stat("./upload"); os.IsNotExist(err) {
		os.Mkdir("upload", 0777)

	}

}
func main() {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		// add header Access-Control-Allow-Origin
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return

		}

		c.Next()

	})

	v1 := r.Group("api/v1")
	{
		v1.POST("/post", UploadHandler)
		v1.GET("/post", func(c *gin.Context) {
			c.JSON(http.StatusMethodNotAllowed,
				gin.H{"message": "Please, POST a file",
					"status": http.StatusMethodNotAllowed})
		})
	}
	r.Static("/png", "./processed/")

	r.Run(":8080")
}
