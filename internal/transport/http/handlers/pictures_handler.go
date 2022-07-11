package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"web_app/internal/domain"
	"web_app/internal/service"

	//sender "web_app/pkg/tempaltes_sender"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

type Handler struct {
	services *service.Services
}

func NewHandler(srv *service.Services) *Handler {
	return &Handler{
		services: srv,
	}
}

func (h *Handler) InitHandlers(router *mux.Router) {
	router.HandleFunc("/page/{page:[0-9]+}/pic/{pic:[0-9]+}", h.ServeDynamicPictures).Methods("GET")
	router.HandleFunc("/page/{page:[0-9]+}/pic/{pic:[0-9]+}/text", h.ServeDynamicPicturesText).Methods("GET")
	router.HandleFunc("/admin/upload", h.ServeAddPictures).Methods("POST")
}

func (h *Handler) ServeDynamicPictures(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pic, _ := strconv.Atoi(vars["pic"])
	page, _ := strconv.Atoi(vars["page"])
	pic_id := page*2 - 2 + pic
	path, err := h.services.Pictures.GetPicturePathById(uint(pic_id))
	if err != nil {
		http.Error(w, "cant get picture from db: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "cant read file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpg")
	w.Header().Set("accept-ranges", "bytes")
	content_length := strconv.Itoa(len(fileBytes))
	w.Header().Set("content-length", content_length)
	w.Write(fileBytes)
	w.WriteHeader(http.StatusAccepted)

}

func (h *Handler) ServeDynamicPicturesText(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page, _ := strconv.Atoi(vars["page"])
	pic, _ := strconv.Atoi(vars["pic"])
	pic_id := page*2 - 2 + pic
	picture, err := h.services.Pictures.GetPictureById(uint(pic_id))
	if err != nil {
		http.Error(w, "cant get picture from db: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, h.services.TemplateParser.GenerateBodyFromHTML(picture))
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) ServeAddPictures(w http.ResponseWriter, r *http.Request) {
	res, err := h.ValidateParams(r)
	if err != nil {
		http.Error(w, "cant validate params: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.services.Pictures.CreatePictureRecord(res)
	if err != nil {
		http.Error(w, "cant add pic to db: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)

}

func (h *Handler) ValidateParams(r *http.Request) (domain.Picture, error) {
	uploadData, handler, err := r.FormFile("my_file")
	if err != nil {
		return domain.Picture{}, errors.Wrap(err, "cant read file")
	}
	fileExtension := filepath.Ext(handler.Filename)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	picture := &domain.Picture{}
	err_ := decoder.Decode(picture, r.MultipartForm.Value)
	if err_ != nil {
		return *picture, errors.Wrap(err, "cant decode msg")
	}
	defer uploadData.Close()
	newFile, file_path, err := h.services.StorageManager.CreateUnicFile(fileExtension)
	picture.Picture_path = file_path
	if err != nil {
		return domain.Picture{}, errors.Wrap(err, "cant create file")
	}
	defer newFile.Close()
	if _, err := io.Copy(newFile, uploadData); err != nil {
		return domain.Picture{}, errors.Wrap(err, "cant copy file from request file")
	}
	return *picture, nil
}
