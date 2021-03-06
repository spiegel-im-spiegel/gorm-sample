//go:build run
// +build run

package main

import (
	"context"
	"fmt"
	"os"
	"sample/files"
	"sample/orm"
	"sample/orm/model"

	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
	"gorm.io/gorm/clause"
)

func Run() exitcode.ExitCode {
	// create gorm.DB instance for PostgreSQL service
	gormCtx, err := orm.NewGORM()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitcode.Abnormal
	}
	defer gormCtx.Close()

	// select all records (with preload)
	data := []model.User{}
	tx := gormCtx.GetDb().WithContext(context.TODO()).Preload(clause.Associations).Find(&data)
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	// output by JSON format
	if err := files.Output(os.Stdout, data); err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
