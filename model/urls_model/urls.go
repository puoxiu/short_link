package urlsmodel

import "time"


type UrlsModel struct {
    ID         int       `json:"id" gorm:"primaryKey;autoIncrement"`       // 主键
    LongLink   string    `json:"longLink" gorm:"type:varchar(255);not null"`  // 长链接
    ShortLink  string    `json:"shortLink" gorm:"type:varchar(10);not null;unique"` // 短链接
    IsCustom   bool      `json:"isCustom" gorm:"type:tinyint(1);default:0"`  // 是否自定义
    ExpireTime time.Time `json:"expireTime" gorm:"type:datetime;not null"`   // 过期时间
    CreateTime time.Time `json:"createTime" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"` // 创建时间
}
