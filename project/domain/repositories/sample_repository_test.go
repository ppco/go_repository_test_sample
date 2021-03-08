package repositories

import (
	"errors"
	"project/domain/model"
	"project/domain/repositories/testutils"
	"project/infrastructures"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func Test_Sample実態DBテスト(t *testing.T) {
	db := testutils.GetTestDBConn()
	mySQLHandler := &infrastructures.MySQLHandler{}
	mySQLHandler.Conn = db
	repo := SampleRepository{IDbHandler: mySQLHandler}

	expectedSample := &model.Sample{
		ID:   1,
		Code: "hoge",
		Name: "ほげ",
	}
	actualSample, err := repo.FindSample(int64(1))

	if !reflect.DeepEqual(expectedSample, actualSample) {
		t.Errorf("%sで取得結果が異なる expected:%v actual:%v", t.Name(), expectedSample, actualSample)
	}

	if err != nil {
		t.Errorf("%sでerrがnilじゃない err:%v", t.Name(), err)
	}
}

func Test_振る舞いテスト_コネクションエラー時のテスト(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mockHandler := &infrastructures.MySQLHandler{}
	mockHandler.Conn = sqlxDB

	expectedErr := errors.New("sql: database is closed")

	sqlxDB.Close()
	sampleRepository := SampleRepository{IDbHandler: mockHandler}

	_, actualErr := sampleRepository.FindSample(int64(1))
	if !reflect.DeepEqual(expectedErr.Error(), actualErr.Error()) {
		t.Errorf("%sで期待結果相違 expectedErr:%v actualErr:%v", t.Name(), expectedErr, actualErr)
	}
}

func Test_振る舞いテスト_Scanエラー時のテスト(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mockHandler := &infrastructures.MySQLHandler{}
	mockHandler.Conn = sqlxDB

	expectedErr := errors.New("missing destination name x in *model.Sample")

	/*
		SELECT・FROMクエリは全て記述
		WHERE等に関しては、ここでは考慮しないので省略してもOK
		想定外のカラム「x」が指定されていた時を想定
	*/
	mock.ExpectQuery(regexp.QuoteMeta(`
SELECT
	 id
	,code
	,name
FROM
	sample
`)).WillReturnRows(sqlmock.NewRows([]string{"x"}).AddRow(1))

	sampleRepository := SampleRepository{IDbHandler: mockHandler}

	_, actualErr := sampleRepository.FindSample(int64(1))

	if !reflect.DeepEqual(expectedErr.Error(), actualErr.Error()) {
		t.Errorf("%sで期待結果相違 expectedErr:%v actualErr:%v", t.Name(), expectedErr, actualErr)
	}
}

func Test_振る舞いテスト_行取得時エラー時のテスト(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mockHandler := &infrastructures.MySQLHandler{}
	mockHandler.Conn = sqlxDB

	expectedErr := errors.New("rows error")

	/*
		クエリは評価しないので記述しない
		Rows自体からエラーが発生した場合を想定
	*/
	mock.ExpectQuery(".*").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).RowError(0, errors.New("rows error")))

	sampleRepository := SampleRepository{IDbHandler: mockHandler}

	_, actualErr := sampleRepository.FindSample(int64(1))

	if !reflect.DeepEqual(expectedErr.Error(), actualErr.Error()) {
		t.Errorf("%sで期待結果相違 expectedErr:%v actualErr:%v", t.Name(), expectedErr, actualErr)
	}
}
