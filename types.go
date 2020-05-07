// Author :		Eric<eehsiao@gmail.com>

package sqlbuilder

// SQLBuilder is the main struct type
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
	values [][]interface{}

	// for update
	sets []Set
}

// SQLVar can that you sql internal function via NewSQLVar()
type SQLVar struct {
	VarS string
}

// Set is a type that define the update set struct
type Set struct {
	K string
	V interface{}
}

// SubCond struct for join condition
type SubCond struct {
	c bool // true is and else or
	s string
	o string
	v interface{}
}
