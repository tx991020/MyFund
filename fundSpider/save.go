package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	pg "github.com/go-pg/pg/v10"
	"github.com/tx991020/utils"
)

var (
	db    *pg.DB
	data  = strings.Split(utils.Time2Date(time.Now().AddDate(0, 0, -1)), "-")
	table = fmt.Sprintf("fund%s%s", data[1], data[2])

	createTable = `CREATE TABLE "public"."%s" (
  "c1" text COLLATE "pg_catalog"."default",
  "c2" text COLLATE "pg_catalog"."default",
  "c3" text COLLATE "pg_catalog"."default",
  "c4" text COLLATE "pg_catalog"."default",
  "c5" text COLLATE "pg_catalog"."default",
  "c6" text COLLATE "pg_catalog"."default",
  "c7" text COLLATE "pg_catalog"."default",
  "c8" text COLLATE "pg_catalog"."default",
  "c9" text COLLATE "pg_catalog"."default",
  "c10" text COLLATE "pg_catalog"."default",
  "c11" text COLLATE "pg_catalog"."default",
  "c12" text COLLATE "pg_catalog"."default",
  "c13" text COLLATE "pg_catalog"."default",
  "c14" text COLLATE "pg_catalog"."default",
  "c15" text COLLATE "pg_catalog"."default",
  "c16" text COLLATE "pg_catalog"."default",
  "c17" text COLLATE "pg_catalog"."default"
)
;
`

	rankTable = `CREATE TABLE "public"."%s" (
  "c1" text COLLATE "pg_catalog"."default",
  "c2" text COLLATE "pg_catalog"."default",
  "c3" text COLLATE "pg_catalog"."default",
  "c4" text COLLATE "pg_catalog"."default",
  "c5" text COLLATE "pg_catalog"."default",
  "c6" text COLLATE "pg_catalog"."default",
  "c7" text COLLATE "pg_catalog"."default",
  "c8" text COLLATE "pg_catalog"."default"
)
;
`
)

func DBInit() {




	db = pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Password: "123456",
		Database: "prest",
		PoolSize: 10,
	})

}

func CreateTableIfNotExist(db *pg.DB, table, sql string) (err error) {

	_, err = db.Exec(fmt.Sprintf(sql, table))
	if err != nil {
		return

	}
	return
}

func CopyToPG(path, table ,sql1 string) (err error) {

	err = CreateTableIfNotExist(db, table, sql1)
	if err != nil {
		return
	}
	buffer := bytes.NewBuffer(nil)
	defer buffer.Reset()
	file, err := ioutil.ReadFile(path)
	if err != nil {

		return
	}
	buffer.Write(file)
	sql := fmt.Sprintf("COPY %s FROM STDIN WITH CSV", table)
	fmt.Println(sql)
	_, err = db.CopyFrom(buffer, sql)

	return
}
