package main

import (
	"github.com/julienschmidt/httprouter"
	_ "github.com/tiny-go/codec/driver"
	_ "github.com/tiny-go/middleware"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()

	router.GET("/diary_entries", Log(getAllDiaryEntries))
	router.GET("/diary_entries/:id", Log(getDiaryEntry))
	router.POST("/diary_entries", Log(createDiaryEntry))
	router.PUT("/diary_entries/:id", Log(updateDiaryEntry))
	router.DELETE("/diary_entries/:id", Log(deleteDiaryEntry))

	log.Fatal(http.ListenAndServe(":8080", router))
}
