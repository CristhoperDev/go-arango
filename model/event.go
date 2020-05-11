package model

type Event struct {
	Key				string 	`json:"_key"`
	Type 			string 	`json:"type"`
	UserID 			string 	`json:"userID"`
	CreatedAt 		*string `json:"createdAt"`
}