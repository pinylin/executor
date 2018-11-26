package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"sync"
	"time"
)

// 引用包go-sql-driver
// go get "github.com/go-sql-driver/mysql"

// TODO: Fix Bugs

var db *sql.DB

var dbNilError error = errors.New("db is nil")

func checkError(err error) {
	if err != nil {
		fmt.Println("db fail: ", err)
		panic(err)
	}
}

//var once sync.Once

//func New() {
//	if db == nil {
//		once.Do(openDB)
//	} else {
//		fmt.Println("db already exists")
//	}
//}
//
//func GetDB() *sql.DB {
//	if db == nil {
//		openDB()
//	}
//	return db
//}

var mu sync.Mutex

func openDB(string) {
	mu.Lock()
	defer mu.Unlock()
	if db != nil {
		return
	}
	var err error

	//"user=pqgotest dbname=pqgotest sslmode=verify-full"

	db, err = sql.Open("postgres", os.Getenv("POSTGRES"))
	if err != nil {
		errMsg := fmt.Sprintf("数据库连接失败1: %v", err)
		fmt.Printf(errMsg)
		panic(err)
	}
	db.SetConnMaxLifetime(time.Second * 10)
	if err = db.Ping(); err != nil {
		errMsg := fmt.Sprintf("数据库连接失败2: %v", err)
		fmt.Printf(errMsg)
		panic(err)
	}
}

//插入
func Insert(db *sql.DB, sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	stmtIns, err := db.Prepare(sqlstr)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()
	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

//修改和删除
func Exec(db *sql.DB, sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	stmtIns, err := db.Prepare(sqlstr)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()

	result, err := stmtIns.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//取一行数据，注意这类取出来的结果都是string
func FetchRow(db *sql.DB, sqlstr string, args ...interface{}) (*map[string]string, error) {

	if db == nil {
		return nil, dbNilError
	}

	stmtOut, err := db.Prepare(sqlstr)
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	ret := make(map[string]string, len(scanArgs))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var value string

		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			ret[columns[i]] = value
		}
		break //get the first row only
	}
	return &ret, nil
}

//取多行，<span style="font-family: Arial, Helvetica, sans-serif;">注意这类取出来的结果都是string </span>
func FetchRows(db *sql.DB, sqlstr string, args ...interface{}) (*[]map[string]string, error) {
	if db == nil {
		return nil, dbNilError
	}
	stmtOut, err := db.Prepare(sqlstr)
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	ret := make([]map[string]string, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var value string
		vmap := make(map[string]string, len(scanArgs))
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret, vmap)
	}
	return &ret, nil
}

//插入
func InsertD(sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	result, err := db.Exec(sqlstr, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

//修改和删除
func ExecD(sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	result, err := db.Exec(sqlstr, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// 插入
func ExecInsertD(sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	result, err := db.Exec(sqlstr, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

//插入(有事务处理) 调用时不用Rollback
func InsertTx(tx *sql.Tx, sqlstr string, args ...interface{}) (int64, error) {

	if db == nil {
		return 0, dbNilError
	}
	result, err := tx.Exec(sqlstr, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	n, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return n, nil
}

//修改和删除(有事务处理) 调用时不Rollback
func ExecTx(tx *sql.Tx, sqlstr string, args ...interface{}) (int64, error) {
	if db == nil {
		return 0, dbNilError
	}
	result, err := tx.Exec(sqlstr, args...)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return n, nil
}

//查询一行数据
func FetchRowD(sqlstr string, args ...interface{}) (map[string]string, error) {
	if db == nil {
		return nil, dbNilError
	}
	rows, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret, err := fetchData(rows)
	if err != nil {
		return nil, err
	}
	if len(ret) > 0 {
		return ret[0], nil
	}
	return map[string]string{}, nil
}

//查询多行数据
func FetchRowsD(sqlstr string, args ...interface{}) ([]map[string]string, error) {
	if db == nil {
		return nil, dbNilError
	}
	rows, err := db.Query(sqlstr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret, err := fetchData(rows)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func fetchData(rows *sql.Rows) ([]map[string]string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	ret := make([]map[string]string, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var value string
		vmap := make(map[string]string, len(scanArgs))
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret, vmap)
	}
	return ret, nil
}

//批量删除了
//func BatchInsertD(sqlstr string, args ...interface{}) (int64, error) {
//	if db == nil {
//		return 0, dbNilError
//	}
//	stmtOut, err := db.Prepare(sqlstr)
//	if err != nil {
//		return 0, err
//	}
//	defer stmtOut.Close()
//
//	_, err = stmtOut.Exec(args...)
//	if err != nil {
//		fmt.Println("=======", err)
//	}
//
//	return 0, nil
//}
//
//func NewTest() {
//	if db == nil {
//		openDB()
//	} else {
//		fmt.Println("db no open")
//	}
//}
