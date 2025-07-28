package route

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fawrwebservice/model"
	"gorm.io/gorm"
)

const imagesFolder = "/root/data/images"

type CommentRoute struct {
	DB *gorm.DB
}

func (route *CommentRoute) Register(parent *http.ServeMux) {
	parent.HandleFunc("/comment", route.commentHandler)
	parent.HandleFunc("/comment/img", route.imgHandler)
}

func (route *CommentRoute) commentHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		route.get(w, req)
	case http.MethodPost:
		route.add(w, req)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (route *CommentRoute) add(w http.ResponseWriter, req *http.Request) {
	in, header, err := req.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer in.Close()

	comment := model.Comment{
		Bought: false,
	}
	err = route.DB.Create(&comment).Error
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::add] DB create error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll(imagesFolder, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::add] Create dir error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	comment.ImgPath = filepath.Join(imagesFolder, strconv.Itoa(comment.ID), filepath.Ext(header.Filename))
	out, err := os.Create(comment.ImgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::add] Out file error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer out.Close()

	io.Copy(out, in)

	err = route.DB.Updates(&comment).Error
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::add] DB update error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(&comment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::add] JSON error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (route *CommentRoute) get(w http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("id") == "" {
		route.getAll(w, req)
	} else {
		route.getSingle(w, req)
	}
}

func (route *CommentRoute) getAll(w http.ResponseWriter, _ *http.Request) {
	var comments []model.Comment
	err := route.DB.Find(&comments).Error
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::getAll] DB error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(&comments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::getAll] JSON error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (route *CommentRoute) getSingle(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	var comment model.Comment
	err := route.DB.Model(&model.Comment{}).Where("id = ?", id).First(&comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprintf(os.Stderr, "[CommentRoute::getSingle] DB error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(&comment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::getSingle] JSON error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (route *CommentRoute) imgHandler(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	var comment model.Comment
	err := route.DB.Model(&model.Comment{}).Where("id = ?", id).First(&comment).Error
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::imgHandler] DB error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file, err := os.Open(comment.ImgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[CommentRoute::imgHandler] File open error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	http.ServeContent(w, req, file.Name(), time.Now(), file)
}
