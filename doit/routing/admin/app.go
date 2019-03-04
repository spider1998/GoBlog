package admin

import (
	"net/http"
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/entity"
	"Project/doit/service"
	"sort"
	"strconv"
	"strings"
	"time"

	"Project/doit/util"

	"github.com/go-ozzo/ozzo-dbx"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/robfig/cron"
)

func makeBackup(c *routing.Context) error {
	return service.BackupApp.Backup()
}

func restoreBackup(c *routing.Context) error {
	var req struct {
		Datetime string `json:"datetime"`
	}
	err := c.Read(&req)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest)
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", req.Datetime, time.Local)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err("invalid datetime.")
	}

	err = service.BackupApp.Restore(service.BackupApp.TranslateTimeToBackupFilename(t))
	if err != nil {
		return err
	}
	return c.Write(http.StatusOK)
}

type Spec struct {
	Hour           string `json:"hour"`
	Day            string `json:"day"`
	Month          string `json:"month"`
	Week           string `json:"week"`
	LastExecutedAt string `json:"last_executed_at"`
	NextExecutedAt string `json:"next_executed_at"`
}

func updateSchedule(c *routing.Context) error {
	key := c.Param("key")
	if key != service.CronDBBackupKey {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err("invalid schedule key.")
	}

	var specReq Spec

	err := c.Read(&specReq)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}

	spec := strings.Join([]string{"0", "0", specReq.Hour, specReq.Day, specReq.Month, specReq.Week}, " ")
	_, err = cron.Parse(spec)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err("invalid spec.")
	}

	_, err = app.DB.Update(entity.TableSysCron, dbx.Params{"update_time": util.DateTimeStd(), "spec": spec}, dbx.HashExp{"key": service.CronDBBackupKey}).Execute()
	if err != nil {
		return err
	}

	service.BackupApp.ReloadCron()

	return nil
}

func getSchedule(c *routing.Context) error {
	key := c.Param("key")
	if key != service.CronDBBackupKey {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err("invalid schedule key.")
	}

	var sysCron entity.SysCron
	err := app.DB.Select().Where(dbx.HashExp{"key": key}).One(&sysCron)
	if err != nil {
		return err
	}

	var spec Spec
	comps := strings.Split(sysCron.Spec, " ")
	if len(comps) != 6 {
		app.Logger.Error().Str("spec", sysCron.Spec).Msg("invalid spec of cron job.")
		return code.New(http.StatusInternalServerError, code.CodeServerErr)
	}
	spec.Hour = comps[2]
	spec.Day = comps[3]
	spec.Month = comps[4]
	spec.Week = comps[5]
	spec.LastExecutedAt = sysCron.LastExecutedAt
	schedule, err := cron.Parse(sysCron.Spec)
	if err != nil {
		return err
	}
	spec.NextExecutedAt = schedule.Next(time.Now()).Format("2006-01-02 15:04:05")

	return c.Write(spec)
}

func listBackups(c *routing.Context) error {
	page, size, err := util.ParsePagination(c)
	if err != nil {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err(err)
	}

	files, err := service.BackupApp.ListBackupFiles()
	if err != nil {
		return err
	}

	backupTimes := make([]string, 0)
	for _, file := range files {
		t, err := service.BackupApp.TranslateTimeFromBackupFilename(file)
		if err != nil {
			return err
		}
		backupTimes = append(backupTimes, t.Format("2006-01-02 15:04:05"))
	}

	sort.Sort(sort.Reverse(sort.StringSlice(backupTimes)))

	total := len(backupTimes)

	c.Response.Header().Set("X-Total-Count", strconv.Itoa(total))

	from := int((page - 1) * size)
	if from > total {
		from = total
	}
	to := from + int(size)
	if to > total {
		to = total
	}

	return c.Write(backupTimes[from:to])
}
