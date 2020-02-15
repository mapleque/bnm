package server

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type DBConn interface {
	Begin() (DBConn, error)
	Rollback() error
	Commit() error
	Execute(sql string, params ...interface{}) DBExecResult
	Query(sql string, params ...interface{}) DBQueryResult
}

type MysqlConn struct {
	conn *sql.DB
	log  Logger
}

type MysqlTxConn struct {
	conn  *sql.Tx
	log   Logger
	token string
}

func NewDB(dsn string) DBConn {
	log := NewStdLogger("[DB]")
	log.Trace("open db", dsn)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err := conn.Ping(); err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	return &MysqlConn{
		conn: conn,
		log:  log,
	}
}

func (c *MysqlConn) Begin() (DBConn, error) {
	token := GetSessionToken()
	c.log.Trace("start a trancsaction", token)
	conn, err := c.conn.Begin()
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	return &MysqlTxConn{
		conn:  conn,
		log:   c.log,
		token: token,
	}, nil
}

var errMethodUnsupport = errors.New("this method is unsupport")

func (c *MysqlConn) Rollback() error { return errMethodUnsupport }
func (c *MysqlConn) Commit() error   { return errMethodUnsupport }

func (c *MysqlConn) Execute(sql string, params ...interface{}) DBExecResult {
	c.log.Trace(sql, params)
	sr, err := c.conn.Exec(sql, params...)
	if err != nil {
		c.log.Error(err)
	}
	return &MysqlExecResult{
		sr:  sr,
		err: err,
	}
}
func (c *MysqlConn) Query(sql string, params ...interface{}) DBQueryResult {
	c.log.Trace(sql, params)
	ret := &MysqlQueryResult{}
	stmt, err := c.conn.Prepare(sql)
	if err != nil {
		ret.err = err
		c.log.Error(err)
		return ret
	}

	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		c.log.Error(err)
	}
	ret.err = err
	ret.rows = rows
	return ret
}

func (c *MysqlTxConn) Begin() (DBConn, error) { return nil, errMethodUnsupport }

func (c *MysqlTxConn) Rollback() error {
	c.log.Trace(c.token, "rollback")
	return c.conn.Rollback()
}
func (c *MysqlTxConn) Commit() error {
	c.log.Trace(c.token, "commit")
	return c.conn.Commit()
}

func (c *MysqlTxConn) Execute(sql string, params ...interface{}) DBExecResult {
	c.log.Trace(c.token, sql, params)
	sr, err := c.conn.Exec(sql, params...)
	if err != nil {
		c.log.Error(c.token, err)
	}
	return &MysqlExecResult{
		sr:  sr,
		err: err,
	}
}
func (c *MysqlTxConn) Query(sql string, params ...interface{}) DBQueryResult {
	c.log.Trace(c.token, sql, params)
	ret := &MysqlQueryResult{}
	stmt, err := c.conn.Prepare(sql)
	if err != nil {
		c.log.Error(c.token, err)
		ret.err = err
		return ret
	}

	rows, err := stmt.Query(params...)
	if err != nil {
		c.log.Error(c.token, err)
	}
	ret.err = err
	ret.rows = rows
	ret.stmt = stmt
	return ret
}

type DBExecResult interface {
	Error() error
	LastId() (int64, error)
	EffectRows() (int64, error)
}

type MysqlExecResult struct {
	sr  sql.Result
	err error
}

func (r *MysqlExecResult) Error() error {
	return r.err
}

func (r *MysqlExecResult) LastId() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.sr.LastInsertId()
}

func (r *MysqlExecResult) EffectRows() (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.sr.RowsAffected()
}

type DBQueryResult interface {
	Exist() (bool, error)
	Scan(tar ...interface{}) error
	ScanFunc(func(Scanner) error) error
}

type MysqlQueryResult struct {
	err  error
	rows *sql.Rows
	stmt *sql.Stmt
}

func (r *MysqlQueryResult) Exist() (bool, error) {
	if r.err != nil {
		return false, r.err
	}
	if r.stmt != nil {
		defer r.stmt.Close()
	}
	if r.rows != nil {
		defer r.rows.Close()
	}
	return r.rows.Next(), nil
}

func (r *MysqlQueryResult) Scan(tar ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	if r.stmt != nil {
		defer r.stmt.Close()
	}
	if r.rows != nil {
		defer r.rows.Close()
	}
	if r.rows.Next() {
		return r.rows.Scan(tar...)
	}
	return sql.ErrNoRows
}

type Scanner func(...interface{}) error

func (r *MysqlQueryResult) ScanFunc(f func(Scanner) error) error {
	if r.err != nil {
		return r.err
	}
	if r.stmt != nil {
		defer r.stmt.Close()
	}
	if r.rows != nil {
		defer r.rows.Close()
	}
	for r.rows.Next() {
		if err := f(r.rows.Scan); err != nil {
			return err
		}
	}
	return nil
}
