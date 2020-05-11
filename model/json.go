package model

type JsonResult struct {
	Status int 			`json:"status"`
	Data   interface{} 	`json:"data"`
}