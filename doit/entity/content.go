package entity

const TableContent = "content"

type Content struct {
	ArtId    string `json:"art_id"`                  //文章Id
	UserId   string `json:"user_id"`                 //用户id
	Version  int    `json:"version"`                 //版本标识
	HeadUuid string `json:"head_uuid"`               //头标识
	TailUuid string `json:"tail_uuid"`               //尾标识
	Detail   string `json:"detail"`                  //内容
	Changed  bool   `json:"changed" default:"false"` //改动标识
}

func (Content) TableName() string {
	return TableContent
}
