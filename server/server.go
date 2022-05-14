package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	engine   *gin.Engine
	logger   *Logger
	addr     string
	basePath string
}

func (s *Server) Run() {
	err := s.engine.Run(s.addr)
	if err != nil {
		s.logger.Error("Error starting cloud file server.")
		return
	}
}

func NewServer(c Config) *Server {
	logger := NewLogger()
	router := gin.New()
	server := &Server{
		engine:   router,
		logger:   logger,
		addr:     c.addr,
		basePath: c.basePath,
	}
	router.GET("/", server.browse)
	router.POST("/upload", server.upload)
	router.POST("/delete/:filename", server.delete)
	router.GET("/download/:filename", server.download)
	return server
}

func link(name string) string {
	return fmt.Sprintf("<a download=\"%s\" href=\"/download/%s\">%s</a><form action=\"/delete/%s\" method=\"post\" multiple enctype=\"multipart/form-data\"><input type=\"submit\" value=\"delete\" id=\"file-del\"></form>\n\n", name, name, name, name)
}

func (s *Server) browse(c *gin.Context) {
	files, err := ioutil.ReadDir(s.basePath)
	reason := ""
	if err != nil {
		reason = fmt.Sprintf("Error checking directory.")
		s.logger.Error(reason)
		c.String(http.StatusBadRequest, reason)
		return
	}
	template := "<ul>\n%s</ul>\n"
	ret := ""

	for _, file := range files {
		if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			ret += "<li>" + link(file.Name()) + "</li>\n"
		}
	}
	ret = fmt.Sprintf(template, ret)
	ret += "<form action=\"/upload\" method=\"post\" multiple enctype=\"multipart/form-data\" accept-charset=\"UTF-8\">\n    <div><input type=\"file\" name=\"file\" id=\"file\"></div>\n    <div><input type=\"submit\" value=\"upload\" id=\"fileUpload\"></div>\n</form>"
	c.Data(http.StatusOK, "text/html;charset=utf-8", []byte(ret))
}

func (s *Server) upload(c *gin.Context) {
	s.logger.Log("attempting upload!")
	reason := ""
	file, err := c.FormFile("file")
	if err != nil {
		reason = fmt.Sprintf("Error receiving file. %s", err.Error())
		s.logger.Error(reason)
		c.String(http.StatusBadRequest, reason)
		return
	}
	if strings.HasPrefix(file.Filename, ".") {
		reason = "Cannot upload file with prefix '.'"
		s.logger.Error(reason)
		c.String(http.StatusBadRequest, reason)
		return
	}
	fmt.Println("This is file name : ", file.Filename)
	err = c.SaveUploadedFile(file, s.basePath+file.Filename)
	if err != nil {
		reason = fmt.Sprintf("Error uploading file [%s]. %s", file.Filename, err.Error())
		s.logger.Error(reason)
		c.String(http.StatusBadRequest, reason)
		return
	}
	reason = fmt.Sprintf("File [%s] uploaded!", file.Filename)
	s.logger.Log(reason)
	c.String(http.StatusOK, reason)
}

func (s *Server) download(c *gin.Context) {
	fmt.Println("attempting download!")
	name := c.Param("filename")
	reason := ""
	if strings.HasPrefix(name, ".") {
		reason = "Cannot download file with prefix '.'"
		s.logger.Error(reason)
		c.String(http.StatusBadRequest, reason)
		return
	}
	files, err := ioutil.ReadDir(s.basePath)
	if err != nil {
		reason = fmt.Sprintf("Error checking directory.")
		s.logger.Error(reason)
		c.String(http.StatusBadRequest, reason)
		return
	}
	flag := false
	for _, file := range files {
		if file.Name() == name && !file.IsDir() {
			flag = true
		}
	}
	if flag == false {
		reason = fmt.Sprintf("Cannot find file.")
		s.logger.Error(reason)
		c.String(http.StatusNotFound, reason)
		return
	}
	fmt.Println("transferring!", name)
	c.File(s.basePath + name)
}

func (s *Server) delete(c *gin.Context) {
	name := c.Param("filename")
	reason := ""
	if strings.HasPrefix(name, ".") {
		reason = "Cannot delete file with prefix '.'"
		s.logger.Error(reason)
		c.String(http.StatusBadRequest, reason)
		return
	}
	files, err := ioutil.ReadDir(s.basePath)
	if err != nil {
		reason = fmt.Sprintf("Error checking directory.")
		s.logger.Error(reason)
		c.String(http.StatusBadRequest, reason)
		return
	}
	flag := false
	for _, file := range files {
		if file.Name() == name && !file.IsDir() {
			flag = true
		}
	}
	if flag == false {
		reason = fmt.Sprintf("Cannot find file.")
		s.logger.Error(reason)
		c.String(http.StatusNotFound, reason)
		return
	}
	err = os.Remove(s.basePath + name)
	if err != nil {
		reason = fmt.Sprintf("Cannot delete file.")
		s.logger.Error(reason)
		c.String(http.StatusBadRequest, reason)
		return
	}
	reason = fmt.Sprintf("File %s is deleted.", name)
	c.String(http.StatusOK, reason)
}
