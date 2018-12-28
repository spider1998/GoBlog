package entity

const TableArtChange = "art_change"

type ArtChange struct {
	ArtId      string `json:"art_id"`      //文章Id
	ChangeId   string `json:"change_id"`   //改动人Id
	Name       string `json:"name"`        //改动人姓名
	State      string `json:"state"`       //改动状态
	UpdateTime string `json:"update_time"` //修改时间
	AgreeTime  string `json:"agree_time"`  //同意修改时间
}

func (ArtChange) TableName() string {
	return TableArtChange
}
