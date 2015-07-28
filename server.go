package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	reader, err := r.MultipartReader()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := JsonResponse{make(map[string]string)}
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, part); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Spectrograms[part.FileName()] = r.Host + "/png/" + part.FileName() + ".png"
	}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	BuildSpectrograms()
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

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

	v1 := r.Group("api/v1")
	{
		v1.POST("/post", gin.WrapF(UploadHandler))
		v1.GET("/post", func(c *gin.Context) {
			c.JSON(http.StatusMethodNotAllowed,
				gin.H{"message": "Please, POST a file",
					"status": http.StatusMethodNotAllowed})
		})
	}
	r.Static("/png", "./processed/")

	r.Run(":8080")
}
