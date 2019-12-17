// Author :		Eric<eehsiao@gmail.com>

package sqlbuilder

import (
	"testing"
)

func TestSQLBuilder_BuildedSQL(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(sb *SQLBuilder)
		wantSql string
	}{
		{
			name: "case 1 : UPDATE",
			fn: func(sb *SQLBuilder) {
				sb.Set([]Set{{"foo", 1}, {"bar", "\"2\""}, {"te\"st", true}, {"testNil", nil}}).From("user").Where("abc", "=", 1).WhereOr("def", "=", true).WhereAnd("ghi", "like", "%ghi%").WhereAnd("jkl", "is", nil).WhereAnd("mno", "is not", nil).BuildUpdateSQL()
			},
			wantSql: "UPDATE user SET foo=1,bar='\\\"2\\\"',te\\\"st=true,testNil=NULL WHERE abc = 1 OR def = true AND ghi like '%ghi%' AND jkl is NULL AND mno is not NULL",
		},
		{
			name: "case 2 : JOIN",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").From("user").Join("company").JoinOn("priv", "abc", "=", 1).Limit(1).BuildSelectSQL()
			},
			wantSql: `SELECT Host,User,Select_priv FROM user JOIN company JOIN priv ON abc = 1 LIMIT 1`,
		},
		{
			name: "case 3 : INNER JOIN, LEFT JOIN",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").From("user").InnerJoin("company").InnerJoinOn("priv", "abc", "=", 1).LeftJoin("company").LeftJoinOn("priv", "abc", "=", 1).Limit(1).BuildSelectSQL()
			},
			wantSql: `SELECT Host,User,Select_priv FROM user INNER JOIN company INNER JOIN priv ON abc = 1 LEFT JOIN company LEFT JOIN priv ON abc = 1 LIMIT 1`,
		},
		{
			name: "case 4 : RIGHT JOIN, FULL JOIN",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").From("user").RightJoin("company").RightJoinOn("priv", "abc", "=", 1).FullJoin("company").FullJoinOn("priv", "abc", "=", 1).Limit(1).BuildSelectSQL()
			},
			wantSql: `SELECT Host,User,Select_priv FROM user RIGHT JOIN company RIGHT JOIN priv ON abc = 1 FULL JOIN company FULL JOIN priv ON abc = 1 LIMIT 1`,
		},
		{
			name: "case 4 : GroupBy OrderBy Having",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").From("user").OrderBy("Host").OrderByAsc("User").OrderByDesc("Select_priv").GroupBy("Host", "User", "Select_priv").Having("count(Host)", ">", 1).BuildSelectSQL()
			},
			wantSql: `SELECT Host,User,Select_priv FROM user ORDER BY Host ASC,User ASC,Select_priv DESC GROUP BY Host,User,Select_priv HAVING count(Host) > 1`,
		},
		{
			name: "case 5 : DELETE",
			fn: func(sb *SQLBuilder) {
				sb.From("user").BuildDeleteSQL()
			},
			wantSql: `DELETE FROM user`,
		},
		{
			name: "case 6 : INSERT",
			fn: func(sb *SQLBuilder) {
				sb.Fields("Host", "User", "Select_priv", "testNil").Values(1, "\"2", true, nil).Into("user").BuildInsertSQL()
			},
			wantSql: `INSERT INTO user (Host,User,Select_priv,testNil) VALUES (1,'\"2',true,NULL)`,
		},
		{
			name: "case 7 : Bulk INSERT",
			fn: func(sb *SQLBuilder) {
				sb.Fields("Host", "User", "Select_priv", "testNil").Values(1, "\"2", true, nil).Values(2, "\"22", true, nil)
				sb.Values(3, "\"32", false, nil).Into("user").BuildBulkInsertSQL()
			},
			wantSql: `INSERT INTO user (Host,User,Select_priv,testNil) VALUES (1,'\"2',true,NULL),(2,'\"22',true,NULL),(3,'\"32',false,NULL)`,
		},
		{
			name: "case 8 : Where",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").From("user").Where("company", "=", "a").WhereStr("company!='b'").WhereOrStr("user!='b'").BuildSelectSQL()
			},
			wantSql: `SELECT Host,User,Select_priv FROM user WHERE company = 'a' AND company!='b' OR user!='b'`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSQLBuilder("mysql")
			tt.fn(sb)
			if gotSql := sb.BuildedSQL(); gotSql != tt.wantSql {
				t.Errorf("SQLBuilder.BuildedSQL() = %v, want %v", gotSql, tt.wantSql)
			}
			sb.Release()
		})
	}
}
