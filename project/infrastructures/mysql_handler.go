package infrastructures

import (
	"project/interfaces"

	//MySQL
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//MySQLHandler .
type MySQLHandler struct {
	Conn *sqlx.DB
}

//MySQLRow .
type MySQLRow struct {
	Rows *sqlx.Rows
}

//Query .
func (handler *MySQLHandler) Query(query string, arg interface{}) (interfaces.IRow, error) {
	rows, err := handler.Conn.NamedQuery(query, arg)
	if err != nil {
		return new(MySQLRow), err
	}
	row := new(MySQLRow)
	row.Rows = rows

	return row, nil
}

//StructScan .
func (r MySQLRow) StructScan(dest interface{}) error {
	err := r.Rows.StructScan(dest)
	if err != nil {
		return err
	}

	return nil
}

//Next .
func (r MySQLRow) Next() bool {
	return r.Rows.Next()
}

//Close .
func (r MySQLRow) Close() error {
	err := r.Rows.Close()
	if err != nil {
		return err
	}

	return nil
}

//Err .
func (r MySQLRow) Err() error {
	err := r.Rows.Err()
	if err != nil {
		return err
	}

	return nil
}
