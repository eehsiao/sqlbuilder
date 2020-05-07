// Author :		Eric<eehsiao@gmail.com>

// Package sqlbuilder is a builder that easy to use to build SQL string
package sqlbuilder

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	errLog     = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	isErrorLog = false
)

// NewSQLBuilder can create a sqlbuilder object
func NewSQLBuilder(d ...string) (b *SQLBuilder) {
	driverName := "mysql"
	if len(d) > 0 && checkDriveType(d[0]) {
		driverName = d[0]
	}
	b = &SQLBuilder{
		driverType: driverName,
	}
	b.ClearBuilder()

	return
}

func checkDriveType(d string) (b bool) {
	switch d {
	case "mysql":
		fallthrough
	case "mssql":
		fallthrough
	case "oracle":
		fallthrough
	case "postgresql":
		fallthrough
	case "SQLite":
		b = true
	}

	return
}

// EscapeStr is a internal function
func EscapeStr(value string, isMysql ...bool) string {
	replace := []struct {
		org string
		rep string
	}{
		{`\\`, `\\\\`},
		{`'`, `\'`},
		{`"`, `\"`},
		{`\\0`, "\\\\0"},
		{`\n`, `\\n`},
		{`\r`, `\\r`},
		{`\x1a`, `\\Z`},
	}

	if len(isMysql) > 0 && !isMysql[0] {
		replace = []struct {
			org string
			rep string
		}{
			{`\\`, `\\\\`},
			{`'`, `''`},
			{`"`, `""`},
			{`\\0`, "\\\\0"},
			{`\n`, `\\n`},
			{`\r`, `\\r`},
			{`\x1a`, `\\Z`},
		}
	}

	for _, r := range replace {
		value = strings.Replace(value, r.org, r.rep, -1)
	}

	return value
}

// SwitchPanicToErrorLog is a internal function
func (sb *SQLBuilder) SwitchPanicToErrorLog(b bool) {
	isErrorLog = b
}

// PanicOrErrorLog is a internal function
func (sb *SQLBuilder) PanicOrErrorLog(s string) {
	if isErrorLog {
		errLog.Println(s)
	} else {
		panic(s)
	}
}

// SetDriverType is a internal function
func (sb *SQLBuilder) SetDriverType(t string) {
	if checkDriveType(t) {
		sb.driverType = t
	}
}

// ClearBuilder that reset the sqlbuilder
func (sb *SQLBuilder) ClearBuilder() {
	sb.distinct = false
	sb.buildedStr = ""
	sb.selects = make([]string, 0)
	sb.froms = make([]string, 0)
	sb.joins = make([]string, 0)
	sb.wheres = make([]string, 0)
	sb.orders = make([]string, 0)
	sb.groups = make([]string, 0)
	sb.havings = ""
	sb.limit = ""
	sb.top = ""
	sb.into = ""
	sb.fields = make([]string, 0)
	sb.values = make([][]interface{}, 0)
	sb.sets = make([]Set, 0)
}

// SetDbName set a default db name
func (sb *SQLBuilder) SetDbName(s string) {
	if s == "" {
		return
	}

	sb.dbName = s
}

// SetTbName set a default table name
func (sb *SQLBuilder) SetTbName(s string) {
	if s == "" {
		return
	}

	sb.tbName = s
}

// GetDbName return the default db name
func (sb *SQLBuilder) GetDbName() string {
	return sb.dbName
}

// GetTbName return the default table name
func (sb *SQLBuilder) GetTbName() string {
	return sb.tbName
}

// IsMysql return the builder engine is for mysql
func (sb *SQLBuilder) IsMysql() bool {
	return sb.driverType == "mysql"
}

// IsMssql return the builder engine is for mssql
func (sb *SQLBuilder) IsMssql() bool {
	return sb.driverType == "mssql"
}

// IsOracle return the builder engine is for oracle
func (sb *SQLBuilder) IsOracle() bool {
	return sb.driverType == "oracle"
}

// IsPostgresql return the builder engine is for postgresql
func (sb *SQLBuilder) IsPostgresql() bool {
	return sb.driverType == "postgresql"
}

// IsSQLite return the builder engine is for SQLite
func (sb *SQLBuilder) IsSQLite() bool {
	return sb.driverType == "SQLite"
}

// IsDistinct is internal function
func (sb *SQLBuilder) IsDistinct() bool {
	return sb.distinct
}

// IsHasSelects is internal function
func (sb *SQLBuilder) IsHasSelects() bool {
	return len(sb.selects) > 0
}

// IsHasDbName is internal function
func (sb *SQLBuilder) IsHasDbName() bool {
	return sb.dbName != ""
}

// IsHasTbName is internal function
func (sb *SQLBuilder) IsHasTbName() bool {
	return sb.tbName != ""
}

// IsHasOneFroms is internal function
func (sb *SQLBuilder) IsHasOneFroms() bool {
	return len(sb.froms) == 1
}

// IsHasFroms is internal function
func (sb *SQLBuilder) IsHasFroms() bool {
	return len(sb.froms) > 0
}

// IsHasJoins is internal function
func (sb *SQLBuilder) IsHasJoins() bool {
	return len(sb.joins) > 0
}

// IsHasWheres is internal function
func (sb *SQLBuilder) IsHasWheres() bool {
	return len(sb.wheres) > 0
}

// IsHasOrders is internal function
func (sb *SQLBuilder) IsHasOrders() bool {
	return len(sb.orders) > 0
}

// IsHasGroups is internal function
func (sb *SQLBuilder) IsHasGroups() bool {
	return len(sb.groups) > 0
}

// IsHasHavings is internal function
func (sb *SQLBuilder) IsHasHavings() bool {
	return sb.havings != ""
}

// IsHasLimit is internal function
func (sb *SQLBuilder) IsHasLimit() bool {
	return (sb.IsMysql() || sb.IsSQLite()) && sb.limit != ""
}

// IsHasTop is internal function
func (sb *SQLBuilder) IsHasTop() bool {
	return sb.IsMssql() && sb.top != ""
}

// IsHasInto is internal function
func (sb *SQLBuilder) IsHasInto() bool {
	return sb.into != ""
}

// IsHasFields is internal function
func (sb *SQLBuilder) IsHasFields() bool {
	return len(sb.fields) > 0
}

// IsHasValues is internal function
func (sb *SQLBuilder) IsHasValues() bool {
	return len(sb.values) > 0
}

// IsHasSets is internal function
func (sb *SQLBuilder) IsHasSets() bool {
	return len(sb.sets) > 0
}

// IsHadBuildedSQL can know has BuildedSQL string or not
func (sb *SQLBuilder) IsHadBuildedSQL() bool {
	return sb.buildedStr != ""
}

// GetFieldsCount is internal function
func (sb *SQLBuilder) GetFieldsCount() int {
	return len(sb.fields)
}

// BuildedSQL return the builded SQL string
func (sb *SQLBuilder) BuildedSQL() (sql string) {
	return sb.buildedStr
}

// BuildDeleteSQL do build the `delete` SQL string
func (sb *SQLBuilder) BuildDeleteSQL() *SQLBuilder {
	if !sb.CanBuildDelete() {
		sb.PanicOrErrorLog("must be have only one from table or default TbName")
	}
	sql := ""
	if sb.IsHasOneFroms() {
		sql = "DELETE FROM " + sb.froms[0]
	} else if sb.IsHasTbName() {
		sql = "DELETE FROM " + sb.tbName
	}

	if sb.IsHasWheres() {
		sql += " WHERE " + strings.Join(sb.wheres, " ")
	}
	sb.buildedStr = sql

	return sb
}

// BuildSelectSQL do build the `select` SQL string
func (sb *SQLBuilder) BuildSelectSQL() *SQLBuilder {
	if !sb.CanBuildSelect() {
		sb.PanicOrErrorLog("Without selects or from table is not set")
	}

	sql := "SELECT"

	if sb.IsDistinct() {
		sql += " DISTINCT"
	}

	if sb.IsHasTop() {
		sql += " TOP " + sb.top
	}

	sql += " " + strings.Join(sb.selects, ",")
	if sb.IsHasFroms() {
		sql += " FROM " + strings.Join(sb.froms, ",")
	} else if sb.IsHasTbName() {
		sql += " FROM " + sb.tbName
	}

	if sb.IsHasJoins() {
		sql += " " + strings.Join(sb.joins, " ")
	}

	if sb.IsHasWheres() {
		sql += " WHERE " + strings.Join(sb.wheres, " ")
	}

	if sb.IsHasOrders() {
		sql += " ORDER BY " + strings.Join(sb.orders, ",")
	}

	if sb.IsHasGroups() {
		sql += " GROUP BY " + strings.Join(sb.groups, ",")
	}

	if sb.IsHasHavings() {
		sql += " HAVING " + sb.havings
	}

	if sb.IsHasLimit() {
		sql += " LIMIT " + sb.limit
	}

	sb.buildedStr = sql

	return sb
}

// BuildUpdateSQL do build the `update` SQL string
func (sb *SQLBuilder) BuildUpdateSQL() *SQLBuilder {
	if !sb.CanBuildUpdate() {
		sb.PanicOrErrorLog("Without update table or default TbName")
	}

	sql := "UPDATE "
	if sb.IsHasOneFroms() {
		sql += sb.froms[0] + " "
	} else if sb.IsHasTbName() {
		sql += sb.tbName + " "
	}

	setStr := "SET "
	for _, set := range sb.sets {
		switch set.V.(type) {
		case SQLVar:
			setStr += fmt.Sprintf("%s=%s,", EscapeStr(set.K, sb.IsMysql()), (set.V.(SQLVar)).VarS)
		case string:
			setStr += fmt.Sprintf("%s='%v',", EscapeStr(set.K, sb.IsMysql()), EscapeStr(set.V.(string), sb.IsMysql()))
		default:
			if set.V == nil {
				setStr += fmt.Sprintf("%s=NULL,", EscapeStr(set.K, sb.IsMysql()))
			} else {
				setStr += fmt.Sprintf("%s=%v,", EscapeStr(set.K, sb.IsMysql()), set.V)
			}
		}
	}
	sql += strings.Trim(setStr, ",")

	if sb.IsHasWheres() {
		sql += " WHERE " + strings.Join(sb.wheres, " ")
	}
	sb.buildedStr = sql

	return sb
}

// BuildInsertSQL do build the `insert` SQL string
func (sb *SQLBuilder) BuildInsertSQL() *SQLBuilder {
	if !sb.CanBuildInsert() {
		sb.PanicOrErrorLog("Without insert table or default TbName")
	}

	sql := "INSERT INTO "
	if sb.IsHasInto() {
		sql += sb.into
	} else if sb.IsHasTbName() {
		sql += sb.tbName
	}

	sql += " (" + strings.Join(sb.fields, ",") + ") VALUES "
	vals := "("
	for _, v := range sb.values[0] {
		switch v.(type) {
		case SQLVar:
			vals += fmt.Sprintf("%s,", v.(SQLVar).VarS)
		case string:
			vals += fmt.Sprintf("'%v',", EscapeStr(v.(string), sb.IsMysql()))
		default:
			if v == nil {
				vals += "NULL,"
			} else {
				vals += fmt.Sprintf("%v,", v)
			}
		}
	}
	sql += strings.Trim(vals, ",") + ")"

	sb.buildedStr = sql

	return sb
}

// BuildBulkInsertSQL do build the `insert` SQL string with bulk values
func (sb *SQLBuilder) BuildBulkInsertSQL() *SQLBuilder {
	if !sb.CanBuildInsert() {
		sb.PanicOrErrorLog("Without insert table or default TbName")
	}

	sql := "INSERT INTO "
	if sb.IsHasInto() {
		sql += sb.into
	} else if sb.IsHasTbName() {
		sql += sb.tbName
	}

	sql += " (" + strings.Join(sb.fields, ",") + ") VALUES "
	vals := ""
	for _, vs := range sb.values {
		vals += "("
		for _, v := range vs {
			switch v.(type) {
			case SQLVar:
				vals += fmt.Sprintf("%s,", v.(SQLVar).VarS)
			case string:
				vals += fmt.Sprintf("'%v',", EscapeStr(v.(string), sb.IsMysql()))
			default:
				if v == nil {
					vals += "NULL,"
				} else {
					vals += fmt.Sprintf("%v,", v)
				}
			}
		}
		vals = strings.Trim(vals, ",") + "),"
	}
	sql += strings.Trim(vals, ",")

	sb.buildedStr = sql

	return sb
}

// BuildInsertOrReplaceSQL do build the `insert or replace into` SQL string
// only for SQLite
func (sb *SQLBuilder) BuildInsertOrReplaceSQL() *SQLBuilder {
	if !sb.IsSQLite() {
		sb.PanicOrErrorLog("InsertOrReplace only support SQLite")
	}
	if !sb.CanBuildInsert() {
		sb.PanicOrErrorLog("Without insert table or default TbName")
	}

	sql := "INSERT OR REPLACE INTO "
	if sb.IsHasInto() {
		sql += sb.into
	} else if sb.IsHasTbName() {
		sql += sb.tbName
	}

	sql += " (" + strings.Join(sb.fields, ",") + ") VALUES "

	vals := "("
	for _, v := range sb.values[0] {
		switch v.(type) {
		case string:
			vals += fmt.Sprintf("'%v',", EscapeStr(v.(string), sb.IsMysql()))
		default:
			if v == nil {
				vals += "NULL,"
			} else {
				vals += fmt.Sprintf("%v,", v)
			}
		}
	}
	sql += strings.Trim(vals, ",") + ")"
	sb.buildedStr = sql

	return sb
}

// CanBuildSelect is internal function
// that sqlbuilder to know can build select string or not
func (sb *SQLBuilder) CanBuildSelect() bool {
	return sb.IsHasSelects() && (sb.IsHasFroms() || sb.IsHasTbName())
}

// CanBuildDelete is internal function
// that sqlbuilder to know can build delete string or not
func (sb *SQLBuilder) CanBuildDelete() bool {
	return sb.IsHasFroms() || sb.IsHasTbName()
}

// CanBuildUpdate is internal function
// that sqlbuilder to know can build update string or not
func (sb *SQLBuilder) CanBuildUpdate() bool {
	return (sb.IsHasOneFroms() || sb.IsHasTbName()) && sb.IsHasSets()
}

// CanBuildInsert is internal function
// that sqlbuilder to know can build insert string or not
func (sb *SQLBuilder) CanBuildInsert() bool {
	return (sb.IsHasInto() || sb.IsHasTbName()) && sb.IsHasFields() && sb.IsHasValues()
}

// Release this object
func (sb *SQLBuilder) Release() {
	sb = nil
}
