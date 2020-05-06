# sqlbuilder
`sqlbuilder` is a simple sql query string builder

sqlbuilder its recursive struct call, that you can easy to build sql string

[![GoDoc](https://godoc.org/github.com/eehsiao/sqlbuilder?status.svg)](https://godoc.org/github.com/eehsiao/sqlbuilder)

simple sample : 
```go
    b := sb.NewSQLBuilder("SQLite")
    b.Select("a", "b").
        From("tblA").
        Join("tblB").
        Where("a", "=", 1).
        Where("b", "=", "str").
        Limit(100).
        BuildSelectSQL()
    fmt.Println(b.BuildedSQL())

    b.Release()
```

more case can see [test case](https://github.com/eehsiao/sqlbuilder/blob/master/sqlbuilder_test.go)

# go-model

how to use sqlbuilder with [go-model](https://github.com/eehsiao/go-models) :
```go
package db

import (
	"database/sql"
	"strconv"

	model "github.com/eehsiao/go-models"
	sb "github.com/eehsiao/sqlbuilder"
)

func (dao *Dao) GetExgs(params ...string) (e *[]*Exg, err error) {
	var (
		b    = sb.NewSQLBuilder("SQLite")
		rows *sql.Rows
	)
	e = &[]*Exg{}
	defer func() {
		if rows != nil {
			rows.Close()
		}
		b.Release()
		rows = nil
	}()

	b.Select(model.Inst2Fields(Exg{})...).
		From(TbExgs).
		Where("is_active", "=", 1)
	if len(params) > 0 && params[0] != "" {
		b.Where("exg_code", "=", params[0])
	}
	if len(params) > 1 && params[1] != "" {
		if i, err := strconv.Atoi(params[1]); err == nil && i > 0 {
			b.Limit(i)
		}
	}
	if rows, err = dao.Get(b); err == nil {
		for rows.Next() {
			exg := Exg{}
			if err = rows.Scan(model.Struct4Scan(&exg)...); err == nil {
				*e = append(*e, &exg)
			} else {
				return
			}
		}
	}
	return
}

func (dao *Dao) UpdateOrInsertExg(t *Exg) (r sql.Result, err error) {
	var (
		b          = sb.NewSQLBuilder("SQLite")
		cnt        int64
		withoutKey = []string{"id", "exg_code"}
	)
	defer func() {
		b.Release()
	}()

	b.Set(model.Inst2Set(*t, withoutKey...)).
		From(TbExgs).
		Where("exg_code", "=", t.ExgCode).
		BuildUpdateSQL()
	if r, err = dao.ExecBuilder(b); err == nil {
		if cnt, err = r.RowsAffected(); err == nil && cnt == 0 {
			b.Fields(model.Inst2FieldWithoutID(*t)...).
				Values(model.Inst2Values(*t, "id")...).
				Into(TbExgs).
				BuildInsertSQL()
			r, err = dao.ExecBuilder(b)
		}
	}

	return
}
```
