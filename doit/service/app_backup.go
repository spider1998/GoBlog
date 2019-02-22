package service

import (
	"Project/doit/app"
	"Project/doit/code"
	"Project/doit/util"
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/Go-SQL-Driver/mysql"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	backupDir         = "/backups"
	backupSuffix      = ".dat"
	backupFilenameFMT = "20060102150405" + backupSuffix
)

var backupLock util.Lock

func (s *BackupAppService) Backup() (err error) {
	if !backupLock.Lock() {
		return code.New(http.StatusConflict, code.CodeTaskIsInProgress)
	}
	defer backupLock.Unlock()

	path := filepath.Join(app.Conf.RuntimePath, backupDir)
	err = util.MakeDirectory(path)
	if err != nil {
		return err
	}

	cfg, err := mysql.ParseDSN(app.Conf.Mysql)
	if err != nil {
		return err
	}
	hostAndPort := strings.Split(cfg.Addr, ":")
	if len(hostAndPort) != 2 {
		return errors.New("invalid addr")
	}
	backFile := filepath.Join(path, time.Now().Format(backupFilenameFMT))
	cmd := exec.Command("mysqldump", "-h"+hostAndPort[0], "-P"+hostAndPort[1], "-u"+cfg.User, "-p"+cfg.Passwd, "--single-transaction", cfg.DBName)
	file, err := os.Create(backFile)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		if err != nil {
			os.Remove(backFile)
		}
	}()

	w := gzip.NewWriter(file)
	defer w.Close()
	stderr := new(bytes.Buffer)
	cmd.Stdout = w
	cmd.Stderr = stderr
	err = cmd.Run()
	if err != nil {
		app.Logger.Error().Str("stderr", stderr.String()).Msg("fail to backup db.")
		return err
	}
	return nil
}

func (s *BackupAppService) Restore(filename string) (err error) {
	if !backupLock.Lock() {
		return code.New(http.StatusConflict, code.CodeTaskIsInProgress)
	}
	defer backupLock.Unlock()

	cfg, err := mysql.ParseDSN(app.Conf.Mysql)
	if err != nil {
		return err
	}
	hostAndPort := strings.Split(cfg.Addr, ":")
	if len(hostAndPort) != 2 {
		return errors.New("invalid addr")
	}

	path := filepath.Join(app.Conf.RuntimePath, backupDir, filename)

	if !util.PathExist(path) {
		return code.New(http.StatusBadRequest, code.CodeBadRequest).Err("file does not exist.")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	r, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer r.Close()

	cmd := exec.Command("mysql", "-h"+hostAndPort[0], "-P"+hostAndPort[1], "-u"+cfg.User, "-p"+cfg.Passwd, cfg.DBName)
	cmd.Stdin = r

	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	cmd.Stderr = buf
	err = cmd.Run()
	if err != nil {
		app.Logger.Error().Err(err).Str("stderr", buf.String()).Msg("fail to backup db.")
		return err
	}

	s.ReloadCron()

	return nil
}

func (s *BackupAppService) ListBackupFiles() (filenames []string, err error) {
	backupPath := filepath.Join(app.Conf.RuntimePath, backupDir)
	if !util.PathExist(backupPath) {
		return nil, nil
	}

	var files []string
	err = filepath.Walk(backupPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			name := info.Name()
			if strings.Contains(name, backupSuffix) {
				files = append(files, name)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *BackupAppService) TranslateTimeFromBackupFilename(name string) (time.Time, error) {
	return time.ParseInLocation(backupFilenameFMT, name, time.Local)
}

func (s *BackupAppService) TranslateTimeToBackupFilename(t time.Time) string {
	return t.Format(backupFilenameFMT)
}
