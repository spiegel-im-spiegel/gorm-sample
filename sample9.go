//go:build run
// +build run

package main

import (
	"context"
	"fmt"
	"os"
	"sample/orm"
	"sample/orm/model"

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

	// select data for 'Bob 2nd'
	data := model.User{}
	tx := gormCtx.GetDb().WithContext(context.TODO()).Where(&model.User{Username: "Bob 2nd"}).First(&data)
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	// delete data in users table (dry run)
	tx = gormCtx.GetDb().Session(&gorm.Session{DryRun: true}).Delete(&data)
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
