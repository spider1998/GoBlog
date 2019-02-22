package entity

const TableSysCron = "sys_cron"

type SysCron struct {
	ID             int    `json:"id"`
	Key            string `json:"key"`
	Spec           string `json:"spec"`
	LastExecutedAt string `json:"last_executed_at"`
	DatetimeAware
}

func (SysCron) TableName() string {
	return TableSysCron
}
