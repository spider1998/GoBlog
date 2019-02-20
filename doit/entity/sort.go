package entity

const TableSort = "sort"

type Sort struct {
	Name 		string 		`json:"name"`
	Operator 	string 		`json:"operator"`
	CreateTime	string 		`json:"create_time"`
}


func (Sort) TableName() string {
	return TableSort
}

