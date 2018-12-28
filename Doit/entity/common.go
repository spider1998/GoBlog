package entity

import "Project/Doit/util"

const (
	GenderMale   = 1
	GenderFemale = 2
)

var GenderList = util.Pairs{
	{GenderMale, "男"},
	{GenderFemale, "女"},
}

type DatetimeAware struct {
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func DatetimeAwareNow() DatetimeAware {
	return DatetimeAware{util.DateTimeStd(), util.DateTimeStd()}
}
