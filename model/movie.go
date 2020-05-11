package model

type Movie struct {
	Key 			*string 	`json:"_key"`
	Title 			string 		`json:"title"`
	Description 	string 		`json:"description"`
	CreatedAt 		*string 	`json:"createdAt"`
	UpdatedAt 		*string 	`json:"updatedAt"`
}
