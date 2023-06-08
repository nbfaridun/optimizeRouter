package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

var diaryEntries []DiaryEntry

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
