// Author :		Eric<eehsiao@gmail.com>

package sqlbuilder

import (
	"fmt"
	"strconv"
	"strings"
)

// NewSQLVar gen you want string int builder
func NewSQLVar(s string) SQLVar {
	return Var(s)
}

// Var same as NewSQLVar
func Var(s string) SQLVar {
	return SQLVar{VarS: s}
}

// On return a sub condition for join or having
// same as OnAnd
func On(s string, o string, v interface{}) SubCond {
	return OnAnd(s, o, v)
}

// OnAnd its will be `and` sub condition
func OnAnd(s string, o string, v interface{}) SubCond {
	return SubCond{c: true, s: s, o: o, v: v}
}

// OnOr its will be `or` sub condition
func OnOr(s string, o string, v interface{}) SubCond {
	return SubCond{c: false, s: s, o: o, v: v}
}

// Distinct set builder for `distinct`
func (sb *SQLBuilder) Distinct(b bool) *SQLBuilder {
	sb.distinct = b

	return sb
}

// Limit set builder for `limit`
// only for Mysql, SQLite
func (sb *SQLBuilder) Limit(i ...int) *SQLBuilder {
	if !(sb.IsMysql() || sb.IsSQLite()) {
		sb.PanicOrErrorLog("limit only support mysql or sqlite")
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

// Top set builder for `top`
// only for Mssql
func (sb *SQLBuilder) Top(i int) *SQLBuilder {
	if !sb.IsMssql() {
		sb.PanicOrErrorLog("top only support mssql")
	}
	if i <= 0 {
		sb.PanicOrErrorLog("must have >=1 value for top")
	}

	sb.top = strconv.Itoa(i)

	return sb
}

// Select set builder for `select`
// params must lest one or more
// ex :
// ```
// Select('fieldA', 'fieldB', 'fieldC')
// ```
func (sb *SQLBuilder) Select(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support fileds")
	}

	for _, v := range s {
		sb.selects = append(sb.selects, EscapeStr(v, sb.IsMysql()))
	}

	return sb
}

// From set builder for `from`
// params must lest one or more
// ex :
// ```
// From('tblA', 'tblB', 'tblC')
// ```
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

// FromOne set builder for `from`
// just one param for one table
// ex :
// ```
// FromOne('tblA')
// ```
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

// WhereStr is equal WhereAndStr
// just one param for add a where string
// if this not first time use, its will be `and` condition
// ex :
// ```
// WhereStr('fieldA = 0')
// ```
func (sb *SQLBuilder) WhereStr(s string) *SQLBuilder {
	return sb.WhereAndStr(s)
}

// WhereAndStr same as WhereStr
func (sb *SQLBuilder) WhereAndStr(s string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support conditions")
	}

	if sb.IsHasWheres() {
		sb.wheres = append(sb.wheres, fmt.Sprintf("AND %s", s))
	} else {
		sb.wheres = append(sb.wheres, s)
	}

	return sb
}

// WhereOrStr just one param for add a where string
// if this not first time use, its will be `or` condition
// it its first time use, just same as WhereStr
// ex :
// ```
// WhereOrStr('fieldA = 0')
// ```
func (sb *SQLBuilder) WhereOrStr(s string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support conditions")
	}

	if sb.IsHasWheres() {
		sb.wheres = append(sb.wheres, fmt.Sprintf("OR %s", s))
	} else {
		sb.wheres = append(sb.wheres, s)
	}

	return sb
}

// Where is equal WhereAnd
// if this not first time use, its will be `and` condition
// 3 params :
// s is mean filed
// o is a operator, ex : `=`, `>`, ...
// v is a value , it type interface{}
// ex :
// ```
// Where('fieldA', '=', 0)
// ```
func (sb *SQLBuilder) Where(s string, o string, v interface{}) *SQLBuilder {
	return sb.WhereAnd(s, o, v)
}

// WhereAnd same as Where
func (sb *SQLBuilder) WhereAnd(s string, o string, v interface{}) *SQLBuilder {
	if s == "" || o == "" {
		sb.PanicOrErrorLog("must be support conditions")
	}

	c := ""
	if sb.IsHasWheres() {
		c = "AND "
	}
	switch v.(type) {
	case SQLVar:
		sb.wheres = append(sb.wheres, fmt.Sprintf("%s%s %s %s", c, EscapeStr(s, sb.IsMysql()), EscapeStr(o, sb.IsMysql()), EscapeStr(v.(SQLVar).VarS, sb.IsMysql())))
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

// WhereOr if this not first time use, its will be `or` condition
// 3 params :
// s is mean filed
// o is a operator, ex : `=`, `>`, ...
// v is a value , it type interface{}
// ex :
// ```
// WhereOr('fieldA', '=', 0)
// ```
func (sb *SQLBuilder) WhereOr(s string, o string, v interface{}) *SQLBuilder {
	if s == "" || o == "" {
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

func (sb *SQLBuilder) joinOn(p string, t string, j ...SubCond) *SQLBuilder {
	if t == "" || len(j) == 0 {
		sb.PanicOrErrorLog("must be support join table or without condition")
	}
	jStr, c := "", ""

	for _, con := range j {
		if jStr == "" {
			switch con.v.(type) {
			case SQLVar:
				jStr = fmt.Sprintf("%sJOIN %s ON %s %s %s", p, t, EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), EscapeStr(con.v.(SQLVar).VarS, sb.IsMysql()))
			case string:
				jStr = fmt.Sprintf("%sJOIN %s ON %s %s '%s'", p, t, EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), EscapeStr(con.v.(string), sb.IsMysql()))
			default:
				jStr = fmt.Sprintf("%sJOIN %s ON %s %s %v", p, t, EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), con.v)
			}
		} else {
			if con.c {
				c = " AND"
			} else {
				c = " OR"
			}
			switch con.v.(type) {
			case SQLVar:
				jStr += fmt.Sprintf("%s %s %s %s", c, EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), EscapeStr(con.v.(SQLVar).VarS, sb.IsMysql()))
			case string:
				jStr += fmt.Sprintf("%s %s %s '%s'", c, EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), EscapeStr(con.v.(string), sb.IsMysql()))
			default:
				jStr += fmt.Sprintf("%s %s %s %v", c, EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), con.v)
			}
		}

	}

	if jStr != "" {
		sb.joins = append(sb.joins, jStr)
	}

	return sb
}

// Join is a natural join
func (sb *SQLBuilder) Join(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("", j)
}

// JoinOn the join with one condition
func (sb *SQLBuilder) JoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("", j, On(s, o, v))
}

// JoinOns the join with multi condition
func (sb *SQLBuilder) JoinOns(j string, on ...SubCond) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("", j, on...)
}

// InnerJoin the join with natural fileds
func (sb *SQLBuilder) InnerJoin(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("INNER ", j)
}

// InnerJoinOn the join with one condition
func (sb *SQLBuilder) InnerJoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("INNER ", j, On(s, o, v))
}

// InnerJoinOns the join with multi condition
func (sb *SQLBuilder) InnerJoinOns(j string, on ...SubCond) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("INNER ", j, on...)
}

// LeftJoin the join with natural fileds
func (sb *SQLBuilder) LeftJoin(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("LEFT ", j)
}

// LeftJoinOn the join with one condition
func (sb *SQLBuilder) LeftJoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("LEFT ", j, On(s, o, v))
}

// LeftJoinOns the join with multi condition
func (sb *SQLBuilder) LeftJoinOns(j string, on ...SubCond) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("LEFT ", j, on...)
}

// RightJoin the join with natural fileds
func (sb *SQLBuilder) RightJoin(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("RIGHT ", j)
}

// RightJoinOn the join with one condition
func (sb *SQLBuilder) RightJoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("RIGHT ", j, On(s, o, v))
}

// RightJoinOns the join with multi condition
func (sb *SQLBuilder) RightJoinOns(j string, on ...SubCond) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("RIGHT ", j, on...)
}

// FullJoin the join with natural fileds
func (sb *SQLBuilder) FullJoin(j string) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.join("FULL ", j)
}

// FullJoinOn the join with one condition
func (sb *SQLBuilder) FullJoinOn(j string, s string, o string, v interface{}) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("FULL ", j, On(s, o, v))
}

// FullJoinOns the join with multi condition
func (sb *SQLBuilder) FullJoinOns(j string, on ...SubCond) *SQLBuilder {
	if j == "" {
		sb.PanicOrErrorLog("must be support join table")
	}

	return sb.joinOn("FULL ", j, on...)
}

// GroupBy with fileds
// fileds must be same as select fileds
func (sb *SQLBuilder) GroupBy(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support group fileds")
	}

	for _, v := range s {
		sb.groups = append(sb.groups, EscapeStr(v, sb.IsMysql()))
	}

	return sb
}

// OrderBy with fileds
// default Asc
func (sb *SQLBuilder) OrderBy(s ...string) *SQLBuilder {
	return sb.OrderByAsc(s...)
}

// OrderByAsc with fileds
func (sb *SQLBuilder) OrderByAsc(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support order fileds")
	}

	sb.orders = append(sb.orders, EscapeStr(strings.Join(s, ","), sb.IsMysql())+" ASC")

	return sb
}

// OrderByDesc with fileds
func (sb *SQLBuilder) OrderByDesc(s ...string) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support order fileds")
	}

	sb.orders = append(sb.orders, EscapeStr(strings.Join(s, ","), sb.IsMysql())+" DESC")

	return sb
}

func (sb *SQLBuilder) having(h ...SubCond) *SQLBuilder {
	if len(h) == 0 {
		sb.PanicOrErrorLog("without condition")
	}
	if !sb.IsHasGroups() {
		sb.PanicOrErrorLog("must be set group by first")
	}

	c := ""

	for _, con := range h {
		if sb.havings == "" {
			switch con.v.(type) {
			case SQLVar:
				sb.havings = fmt.Sprintf("%s %s %s", EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), EscapeStr(con.v.(SQLVar).VarS, sb.IsMysql()))
			case string:
				sb.havings = fmt.Sprintf("%s %s '%s'", EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), EscapeStr(con.v.(string), sb.IsMysql()))
			default:
				sb.havings = fmt.Sprintf("%s %s %v", EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), con.v)
			}
		} else {
			if con.c {
				c = " AND"
			} else {
				c = " OR"
			}
			switch con.v.(type) {
			case SQLVar:
				sb.havings += fmt.Sprintf("%s %s %s %s", c, EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), EscapeStr(con.v.(SQLVar).VarS, sb.IsMysql()))
			case string:
				sb.havings += fmt.Sprintf("%s %s %s '%s'", c, EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), EscapeStr(con.v.(string), sb.IsMysql()))
			default:
				sb.havings += fmt.Sprintf("%s %s %s %v", c, EscapeStr(con.s, sb.IsMysql()), EscapeStr(con.o, sb.IsMysql()), con.v)
			}
		}

	}

	return sb
}

// Having with one condition
// its will overwrite having section
func (sb *SQLBuilder) Having(s string, o string, v interface{}) *SQLBuilder {
	if s == "" || !sb.IsHasGroups() {
		sb.PanicOrErrorLog("must be support having condition or set group by first")
	}

	return sb.having(On(s, o, v))
}

// Havings with multi conditions
func (sb *SQLBuilder) Havings(h ...SubCond) *SQLBuilder {
	if len(h) == 0 {
		sb.PanicOrErrorLog("without condition")
	}
	if !sb.IsHasGroups() {
		sb.PanicOrErrorLog("must be set group by first")
	}

	return sb.having(h...)
}

// Set with Set{K string, V interface{}} structs
func (sb *SQLBuilder) Set(s []Set) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support set fileds : values")
	}

	for _, set := range s {
		sb.sets = append(sb.sets, set)
	}

	return sb
}

// Into for set insert table
func (sb *SQLBuilder) Into(s string) *SQLBuilder {
	if s == "" {
		sb.PanicOrErrorLog("must be support table")
	}

	sb.into = EscapeStr(s, sb.IsMysql())

	return sb
}

// Fields for set update fields
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

// Values for set update values
func (sb *SQLBuilder) Values(s ...interface{}) *SQLBuilder {
	if len(s) == 0 {
		sb.PanicOrErrorLog("must be support fileds")
	}

	fieldCnt := sb.GetFieldsCount()
	if len(s) == fieldCnt {
		vs := make([]interface{}, 0)
		for _, v := range s {
			vs = append(vs, v)
		}
		sb.values = append(sb.values, vs)
	} else {
		sb.PanicOrErrorLog("values count not equal fileds count")
	}

	return sb
}
