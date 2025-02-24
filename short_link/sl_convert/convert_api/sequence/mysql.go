package sequence

import (
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const sqlReplaceIntoStub = `REPLACE INTO sequence (stub) VALUES ('x')`

type MySQL struct {
	conn sqlx.SqlConn
}

func NewMySQL(dsn string) Sequence {
	return &MySQL{
		conn: sqlx.NewMysql(dsn),
	}
}

// NextNumber 取号器 获取下一个序列号
func (m *MySQL) NextNumber() (seq uint64, err error) {
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlReplaceIntoStub)
	if err != nil {
		logx.Errorw("prepare failed", logx.LogField{Key: "err", Value: err.Error()})
		return
	}
	defer stmt.Close()
	// 执行
	var res sql.Result
	res, err = stmt.Exec()
	if err != nil {
		logx.Errorw("exec failed", logx.LogField{Key: "err", Value: err.Error()})
		return
	}

	// 获取最后插入的id
	var lid int64
	lid, err = res.LastInsertId()
	if err != nil {
		logx.Errorw("last insert id failed", logx.LogField{Key: "err", Value: err.Error()})
		return
	}

	return uint64(lid), nil
}