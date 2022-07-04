//main.go
package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type HandlerWrapper struct {
	HandlerFunc  func(http.ResponseWriter, *http.Request)
	DB           *sql.DB
	HandlerErorr ApiError
}

type ApiError struct {
	StatusCode int
	Error      error
}

func main() {
	fmt.Println("hello")
	router := mux.NewRouter()
	//regex := `pic_[0-9]+\.jpg+`
	var handlers HandlerWrapper
	handlers.HandlerErorr = ApiError{http.StatusAccepted, nil}
	connStr := "host = db port = 5432 user=admin password=root dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Невозможно подключиться к базе данных")
		return
	}
	handlers.DB = db
	defer handlers.DB.Close()
	router.HandleFunc("/images/page/{page:[0-9]+}/pic/{pic_num:[0-9]+}", handlers.ServeDynamicPictures).Methods("GET")
	router.HandleFunc("/images/page/{page:[0-9]+}/pic/{pic_num:[0-9]+}/text", handlers.ServeDynamicPicturesText).Methods("GET")

	router.HandleFunc("/admin/upload", handlers.ServeAddPictures).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func (hdlr HandlerWrapper) ServeDynamicPicturesText(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------------------------------------------")
	vars := mux.Vars(r)
	page, _ := strconv.Atoi(vars["page"])
	pic_num, _ := strconv.Atoi(vars["pic_num"])
	pic_id := page*2 - 2 + pic_num
	rows, err := hdlr.DB.Query(`SELECT picture_name, picture_description, author, price, is_purchased FROM pictures WHERE picture_id = $1`, pic_id)
	if err != nil {
		http.Error(w, "cant select from database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	PR := PictureRecord{}
	for rows.Next() {
		err = rows.Scan(&PR.Picture_name, &PR.Picture_description, &PR.Author, &PR.Price, &PR.Is_purchased)
		if err != nil {
			http.Error(w, "cant scan data from database: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	fn := func(is_purchased bool) string {
		if is_purchased {
			return "продано"
		} else {
			return "на продаже"
		}
	}
	fmt.Fprint(w, `<div> <h1> Автор: `, PR.Author, `</h1> </div>`,
		`<div> Название: `, PR.Picture_name, `</div>`,
		`<div> Цена: `, PR.Price, ` руб. </div>`,
		`<div> Описание:`, PR.Picture_description, `</div>`,
		`<div> Состояние: `, fn(PR.Is_purchased), `</div>`)
	w.WriteHeader(http.StatusAccepted)
}

func (hdlr HandlerWrapper) ServeDynamicPictures(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----------------------------------------------")
	vars := mux.Vars(r)
	page, _ := strconv.Atoi(vars["page"])
	pic_num, _ := strconv.Atoi(vars["pic_num"])
	pic_id := page*2 - 2 + pic_num
	var path string
	rows, err := hdlr.DB.Query(`SELECT picture_path FROM pictures WHERE picture_id = $1`, pic_id)
	if err != nil {
		http.Error(w, "cant select from database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&path)
		if err != nil {
			http.Error(w, "cant scan data from database: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "cant read file from db: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpg")
	w.Header().Set("accept-ranges", "bytes")
	content_length := strconv.Itoa(len(fileBytes))
	w.Header().Set("content-length", content_length)
	w.Write(fileBytes)

	w.WriteHeader(http.StatusAccepted) //RETURN HTTP CODE 202
}

//добавить нормальную валидацию параметров и чтению в структуру в отдельную функцию/пакет
type PictureRecord struct {
	Picture_name        string
	Picture_description string
	Author              string
	Price               float32
	Is_purchased        bool
	Picture_path        string
}

func AddPictureRecordToDB(db *sql.DB, PR PictureRecord) (sql.Result, error) {
	fmt.Println(PR)
	result, err := db.Exec(`INSERT INTO pictures (picture_name, picture_description, author, price, is_purchased,
		picture_path) VALUES ($1, $2, $3, $4, $5, $6)`, PR.Picture_name, PR.Picture_description, PR.Author, PR.Price,
		PR.Is_purchased, PR.Picture_path)
	fmt.Println(result, "asdsad")
	fmt.Println(err, "aswwwww")
	return result, err
}

func ValidateParams(r *http.Request) (PictureRecord, error) {
	uploadData, handler, err := r.FormFile("my_file")
	if err != nil {
		return PictureRecord{}, errors.Wrap(err, "cant read file")
	}
	defer uploadData.Close()
	picture_path := "images/" + handler.Filename
	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		return PictureRecord{}, errors.Wrap(err, "cant parse (incorrect) price value")
	}
	is_purchased, err := strconv.ParseBool(r.FormValue("is_purchased"))
	if err != nil {
		return PictureRecord{}, errors.Wrap(err, "cant parse (incorrect) purchased value")
	}
	PR := PictureRecord{r.FormValue("picture_name"), r.FormValue("picture_description"), r.FormValue("author"),
		float32(price), is_purchased, picture_path}
	newFile, err := os.Create(picture_path)
	if _, err := io.Copy(newFile, uploadData); err != nil {
		return PictureRecord{}, errors.Wrap(err, "cant save file")
	}
	//newFile.Sync()
	newFile.Close()
	return PR, err
}

func (hdlr HandlerWrapper) ServeAddPictures(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	PR, err := ValidateParams(r)
	if err != nil {
		http.Error(w, "Error in Validating Params"+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("here2")
	res, err := AddPictureRecordToDB(hdlr.DB, PR)
	if err != nil {
		http.Error(w, "cant insert file in database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("here3")
	fmt.Fprintf(w, "Success upload to database %v\n", res)
	w.WriteHeader(http.StatusAccepted)
}
