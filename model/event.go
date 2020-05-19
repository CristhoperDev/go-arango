package model

type Event struct {
	Key				string 		`json:"_key"`
	Action			string 		`json:"action"`
	UserID 			string 		`json:"userID"`
	CreatedAt 		*string 	`json:"createdAt"`
	FlagMail		bool 		`json:"flagMail"`
	ObjectName 		string 		`json:"objectName"`
	ObjectId 		string 		`json:"objectId"`
	Object 			interface{} `json:"object"`
}

type EventLog struct {
	Key	string `json:"_key"`
	Year int `json:"year"`
	Month int `json:"month"`
	Day int `json:"day"`
	Hour int `json:"hour"`
	CreatedAt *string `json:"createdAt"`
	Events []*Event `json:"events"`
}