package handler

import (
	"encoding/json"
	"fmt"
	"github.com/cristhoperdev/events-import/connection"
	"github.com/cristhoperdev/events-import/model"
	"github.com/cristhoperdev/events-import/utils"
	"github.com/labstack/echo"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var store connection.Datastore

func MoviePost(c echo.Context) error {
	utils.ConsoleLog(c)
	var jsonObj model.JsonResult
	var status int

	l := &lumberjack.Logger{
		Filename:   "logs/access.json",
		MaxSize:    0,
		MaxAge:     1,
		MaxBackups: 1,
		LocalTime:  false,
		Compress:   false,
	}
	log.SetOutput(l)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	jsonr, _ := json.MarshalIndent(event, "", "\t")
	log.Println(jsonr)
	go func() {
		for {
			<-ch
			l.Rotate()
		}
	}()


	decoder := json.NewDecoder(c.Request().Body)
	var movie *model.Movie
	if err := decoder.Decode(&movie); err != nil {
		jsonObj.Status = http.StatusBadRequest
		jsonObj.Data = "Body not well formed"
		status = http.StatusBadRequest
		return c.JSON(status, jsonObj)
	}

	if movie.Key == nil {
		key, _ := utils.NewUUID()
		movie.Key = &[]string{key}[0]
		movie.CreatedAt = utils.NowUTC()
		movie.UpdatedAt = utils.NowUTC()
	}

	err := store.Open()
	dbName := "event_dev"
	db, err := store.GetDatabase(dbName)
	if err != nil {
		fmt.Println(err)
	}

	//create movie
	err = db.CreateMovie(movie)

	if err != nil {
		jsonObj.Status = http.StatusBadRequest
		jsonObj.Data = err
		status = http.StatusBadRequest
	} else {
		jsonObj.Status = http.StatusOK
		jsonObj.Data = movie
		status = http.StatusOK
	}

	return c.JSON(status, jsonObj)
}

func MoviePut(c echo.Context) error {
	utils.ConsoleLog(c)
	var jsonObj model.JsonResult
	var status int

	decoder := json.NewDecoder(c.Request().Body)
	var movie *model.Movie
	if err := decoder.Decode(&movie); err != nil {
		jsonObj.Status = http.StatusBadRequest
		jsonObj.Data = "Body not well formed"
		status = http.StatusBadRequest
		return c.JSON(status, jsonObj)
	}

	movie.UpdatedAt = utils.NowUTC()

	err := store.Open()
	dbName := "event_dev"
	db, err := store.GetDatabase(dbName)
	if err != nil {
		fmt.Println(err)
	}

	//Update Movie
	err = db.UpdateMovie(movie)

	if err != nil {
		jsonObj.Status = http.StatusBadRequest
		jsonObj.Data = err
		status = http.StatusBadRequest
	} else {
		jsonObj.Status = http.StatusOK
		jsonObj.Data = movie
		status = http.StatusOK
	}

	return c.JSON(status, jsonObj)
}

func MovieDelete(c echo.Context) error {
	//utils.ConsoleLog(c)
	var jsonObj model.JsonResult
	var status int

	decoder := json.NewDecoder(c.Request().Body)
	var movie *model.Movie
	_ = decoder.Decode(&movie)


	err := store.Open()
	dbName := "event_dev"
	db, err := store.GetDatabase(dbName)
	if err != nil {
		fmt.Println(err)
	}

	//delete movie
	err = db.DeleteMovie(movie)

	if err != nil {
		jsonObj.Status = http.StatusBadRequest
		jsonObj.Data = err
		status = http.StatusBadRequest
	} else {
		jsonObj.Status = http.StatusOK
		jsonObj.Data = fmt.Sprintf("Movie with ID: %s was delete", *movie.Key)
		status = http.StatusOK
	}

	return c.JSON(status, jsonObj)
}

func AllMoviesGet(c echo.Context) error {
	//utils.ConsoleLog(c)
	var jsonObj model.JsonResult
	var status int

	err := store.Open()
	dbName := "event_dev"
	db, err := store.GetDatabase(dbName)
	if err != nil {
		fmt.Println(err)
	}

	//get all movies
	data, _ := db.GetAllMovies()

	if len(data) == 0 {
		jsonObj.Status = http.StatusNotFound
		jsonObj.Data = "Not found elements"
		status = http.StatusNotFound
	} else {
		jsonObj.Status = http.StatusOK
		jsonObj.Data = data
		status = http.StatusOK
	}

	return c.JSON(status, jsonObj)
}