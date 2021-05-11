package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	shorturlFieldNames          = builderx.RawFieldNames(&Shorturl{})
	shorturlRows                = strings.Join(shorturlFieldNames, ",")
	shorturlRowsExpectAutoSet   = strings.Join(stringx.Remove(shorturlFieldNames, "`create_time`", "`update_time`"), ",")
	shorturlRowsWithPlaceHolder = strings.Join(stringx.Remove(shorturlFieldNames, "`shorten`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheShorturlShortenPrefix = "cache#shorturl#shorten#"
)

type (
	ShorturlModel interface {
		Insert(data Shorturl) (sql.Result, error)
		FindOne(shorten string) (*Shorturl, error)
		Update(data Shorturl) error
		Delete(shorten string) error
	}

	defaultShorturlModel struct {
		sqlc.CachedConn
		table string
	}

	Shorturl struct {
		Shorten string `db:"shorten"` // shorten key
		Url     string `db:"url"`     // original url
	}
)

func NewShorturlModel(conn sqlx.SqlConn, c cache.CacheConf) ShorturlModel {
	return &defaultShorturlModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`shorturl`",
	}
}

func (m *defaultShorturlModel) Insert(data Shorturl) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, shorturlRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.Shorten, data.Url)

	return ret, err
}

func (m *defaultShorturlModel) FindOne(shorten string) (*Shorturl, error) {
	shorturlShortenKey := fmt.Sprintf("%s%v", cacheShorturlShortenPrefix, shorten)
	var resp Shorturl
	err := m.QueryRow(&resp, shorturlShortenKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `shorten` = ? limit 1", shorturlRows, m.table)
		return conn.QueryRow(v, query, shorten)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultShorturlModel) Update(data Shorturl) error {
	shorturlShortenKey := fmt.Sprintf("%s%v", cacheShorturlShortenPrefix, data.Shorten)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `shorten` = ?", m.table, shorturlRowsWithPlaceHolder)
		return conn.Exec(query, data.Url, data.Shorten)
	}, shorturlShortenKey)
	return err
}

func (m *defaultShorturlModel) Delete(shorten string) error {

	shorturlShortenKey := fmt.Sprintf("%s%v", cacheShorturlShortenPrefix, shorten)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `shorten` = ?", m.table)
		return conn.Exec(query, shorten)
	}, shorturlShortenKey)
	return err
}

func (m *defaultShorturlModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheShorturlShortenPrefix, primary)
}

func (m *defaultShorturlModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `shorten` = ? limit 1", shorturlRows, m.table)
	return conn.QueryRow(v, query, primary)
}
