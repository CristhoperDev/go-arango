package main

import (
	"fmt"
	"github.com/cristhoperdev/events-import/connection"
	"github.com/cristhoperdev/events-import/model"
	"github.com/cristhoperdev/events-import/utils"
	"time"
)

var (
	store connection.Datastore
)

func main() {
	//open connection
	err := store.Open()
	if err != nil {
		fmt.Println(err)
	}

	//getting database
	dbName := "event_dev"
	db, err := store.GetDatabase(dbName)
	if err != nil {
		fmt.Println(err)
	}

	Events := []*model.Event{}
	for i := 0; i <= 1000; i++ {
		key, _ := utils.NewUUID()
		userID, _ := utils.NewUUID()
		Events = append(Events, &model.Event{
			Key:       key,
			Type:      "insert",
			UserID:    userID,
			CreatedAt: utils.NowUTC(),
		})
	}

	fmt.Println(len(Events))

	if len(Events) >= 100 {
		start := time.Now()
		//bulk import
		err = db.BulkImportEvents(Events[0:1000])
		if err != nil {
			fmt.Println(err)
		} else {
			Events = Events[1000:]
		}
		//elapsed time
		elapsed := time.Since(start)
		fmt.Printf("ArangoDB call took: %s \n", elapsed)

	}

	fmt.Println(len(Events))


	/*e := echo.New()
	e.GET("/movie", handler.AllMoviesGet)
	e.POST("/movie", handler.MoviePost)
	e.PUT("/movie", handler.MoviePut)
	e.DELETE("/movie", handler.MovieDelete)

	e.Logger.Fatal(e.Start(":3000"))*/
}
