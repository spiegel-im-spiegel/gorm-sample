//go:build run
// +build run

package main

import (
	"fmt"
	"os"
	"sample/orm"

	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
	"gorm.io/gorm"
)

func Run() exitcode.ExitCode {
	// create gorm.DB instance for PostgreSQL service
	gormCtx, err := orm.NewGORM()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitcode.Abnormal
	}
	defer gormCtx.Close()

	// delete all data in users table (dry run)
	tx := gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Exec("UPDATE users SET deleted_at=now() WHERE deleted_at IS NULL") // soft delete
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
