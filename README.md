# sqlbuilder
`sqlbuilder` is a sample sql query string builder

sqlbuilder its recursive call function, that you can easy to build sql string

ex: dao.Select().From().Join().Where().Limit()
### SqlBuilder functions
* build select :
    * Select(s ...string)
    * Distinct(b bool)
    * Top(i int)
    * From(s ...string)
    * Where(s string)
    * WhereAnd(s ...string)
    * WhereOr(s ...string)
    * Join(s string, c string)
    * InnerJoin(s string, c string)
    * LeftJoin(s string, c string)
    * RightJoin(s string, c string)
    * FullJoin(s string, c string)
    * GroupBy(s ...string)
    * OrderBy(s ...string)
    * OrderByAsc(s ...string)
    * OrderByDesc(s ...string)
    * Having(s string)
    * BuildSelectSQL()
* build update :
    * Set(s map[string]interface{})
    * FromOne(s string)
    * BuildUpdateSQL()
* build insert : 
    * Into(s string)
    * Fields(s ...string)
    * Values(s ...[]interface{})
    * BuildInsertSQL()
* build delete :
    * BuildDeleteSQL()
* common :
    * ClearBuilder()
    * BuildedSQL()
    * SetDbName(s string)
    * SetTbName(s string)
    * SwitchPanicToErrorLog(b bool)
    * PanicOrErrorLog(s string)