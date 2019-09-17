// Author :		Eric<eehsiao@gmail.com>

package sqlbuilder

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/eehsiao/go-models/lib"
)

var (
	errLog     *log.Logger
	isErrorLog bool
)

type SQLBuilder struct {
	driver_type string

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
	values [][]interface{}

	// for update
	sets map[string]interface{}
}

func NewSQLBuilder() (b *SQLBuilder) {
	errLog = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	isErrorLog = false

	b = &SQLBuilder{
		driver_type: "mysql",
	}
	b.ClearBuilder()

	return
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
	switch t {
	case "mysql":
		fallthrough
	case "mssql":
		fallthrough
	case "oracle":
		fallthrough
	case "postgresql":
		sb.driver_type = t
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
	sb.values = make([][]interface{}, 0)
	sb.sets = make(map[string]interface{}, 0)
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
	return sb.driver_type == "mysql"
}

func (sb *SQLBuilder) IsMssql() bool {
	return sb.driver_type == "mssql"
}

func (sb *SQLBuilder) IsOracle() bool {
	return sb.driver_type == "oracle"
}

func (sb *SQLBuilder) IsPostgresql() bool {
	return sb.driver_type == "postgresql"
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
	return sb.IsMysql() && sb.limit != ""
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

	sql := "SELECT" + lib.Iif(sb.IsDistinct(), " DISTINCT", "").(string)

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
		sql += " HAVING BY " + sb.havings
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
	for k, v := range sb.sets {
		switch v.(type) {
		case string:
			setStr += fmt.Sprintf("%s='%v',", k, v)
		default:
			setStr += fmt.Sprintf("%s=%v,", k, v)
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

	sql += "(" + strings.Join(sb.fields, ",") + ") VALUES "

	vals := ""
	for _, l := range sb.values {
		ls := "("
		for _, v := range l {
			switch v.(type) {
			case string:
				vals += fmt.Sprintf("'%v',", v)
			default:
				vals += fmt.Sprintf("%v,", v)
			}
		}
		ls = strings.Trim(ls, ",") + "),"
	}
	sql += strings.Trim(vals, ",")
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
	if !sb.IsMysql() {
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
		sb.selects = append(sb.selects, v)
	}

	return sb
}

func (sb *SQLBuilder) From(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support tables")
	}

	for _, v := range s {
		sb.froms = append(sb.froms, v)
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

	sb.froms = append(sb.froms, s)

	return sb
}

func (sb *SQLBuilder) Where(s string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support conditions")
	}
	if !sb.IsHasWheres() {
		sb.wheres = append(sb.wheres, s)
	} else {
		s = "AND " + s
		sb.wheres = append(sb.wheres, s)
	}

	return sb
}

func (sb *SQLBuilder) WhereAnd(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support conditions")
	}

	for _, v := range s {
		if !sb.IsHasWheres() {
			sb.wheres = append(sb.wheres, v)
		} else {
			v = "AND " + v
			sb.wheres = append(sb.wheres, v)
		}
	}

	return sb
}

func (sb *SQLBuilder) WhereOr(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support conditions")
	}

	for _, v := range s {
		if !sb.IsHasWheres() {
			sb.wheres = append(sb.wheres, v)
		} else {
			v = "OR " + v
			sb.wheres = append(sb.wheres, v)
		}
	}

	return sb
}

func (sb *SQLBuilder) Join(s string, c string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	if c != "" {
		s += " ON " + c
	}
	sb.joins = append(sb.joins, "JOIN "+s)

	return sb
}

func (sb *SQLBuilder) InnerJoin(s string, c string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	if c != "" {
		s += " ON " + c
	}
	sb.joins = append(sb.joins, "INNER JOIN "+s)

	return sb
}

func (sb *SQLBuilder) LeftJoin(s string, c string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	if c != "" {
		s += " ON " + c
	}
	sb.joins = append(sb.joins, "LEFT JOIN "+s)

	return sb
}

func (sb *SQLBuilder) RightJoin(s string, c string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	if c != "" {
		s += " ON " + c
	}
	sb.joins = append(sb.joins, "RIGHT JOIN "+s)

	return sb
}

func (sb *SQLBuilder) FullJoin(s string, c string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	if c != "" {
		s += " ON " + c
	}
	sb.joins = append(sb.joins, "FULL OUTER JOIN "+s)

	return sb
}

func (sb *SQLBuilder) GroupBy(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support group fileds")
	}

	for _, v := range s {
		sb.groups = append(sb.groups, v)
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

	sb.orders = append(sb.orders, strings.Join(s, ",")+" ASC")

	return sb
}

func (sb *SQLBuilder) OrderByDesc(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support order fileds")
	}

	sb.orders = append(sb.orders, strings.Join(s, ",")+" DESC")

	return sb
}

func (sb *SQLBuilder) Having(s string) *SQLBuilder {
	if s == "" || !sb.IsHasGroups() {
		sb.PanicOrErrorLog("must be support having condition or set group by first")
	}

	sb.havings = s

	return sb
}

func (sb *SQLBuilder) Set(s map[string]interface{}) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support set fileds : values")
	}

	for k, v := range s {
		sb.sets[k] = v
	}

	return sb
}

func (sb *SQLBuilder) Into(s string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support table")
	}

	sb.into = s

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
		sb.fields = append(sb.fields, v)
	}

	return sb
}

func (sb *SQLBuilder) Values(s ...[]interface{}) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support fileds")
	}

	fieldCnt := sb.GetFieldsCount()
	for _, v := range s {
		if len(v) == fieldCnt {
			sb.values = append(sb.values, v)
		} else {
			sb.PanicOrErrorLog("values count not equal fileds count")
		}
	}

	return sb
}
