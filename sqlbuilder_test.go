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
			name: "case 1",
			fn: func(sb *SQLBuilder) {
				sb.Set([]struct {
					k string
					v interface{}
				}{
					{"foo", 1},
					{"bar", `"2"`},
					{`te"st`, true},
				}).From("user").Where("abc=1").BuildUpdateSQL()
			},
			wantSql: `UPDATE user SET foo=1,bar='\\"2\\"',te\\"st=true WHERE abc=1`,
		},
		{
			name: "case 2",
			fn: func(sb *SQLBuilder) {
				sb.Select("Host", "User", "Select_priv").From("user").Where("User='root'").Limit(1).BuildSelectSQL()
			},
			wantSql: `SELECT Host,User,Select_priv FROM user WHERE User='root'`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := NewSQLBuilder()
			tt.fn(sb)
			if gotSql := sb.BuildedSQL(); gotSql != tt.wantSql {
				t.Errorf("SQLBuilder.BuildedSQL() = %v, want %v", gotSql, tt.wantSql)
			}
		})
	}
}
