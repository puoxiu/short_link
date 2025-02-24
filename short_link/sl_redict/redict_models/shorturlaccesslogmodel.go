package redict_models

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ShortUrlAccessLogModel = (*customShortUrlAccessLogModel)(nil)

type (
	// ShortUrlAccessLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customShortUrlAccessLogModel.
	ShortUrlAccessLogModel interface {
		shortUrlAccessLogModel
	}

	customShortUrlAccessLogModel struct {
		*defaultShortUrlAccessLogModel
	}
)

// NewShortUrlAccessLogModel returns a model for the database table.
func NewShortUrlAccessLogModel(conn sqlx.SqlConn) ShortUrlAccessLogModel {
	return &customShortUrlAccessLogModel{
		defaultShortUrlAccessLogModel: newShortUrlAccessLogModel(conn),
	}
}
