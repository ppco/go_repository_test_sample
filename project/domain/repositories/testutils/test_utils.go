package testutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jmoiron/sqlx"
)

// Setup 初期処理
func Setup(m *testing.M) int {
	shutdown := SetupDBConn()
	defer shutdown()

	createTablesIfNotExist()
	setupFixtures()

	return m.Run()
}

var (
	dbConn   *sqlx.DB
	fixtures *testfixtures.Loader
)

const schemaDirRelativePathFormat = "%s/../../../../db/schema.sql"
const fixturesDirRelativePathFormat = "%s/../../../../db/fixtures"

// SetupDBConn DBへの接続を持っておく
func SetupDBConn() func() {
	db, err := sqlx.Connect("mysql", os.Getenv("UNIT_TEST_MYSQL_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Could not connect to mysql: %s", err)
	}

	dbConn = db

	return func() { dbConn.Close() }
}

// GetTestDBConn プールしてあるテスト用のDBコネクションを返す
func GetTestDBConn() *sqlx.DB {
	if dbConn == nil {
		panic("mysql connection is not initialized yet")
	}

	//各テストデータをtrancateしてinsertする
	if err := fixtures.Load(); err != nil {
		log.Fatalf("fixture load error: %v", err)
	}

	return dbConn
}

func createTablesIfNotExist() {
	_, pwd, _, _ := runtime.Caller(0)
	schemaPath := fmt.Sprintf(schemaDirRelativePathFormat, path.Dir(pwd))
	execSchema(schemaPath)
}

func execSchema(fpath string) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatalf("schema reading error: %v", err)
	}

	queries := strings.Split(string(b), ";")

	// 一時的に外部キー制約を外す
	_, err = dbConn.Exec("SET FOREIGN_KEY_CHECKS=0;")
	if err != nil {
		log.Fatalf("exec schema error: %v", err)
	}

	for _, query := range queries {
		if strings.Contains(query, "CREATE TABLE") || strings.Contains(query, "DROP TABLE") {
			_, err = dbConn.Exec(query)
			if err != nil {
				log.Fatalf("exec schema error: %v, query: %s", err, query)
			}
		}
	}

	_, err = dbConn.Exec("SET FOREIGN_KEY_CHECKS=1;")
	if err != nil {
		log.Fatalf("exec schema error: %v", err)
	}
}

func setupFixtures() {
	_, pwd, _, _ := runtime.Caller(0)

	dir := fmt.Sprintf(fixturesDirRelativePathFormat, path.Dir(pwd))
	fix, err := testfixtures.New(
		testfixtures.Database(dbConn.DB),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory(dir),
	)
	if err != nil {
		log.Fatalf("setup fixtures error: %v", err)
	}
	fixtures = fix
}

// SelectItem 指定クエリで単一の結果を取得(INSERTやUPDATEしたときの確認などに使用してください)
func SelectItem(query string, arg map[string]interface{}, dest interface{}) {
	rows, err := dbConn.NamedQuery(query, arg)
	if err != nil {
		log.Fatalf("SelectItem error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(dest)
		if err != nil {
			log.Fatalf("rows.StructScan error: %v", err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("rows.Err error: %v", err)
	}
}

//SelectItems 指定クエリで複数の結果を取得(INSERTやUPDATEしたときの確認などに使用してください)
func SelectItems(query string, arg map[string]interface{}, target interface{}, dest []interface{}) {
	rows, err := dbConn.NamedQuery(query, arg)
	if err != nil {
		log.Fatalf("SelectItem error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(target)
		if err != nil {
			log.Fatalf("rows.StructScan error: %v", err)
		}
		dest = append(dest, target)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("rows.Err error: %v", err)
	}
}
