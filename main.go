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

const LOG_FILE = "eventtest.json"

/*func init() {
	file, err := os.OpenFile(LOG_FILE, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	log.SetOutput(file)

	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)


	//defer file.Close()
}*/

/*func foo(event *model.Event) {
	log.WithFields(log.Fields{
		"_key":     	event.Key,
		"type":	 		event.Type,
		"userID": 		event.UserID,
		"action":	 	event.Action,
		"flagMail":	 	event.FlagMail,
		"createdAt": 	event.CreatedAt,
	}).Info(event.Type)
}*/

func main() {
	/*l := &lumberjack.Logger{
		Filename:   LOG_FILE,
		MaxSize:    1,
		MaxAge:     0,
		MaxBackups: 0,
		LocalTime:  true,
		Compress:   false,
	}
	log.SetOutput(l)
	log.SetFormatter(&log.JSONFormatter{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	go func() {
		for {
			<-c
			l.Rotate()
		}
	}()
	for i := 0; i < 10000; i++ {
		log.WithFields(log.Fields{
			"_key":     	"test",
		}).Info("test")
	}*/
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

	key, _ := utils.NewUUID()
	Now := time.Now()

	EventLog := &model.EventLog{
		Key: key,
		Year:   Now.Year(),
		Month:  int(Now.Month()),
		Day:    Now.Day(),
		Hour:   Now.Hour(),
		CreatedAt: utils.NowUTC(),// to manage ttl
		Events: []*model.Event{},
	}

	_ = db.CreateEvent(EventLog)

	Events := []*model.Event{}
	for i := 0; i < 50; i++ {
		key, _ := utils.NewUUID()
		userID, _ := utils.NewUUID()
		objectID, _ := utils.NewUUID()
		movie := model.Movie{
			Key:         &objectID,
			Title:       "movie",
			Description: "movie test",
			CreatedAt:   utils.NowUTC(),
			UpdatedAt:   utils.NowUTC(),
		}
		Events = append(Events, &model.Event{
			Key:       	key,
			Action:    	"insert",
			FlagMail:  	true,
			UserID:    	userID,
			CreatedAt: 	utils.NowUTC(),
			ObjectName: "Movie",
			ObjectId: 	objectID,
			Object: 	movie,
		})

	}
	for i := 0; i < 50; i++ {
		key, _ := utils.NewUUID()
		userID, _ := utils.NewUUID()
		objectID, _ := utils.NewUUID()
		employee := model.Employee{
			Key:       &objectID,
			Name:      "test 1",
			LastName:  "test 1",
			Degree:    "System",
			FlagWork:  true,
			CreatedAt: utils.NowUTC(),
			UpdatedAt: utils.NowUTC(),
		}
		Events = append(Events, &model.Event{
			Key:       	key,
			Action:    	"insert",
			FlagMail:  	true,
			UserID:    	userID,
			CreatedAt: 	utils.NowUTC(),
			ObjectName: "Profile",
			ObjectId: 	objectID,
			Object: 	employee,
		})

	}

	fmt.Println(*utils.NowUTC())

	fmt.Println("Init update array")
	start := time.Now()
	//bulk import
	err = db.UpdateDataEvent(key, Events)
	if err != nil {
		fmt.Println(err)
	} else {
		Events = nil
	}
	//elapsed time
	elapsed := time.Since(start)
	fmt.Printf("ArangoDB call bull took: %s \n", elapsed)

	/*Events = []*model.Event{}
	for i := 0; i < 5; i++ {
		key, _ := utils.NewUUID()
		userID, _ := utils.NewUUID()
		objectID, _ := utils.NewUUID()
		movie := model.Movie{
			Key:         &objectID,
			Title:       "movie",
			Description: "movie test",
			CreatedAt:   utils.NowUTC(),
			UpdatedAt:   utils.NowUTC(),
		}
		Events = append(Events, &model.Event{
			Key:       	key,
			Action:    	"insert",
			FlagMail:  	true,
			UserID:    	userID,
			CreatedAt: 	utils.NowUTC(),
			ObjectName: "Movie",
			ObjectId: 	objectID,
			Object: 	movie,
		})

	}

	for i := 0; i < 5; i++ {
		key, _ := utils.NewUUID()
		userID, _ := utils.NewUUID()
		objectID, _ := utils.NewUUID()
		employee := model.Employee{
			Key:       &objectID,
			Name:      "test 1",
			LastName:  "test 1",
			Degree:    "System",
			FlagWork:  true,
			CreatedAt: utils.NowUTC(),
			UpdatedAt: utils.NowUTC(),
		}
		Events = append(Events, &model.Event{
			Key:       	key,
			Action:    	"insert",
			FlagMail:  	true,
			UserID:    	userID,
			CreatedAt: 	utils.NowUTC(),
			ObjectName: "Profile",
			ObjectId: 	objectID,
			Object: 	employee,
		})

	}

	fmt.Println("Init 10 data insert")
	start = time.Now()
	for len(Events) > 0 {
		err = db.CreateEvent(Events[0])
		if err != nil {
			fmt.Println(err)
		}
		Events = Events[1:]
	}
	//elapsed time
	elapsed = time.Since(start)
	fmt.Printf("ArangoDB call took: %s \n", elapsed)*/

	/*e := echo.New()
	e.GET("/movie", handler.AllMoviesGet)
	e.POST("/movie", handler.MoviePost)
	e.PUT("/movie", handler.MoviePut)
	e.DELETE("/movie", handler.MovieDelete)

	e.Logger.Fatal(e.Start(":3000"))*/
}
