package utils

import (
	"fmt"
	"github.com/cristhoperdev/events-import/connection"
	"github.com/cristhoperdev/events-import/model"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo"
	"time"
)

var store connection.Datastore
var events = []*model.Event{}

// NewUUID creates a new unique universal identifier
func NewUUID() (string, error) {

	result, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

// NowUTC returns the current date time in UTC and in RFC3339 format
func NowUTC() *string {
	date := time.Now().UTC().Format(time.RFC3339)
	return &date
}

//ConsoleLog comment
func ConsoleLog(c echo.Context) {
	timeMessage := NowUTC()
	fmt.Println(*timeMessage, "-", c.Request().Method, "-", c.Request().URL)
}
