package model

type Movie struct {
	Key 			*string 	`json:"_key"`
	Title 			string 		`json:"title"`
	Description 	string 		`json:"description"`
	CreatedAt 		*string 	`json:"createdAt"`
	UpdatedAt 		*string 	`json:"updatedAt"`
}

type Employee struct {
	Key 			*string 	`json:"_key"`
	Name 			string 		`json:"name"`
	LastName 		string 		`json:"lastName"`
	Degree	 		string 		`json:"degree"`
	FlagWork	 	bool 		`json:"flagWork"`
	CreatedAt 		*string 	`json:"createdAt"`
	UpdatedAt 		*string 	`json:"updatedAt"`
}
