[![GoDoc](https://godoc.org/github.com/eehsiao/sqlbuilder?status.svg)](https://godoc.org/github.com/eehsiao/sqlbuilder)
# sqlbuilder
`sqlbuilder` is a simple sql query string builder

sqlbuilder its recursive struct call, that you can easy to build sql string

ex: dao.Select().From().Join().Where().Limit()
### SqlBuilder functions
* build select :
    * Select(f ...string)
        * f is a fileds list of strings
        * Select("filed1", "filed2", "filed3")
        * Select(lib.Struce4QuerySlice(DaoStructType)...)
            * the library ref : [https://github.com/eehsiao/go-models-lib](https://github.com/eehsiao/go-models-lib)
    * Distinct(b bool)
        * its default in builder is set `false`
    * Top(i int)
        * only support mssql
    * From(t ...string)
        * t is table name
    * Where(c string)
        * c is condition, ex Where("field1=1 and filed2='b'")
    * WhereAnd(c ...string)
    * WhereOr(c ...string)
    * Limit(i ...int)
        * support 2 parms
        * only support mysql
    * Join(t string, c string)
        * t is table name
        * c is condition
    * InnerJoin(t string, c string)
    * LeftJoin(t string, c string)
    * RightJoin(t string, c string)
    * FullJoin(t string, c string)
    * GroupBy(f ...string)
        * f is a fileds list of strings
    * OrderBy(f ...string)
        * f is a fileds list of strings
    * OrderByAsc(f ...string)
    * OrderByDesc(f ...string)
    * Having(s string)
        * s is having condition string
    * BuildSelectSQL()
        * check and build sql string.
        * you can get sql string via `BuildedSQL()`
* build update :
    * Set(s map[string]interface{})
    * FromOne(t string)
        * reset the table for only one
    * BuildUpdateSQL()
        * check and build sql string.
        * you can get sql string via `BuildedSQL()`
* build insert : 
    * Into(t string)
        * set the insert table
    * Fields(f ...string)
        * f is a fileds list of strings
    * Values(v ...[]interface{})
        * v is a values list of `interface{}`
    * BuildInsertSQL()
        * check and build sql string.
        * you can get sql string via `BuildedSQL()`
* build delete :
    * BuildDeleteSQL()
        * check and build sql string.
        * you can get sql string via `BuildedSQL()`
* common :
    * ClearBuilder()
        * reset builder
    * BuildedSQL()
        * return the builded sql string, if build success.
    * SetDbName(s string)
    * SetTbName(s string)
    * SwitchPanicToErrorLog(b bool)
    * PanicOrErrorLog(s string)