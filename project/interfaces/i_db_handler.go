package interfaces

//IDbHandler .
type IDbHandler interface {
	Query(query string, arg interface{}) (IRow, error)
}

//IRow .
type IRow interface {
	StructScan(dest interface{}) error
	Next() bool
	Close() error
	Err() error
}
