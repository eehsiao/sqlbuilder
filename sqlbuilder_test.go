// Author :		Eric<eehsiao@gmail.com>

package sqlbuilder

import (
	"testing"
)

func TestSQLBuilder_BuildedSQL(t *testing.T) {
	tests := []struct {
		name    string
		fn      func(sb *SQLBuilder)
		wantSQL string
	}{
		{
			name: "case 1 : UPDATE",
			fn: func(sb *SQLBuilder) {
				sb.Set([]Set{{"foo", 1}, {"bar", "\"2\""}, {"te\"st", true}, {"testNil", nil}}).
					From("user").Where("abc", "=", 1).
					WhereOr("def", "=", true).
					WhereAnd("ghi", "like", "%ghi%").
					WhereAnd("jkl", "is", nil).
					WhereAnd("mno", "is not", nil).
					BuildUpdateSQL()
			},
			wantSQL: "UPDATE user SET foo=1,bar='\\\"2\\\"',te\\\"st=true,testNil=NULL WHERE abc = 1 OR def = true AND ghi like '%ghi%' AND jkl is NULL AND mno is not NULL",
		},
		{
			name: "case 2 : JOIN",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").
					From("user").Join("company").
					JoinOn("priv", "abc", "=", 1).
					Limit(1).
					BuildSelectSQL()
			},
			wantSQL: `SELECT Host,User,Select_priv FROM user JOIN company JOIN priv ON abc = 1 LIMIT 1`,
		},
		{
			name: "case 3 : INNER JOIN, LEFT JOIN",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").
					From("user").
					InnerJoin("company").
					InnerJoinOn("priv", "abc", "=", 1).
					LeftJoin("company").
					LeftJoinOn("priv", "abc", "=", 1).
					Limit(1).
					BuildSelectSQL()
			},
			wantSQL: `SELECT Host,User,Select_priv FROM user INNER JOIN company INNER JOIN priv ON abc = 1 LEFT JOIN company LEFT JOIN priv ON abc = 1 LIMIT 1`,
		},
		{
			name: "case 4 : RIGHT JOIN, FULL JOIN",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").
					From("user").
					RightJoin("company a").
					RightJoinOn("priv b", "a.abc", "=", Var("b.abc")).
					FullJoin("comp c").
					FullJoinOn("pri d", "c.def", "=", Var("d.def")).
					Limit(1).
					BuildSelectSQL()
			},
			wantSQL: `SELECT Host,User,Select_priv FROM user RIGHT JOIN company a RIGHT JOIN priv b ON a.abc = b.abc FULL JOIN comp c FULL JOIN pri d ON c.def = d.def LIMIT 1`,
		},
		{
			name: "case 4 : GroupBy OrderBy Having",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").
					From("user").
					OrderBy("Host").
					OrderByAsc("User").
					OrderByDesc("Select_priv").
					GroupBy("Host", "User", "Select_priv").
					Having("count(Host)", ">", 1).
					Havings(On("count(User)", ">", 2)).
					BuildSelectSQL()
			},
			wantSQL: `SELECT Host,User,Select_priv FROM user ORDER BY Host ASC,User ASC,Select_priv DESC GROUP BY Host,User,Select_priv HAVING count(Host) > 1 AND count(User) > 2`,
		},
		{
			name: "case 5 : DELETE",
			fn: func(sb *SQLBuilder) {
				sb.From("user").
					BuildDeleteSQL()
			},
			wantSQL: `DELETE FROM user`,
		},
		{
			name: "case 6 : INSERT",
			fn: func(sb *SQLBuilder) {
				sb.Fields("Host", "User", "Select_priv", "testNil", "testDt").
					Values(1, "\"2", true, nil,
						Var("current_timestamp")).
					Into("user").
					BuildInsertSQL()
			},
			wantSQL: `INSERT INTO user (Host,User,Select_priv,testNil,testDt) VALUES (1,'\"2',true,NULL,current_timestamp)`,
		},
		{
			name: "case 7 : Bulk INSERT",
			fn: func(sb *SQLBuilder) {
				sb.Fields("testDt", "Host", "User", "Select_priv", "testNil").
					Values(Var("current_timestamp"), 1, "\"2", true, nil).
					Values(Var("datetime('now','localtime')"), 2, "\"22", true, nil).
					Values(Var("current_timestamp"), 3, "\"32", false, nil).
					Into("user").
					BuildBulkInsertSQL()
			},
			wantSQL: `INSERT INTO user (testDt,Host,User,Select_priv,testNil) VALUES (current_timestamp,1,'\"2',true,NULL),(datetime('now','localtime'),2,'\"22',true,NULL),(current_timestamp,3,'\"32',false,NULL)`,
		},
		{
			name: "case 8 : Where",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").
					From("user").
					Where("company", "=", "a").
					WhereStr("company!='b'").
					WhereOrStr("user!='b'").
					BuildSelectSQL()
			},
			wantSQL: `SELECT Host,User,Select_priv FROM user WHERE company = 'a' AND company!='b' OR user!='b'`,
		},
		{
			name: "case 9 : JOIN On multi condition",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").
					From("user a").
					Join("company b").
					JoinOns("priv c",
						On("b.abc", "=", 1),
						OnAnd("b.def", "=", Var("c.def")),
					).
					Limit(1).
					BuildSelectSQL()
			},
			wantSQL: `SELECT Host,User,Select_priv FROM user a JOIN company b JOIN priv c ON b.abc = 1 AND b.def = c.def LIMIT 1`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSQLBuilder("mysql")
			tt.fn(sb)
			if gotSQL := sb.BuildedSQL(); gotSQL != tt.wantSQL {
				t.Errorf("SQLBuilder.BuildedSQL() = %v, want %v", gotSQL, tt.wantSQL)
			}
			sb.Release()
		})
	}
}
