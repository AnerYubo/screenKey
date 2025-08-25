package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB 是sqlite操作的封装
type SQLiteDB struct {
	db   *sql.DB
	lock sync.Mutex
}

// NewSQLiteDB 打开或创建一个sqlite数据库文件
func NewSQLiteDB(filepath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	// 设置连接池最大打开连接数和空闲连接数
	db.SetMaxOpenConns(1) // sqlite推荐单线程操作，防止数据库锁死
	db.SetMaxIdleConns(1)

	return &SQLiteDB{
		db: db,
	}, nil
}

// Close 关闭数据库连接
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// Exec 执行增删改，返回影响行数
func (s *SQLiteDB) Exec(query string, args ...interface{}) (int64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// Query 执行查询，返回sql.Rows，需要调用方负责关闭rows
func (s *SQLiteDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.db.Query(query, args...)
}

// QueryRow 执行单行查询，返回sql.Row
func (s *SQLiteDB) QueryRow(query string, args ...interface{}) *sql.Row {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.db.QueryRow(query, args...)
}

// CreateTable 动态创建表，fields 格式示例：
//
//	map[string]string{
//	  "id": "INTEGER PRIMARY KEY AUTOINCREMENT",
//	  "name": "TEXT NOT NULL",
//	  "age": "INTEGER",
//	}
func (s *SQLiteDB) CreateTable(tableName string, fields map[string]string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var columns []string
	for k, v := range fields {
		columns = append(columns, fmt.Sprintf("%s %s", k, v))
	}
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s);", tableName, strings.Join(columns, ", "))
	_, err := s.db.Exec(query)
	return err
}

// Insert 插入数据，fields是字段名和对应值的map
func (s *SQLiteDB) Insert(tableName string, fields map[string]interface{}) (int64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var keys []string
	var placeholders []string
	var values []interface{}

	for k, v := range fields {
		keys = append(keys, k)
		placeholders = append(placeholders, "?")
		values = append(values, v)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, strings.Join(keys, ","), strings.Join(placeholders, ","))
	result, err := s.db.Exec(query, values...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAll 查询所有结果，返回[]map[string]interface{}，每个map是字段名对应值
func (s *SQLiteDB) GetAll(query string, args ...interface{}) ([]map[string]interface{}, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columns[i]

			// sqlite驱动会返回[]byte，转换成string更友好
			if b, ok := val.([]byte); ok {
				m[colName] = string(b)
			} else {
				m[colName] = val
			}
		}

		results = append(results, m)
	}

	return results, nil
}
