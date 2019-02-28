package entity

const TableSort = "sort"

type SortState int8

const (
	SortStateAble SortState = iota + 1
	SortStateEnable
)

type Sort struct {
	ID 			string 		`json:"id" gorm:"pk"`
	Name 		string 		`json:"name"`
	Operator 	string 		`json:"operator"`
	CreateTime	string 		`json:"create_time"`
	State 		SortState 	`json:"state"`
	Sum 		int 		`json:"sum"`
}


func (Sort) TableName() string {
	return TableSort
}

