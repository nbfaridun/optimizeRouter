package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	_ "github.com/tiny-go/codec/driver"
	_ "github.com/tiny-go/middleware"
	"log"
	"net/http"
	"strconv"
	"time"
)

type DiaryEntry struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type statusCodeWriter struct {
	http.ResponseWriter
	statusCode int
}

var diaryEntries []DiaryEntry

func Log(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		statusCodeWriter := &statusCodeWriter{ResponseWriter: w}
		next(w, r, ps)
		statusCode := statusCodeWriter.statusCode
		log.Printf("%s %s %s %d", r.RemoteAddr, r.Method, r.URL, statusCode)
		//todo add cors middleware
	}
}

func createDiaryEntry(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var diaryEntry DiaryEntry
	err := json.NewDecoder(request.Body).Decode(&diaryEntry)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	diaryEntry.ID = len(diaryEntries) + 1
	diaryEntry.CreatedAt = time.Now()
	diaryEntry.UpdatedAt = time.Now()
	diaryEntries = append(diaryEntries, diaryEntry)
	json.NewEncoder(writer).Encode(diaryEntry)
}

func getAllDiaryEntries(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	json.NewEncoder(writer).Encode(diaryEntries)
}

func getDiaryEntry(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	for _, d := range diaryEntries {
		if d.ID == id {
			json.NewEncoder(writer).Encode(d)
			return
		}
	}
	http.NotFound(writer, request)
}

func updateDiaryEntry(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	var diaryEntry DiaryEntry
	err = json.NewDecoder(request.Body).Decode(&diaryEntry)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	for i, d := range diaryEntries {
		if d.ID == id {
			diaryEntry.ID = id
			diaryEntry.CreatedAt = d.CreatedAt
			diaryEntry.UpdatedAt = time.Now()
			diaryEntries[i] = diaryEntry
			json.NewEncoder(writer).Encode(diaryEntry)
			return
		}
	}
	http.NotFound(writer, request)
}

func deleteDiaryEntry(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	for i, d := range diaryEntries {
		if d.ID == id {
			diaryEntries = append(diaryEntries[:i], diaryEntries[i+1:]...)
			writer.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(writer, request)
}

func main() {
	router := httprouter.New()

	router.GET("/diary_entries", Log(getAllDiaryEntries))
	router.GET("/diary_entries/:id", Log(getDiaryEntry))
	router.POST("/diary_entries", Log(createDiaryEntry))
	router.PUT("/diary_entries/:id", Log(updateDiaryEntry))
	router.DELETE("/diary_entries/:id", Log(deleteDiaryEntry))

	log.Fatal(http.ListenAndServe(":8080", router))
}

/****/
