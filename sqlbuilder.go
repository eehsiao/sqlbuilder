// Author :		Eric<eehsiao@gmail.com>

package sqlbuilder

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	errLog     = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	isErrorLog = false
)

type Set struct {
	K string
	V interface{}
}
type SQLBuilder struct {
	driverType string

	// default database and table
	dbName     string
	tbName     string
	buildedStr string

	// for select , delete
	distinct bool
	selects  []string
	froms    []string
	joins    []string
	wheres   []string
	orders   []string
	groups   []string
	havings  string
	limit    string
	top      string

	// for insert
	into   string
	fields []string
	values []interface{}

	// for update
	sets []Set
}

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

func (sb *SQLBuilder) SwitchPanicToErrorLog(b bool) {
	isErrorLog = b
}

func (sb *SQLBuilder) PanicOrErrorLog(s string) {
	if isErrorLog {
		errLog.Println(s)
	} else {
		panic(s)
	}
}

func (sb *SQLBuilder) SetDriverType(t string) {
	if checkDriveType(t) {
		sb.driverType = t
	}
}

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
	sb.values = make([]interface{}, 0)
	sb.sets = make([]Set, 0)
}

func (sb *SQLBuilder) SetDbName(s string) {
	if s == "" {
		return
	}

	sb.dbName = s
}

func (sb *SQLBuilder) SetTbName(s string) {
	if s == "" {
		return
	}

	sb.tbName = s
}

func (sb *SQLBuilder) GetDbName() string {
	return sb.dbName
}

func (sb *SQLBuilder) GetTbName() string {
	return sb.tbName
}

func (sb *SQLBuilder) IsMysql() bool {
	return sb.driverType == "mysql"
}

func (sb *SQLBuilder) IsMssql() bool {
	return sb.driverType == "mssql"
}

func (sb *SQLBuilder) IsOracle() bool {
	return sb.driverType == "oracle"
}

func (sb *SQLBuilder) IsPostgresql() bool {
	return sb.driverType == "postgresql"
}

func (sb *SQLBuilder) IsSQLite() bool {
	return sb.driverType == "SQLite"
}

func (sb *SQLBuilder) IsDistinct() bool {
	return sb.distinct
}

func (sb *SQLBuilder) IsHasSelects() bool {
	return len(sb.selects) > 0
}

func (sb *SQLBuilder) IsHasDbName() bool {
	return sb.dbName != ""
}

func (sb *SQLBuilder) IsHasTbName() bool {
	return sb.tbName != ""
}

func (sb *SQLBuilder) IsHasOneFroms() bool {
	return len(sb.froms) == 1
}

func (sb *SQLBuilder) IsHasFroms() bool {
	return len(sb.froms) > 0
}

func (sb *SQLBuilder) IsHasJoins() bool {
	return len(sb.joins) > 0
}

func (sb *SQLBuilder) IsHasWheres() bool {
	return len(sb.wheres) > 0
}

func (sb *SQLBuilder) IsHasOrders() bool {
	return len(sb.orders) > 0
}

func (sb *SQLBuilder) IsHasGroups() bool {
	return len(sb.groups) > 0
}

func (sb *SQLBuilder) IsHasHavings() bool {
	return sb.havings != ""
}

func (sb *SQLBuilder) IsHasLimit() bool {
	return (sb.IsMysql() || sb.IsSQLite()) && sb.limit != ""
}

func (sb *SQLBuilder) IsHasTop() bool {
	return sb.IsMssql() && sb.top != ""
}

func (sb *SQLBuilder) IsHasInto() bool {
	return sb.into != ""
}

func (sb *SQLBuilder) IsHasFields() bool {
	return len(sb.fields) > 0
}

func (sb *SQLBuilder) IsHasValues() bool {
	return len(sb.values) > 0
}

func (sb *SQLBuilder) IsHasSets() bool {
	return len(sb.sets) > 0
}

func (sb *SQLBuilder) IsHadBuildedSQL() bool {
	return sb.buildedStr != ""
}

func (sb *SQLBuilder) GetFieldsCount() int {
	return len(sb.fields)
}

func (sb *SQLBuilder) BuildedSQL() (sql string) {
	return sb.buildedStr
}

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

func (sb *SQLBuilder) BuildInsertSQL() *SQLBuilder {
	if !sb.CanBuildInsert() {
		sb.PanicOrErrorLog("Without insert tableor default TbName")
	}

	sql := "INSERT INTO "
	if sb.IsHasInto() {
		sql += sb.into
	} else if sb.IsHasTbName() {
		sql += sb.tbName
	}

	sql += " (" + strings.Join(sb.fields, ",") + ") VALUES "

	vals := "("
	for _, v := range sb.values {
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

func (sb *SQLBuilder) BuildInsertOrReplaceSQL() *SQLBuilder {
	if !sb.IsSQLite() {
		sb.PanicOrErrorLog("limit only support SQLite")
	}
	if !sb.CanBuildInsert() {
		sb.PanicOrErrorLog("Without insert tableor default TbName")
	}

	sql := "INSERT OR REPLACE INTO "
	if sb.IsHasInto() {
		sql += sb.into
	} else if sb.IsHasTbName() {
		sql += sb.tbName
	}

	sql += " (" + strings.Join(sb.fields, ",") + ") VALUES "

	vals := "("
	for _, v := range sb.values {
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

func (sb *SQLBuilder) CanBuildSelect() bool {
	return sb.IsHasSelects() && (sb.IsHasFroms() || sb.IsHasTbName())
}

func (sb *SQLBuilder) CanBuildDelete() bool {
	return sb.IsHasFroms() || sb.IsHasTbName()
}

func (sb *SQLBuilder) CanBuildUpdate() bool {
	return (sb.IsHasOneFroms() || sb.IsHasTbName()) && sb.IsHasSets()
}

func (sb *SQLBuilder) CanBuildInsert() bool {
	return (sb.IsHasInto() || sb.IsHasTbName()) && sb.IsHasFields() && sb.IsHasValues()
}

func (sb *SQLBuilder) Distinct(b bool) *SQLBuilder {
	sb.distinct = b

	return sb
}

func (sb *SQLBuilder) Limit(i ...int) *SQLBuilder {
	if !(sb.IsMysql() || sb.IsSQLite()) {
		sb.PanicOrErrorLog("limit only support mysql")
	}
	if len(i) == 0 {
		sb.PanicOrErrorLog("must have value for limit")
	}

	if len(i) == 1 {
		sb.limit = strconv.Itoa(i[0])
	} else if len(i) > 1 {
		sb.limit = strconv.Itoa(i[0]) + "," + strconv.Itoa(i[1])
	}

	return sb
}

func (sb *SQLBuilder) Top(i int) *SQLBuilder {
	if !sb.IsMssql() {
		sb.PanicOrErrorLog("limit only support mssql")
	}
	if i <= 0 {
		sb.PanicOrErrorLog("must have >=1 value for top")
	}

	sb.limit = strconv.Itoa(i)

	return sb
}

func (sb *SQLBuilder) Select(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support fileds")
	}

	for _, v := range s {
		sb.selects = append(sb.selects, EscapeStr(v, sb.IsMysql()))
	}

	return sb
}

func (sb *SQLBuilder) From(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support tables")
	}

	for _, v := range s {
		sb.froms = append(sb.froms, EscapeStr(v, sb.IsMysql()))
	}

	return sb
}

func (sb *SQLBuilder) clearFrom() {
	sb.froms = make([]string, 0)
}

func (sb *SQLBuilder) FromOne(s string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support tables")
	}

	if sb.IsHasFroms() {
		sb.clearFrom()
	}

	sb.froms = append(sb.froms, EscapeStr(s, sb.IsMysql()))

	return sb
}

// Where : just equal WhereAnd
func (sb *SQLBuilder) Where(s string, o string, v interface{}) *SQLBuilder {
	return sb.WhereAnd(s, o, v)
}

func (sb *SQLBuilder) WhereAnd(s string, o string, v interface{}) *SQLBuilder {
	if len(s) == 0 && s != "" && o != "" {
		sb.PanicOrErrorLog("must be support conditions")
	}

	c := ""
	if sb.IsHasWheres() {
		c = "AND "
	}
	switch v.(type) {
	case string:
		sb.wheres = append(sb.wheres, fmt.Sprintf("%s%s %s '%s'", c, EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql()), EscapeStr(v.(string), sb.IsMysql())))
	default:
		if v == nil {
			sb.wheres = append(sb.wheres, fmt.Sprintf("%s%s %s NULL", c, EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql())))
		} else {
			sb.wheres = append(sb.wheres, fmt.Sprintf("%s%s %s %v", c, EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql()), v))
		}
	}

	return sb
}

func (sb *SQLBuilder) WhereOr(s string, o string, v interface{}) *SQLBuilder {
	if len(s) == 0 && s != "" && o != "" {
		sb.PanicOrErrorLog("must be support conditions")
	}

	c := ""
	if sb.IsHasWheres() {
		c = "OR "
	}
	switch v.(type) {
	case string:
		sb.wheres = append(sb.wheres, fmt.Sprintf("%s%s %s '%s'", c, EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql()), EscapeStr(v.(string), sb.IsMysql())))
	default:
		if v == nil {
			sb.wheres = append(sb.wheres, fmt.Sprintf("%s%s %s NULL", c, EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql())))
		} else {
			sb.wheres = append(sb.wheres, fmt.Sprintf("%s%s %s %v", c, EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql()), v))
		}
	}

	return sb
}

func (sb *SQLBuilder) join(p string, j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	sb.joins = append(sb.joins, p+"JOIN "+EscapeStr(j, sb.IsMysql()))

	return sb
}

func (sb *SQLBuilder) joinOn(p string, j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	switch v.(type) {
	case string:
		sb.joins = append(sb.joins, fmt.Sprintf("%sJOIN %s ON %s %s '%s'", p, j, EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql()), EscapeStr(v.(string), sb.IsMysql())))

	default:
		sb.joins = append(sb.joins, fmt.Sprintf("%sJOIN %s ON %s %s %v", p, j, EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql()), v))
	}

	return sb
}

func (sb *SQLBuilder) Join(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("", j)
}

func (sb *SQLBuilder) JoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("", j, s, o, v)
}

func (sb *SQLBuilder) InnerJoin(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("INNER ", j)
}

func (sb *SQLBuilder) InnerJoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("INNER ", j, s, o, v)
}

func (sb *SQLBuilder) LeftJoin(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("LEFT ", j)
}

func (sb *SQLBuilder) LeftJoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("LEFT ", j, s, o, v)
}

func (sb *SQLBuilder) RightJoin(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("RIGHT ", j)
}

func (sb *SQLBuilder) RightJoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("RIGHT ", j, s, o, v)
}

func (sb *SQLBuilder) FullJoin(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("FULL ", j)
}

func (sb *SQLBuilder) FullJoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("FULL ", j, s, o, v)
}

func (sb *SQLBuilder) GroupBy(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support group fileds")
	}

	for _, v := range s {
		sb.groups = append(sb.groups, EscapeStr(v, sb.IsMysql()))
	}

	return sb
}

func (sb *SQLBuilder) OrderBy(s ...string) *SQLBuilder {
	return sb.OrderByAsc(s...)
}

func (sb *SQLBuilder) OrderByAsc(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support order fileds")
	}

	sb.orders = append(sb.orders, EscapeStr(strings.Join(s, ","), sb.IsMysql())+" ASC")

	return sb
}

func (sb *SQLBuilder) OrderByDesc(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support order fileds")
	}

	sb.orders = append(sb.orders, EscapeStr(strings.Join(s, ","), sb.IsMysql())+" DESC")

	return sb
}

func (sb *SQLBuilder) Having(s string, o string, v interface{}) *SQLBuilder {
	if s == "" || !sb.IsHasGroups() {
		sb.PanicOrErrorLog("must be support having condition or set group by first")
	}

	sb.havings = EscapeStr(s, sb.IsMysql())

	switch v.(type) {
	case string:
		sb.havings = fmt.Sprintf("%s %s '%s'", EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql()), EscapeStr(v.(string), sb.IsMysql()))

	default:
		sb.havings = fmt.Sprintf("%s %s %v", EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql()), v)
	}

	return sb
}

func (sb *SQLBuilder) Set(s []Set) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support set fileds : values")
	}

	for _, set := range s {
		sb.sets = append(sb.sets, set)
	}

	return sb
}

func (sb *SQLBuilder) Into(s string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support table")
	}

	sb.into = EscapeStr(s, sb.IsMysql())

	return sb
}

func (sb *SQLBuilder) Fields(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support fileds")
	}
	if sb.IsHasValues() {
		sb.PanicOrErrorLog("cannot add fileds after Values()")
	}

	for _, v := range s {
		sb.fields = append(sb.fields, EscapeStr(v, sb.IsMysql()))
	}

	return sb
}

func (sb *SQLBuilder) Values(s ...interface{}) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support fileds")
	}

	fieldCnt := sb.GetFieldsCount()
	if len(s) == fieldCnt {
		for _, v := range s {
			sb.values = append(sb.values, v)
		}

	} else {
		sb.PanicOrErrorLog("values count not equal fileds count")
	}

	return sb
}

func (sb *SQLBuilder) Release() {
	sb = nil
}
