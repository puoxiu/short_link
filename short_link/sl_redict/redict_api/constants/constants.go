package constants

type contextKey string

// 定义上下文键常量
const (
    UserAgentKey contextKey = "User_Agent"
    UserIPKey    contextKey = "User_IP"

    DefaultDBTimeout int64 = 5  // 数据库操作默认超时时间
)