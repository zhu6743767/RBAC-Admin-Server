package examples

// 使用示例：如何在项目中使用新的日志系统

// 示例1：记录数据库连接错误
func ExampleDatabaseError() {
	// 假设在数据库初始化时
	// errors.LogDatabaseConnectionError("192.168.10.199", 3306, "root", err, map[string]interface{}{
	// 	"request_id": "req_12345",
	// 	"user_id":    "system",
	// })
}

// 示例2：记录用户注册操作
func ExampleUserRegister() {
	// opLogger := audit.NewOperationLogger()
	// opLogger.LogUserRegister(
	// 	"user_123",
	// 	"zhangsan",
	// 	"192.168.1.100",
	// 	audit.Success,
	// 	map[string]interface{}{
	// 		"email": "zhangsan@example.com",
	// 		"phone": "13800138000",
	// 	},
	// )
}

// 示例3：记录用户登录操作
func ExampleUserLogin() {
	// opLogger := audit.NewOperationLogger()
	// opLogger.LogUserLogin(
	// 	"user_123",
	// 	"zhangsan",
	// 	"192.168.1.100",
	// 	audit.Success,
	// 	map[string]interface{}{
	// 		"login_method": "password",
	// 		"user_agent":    "Mozilla/5.0...",
	// 	},
	// )
}

// 示例4：记录数据修改操作
func ExampleDataUpdate() {
	// opLogger := audit.NewOperationLogger()
	// before := map[string]interface{}{
	// 	"username": "old_name",
	// 	"email":    "old@example.com",
	// }
	// after := map[string]interface{}{
	// 	"username": "new_name",
	// 	"email":    "new@example.com",
	// }
	// opLogger.LogDataOperation(
	// 	audit.ProfileUpdate,
	// 	"user_123",
	// 	"/api/users/profile",
	// 	before,
	// 	after,
	// 	audit.Success,
	// 	nil,
	// )
}

// 示例5：记录网络错误
func ExampleNetworkError() {
	// errors.LogNetworkConnectionError(
	// 	"Redis",
	// 	"redis://localhost:6379",
	// 	err,
	// 	map[string]interface{}{
	// 		"request_id": "req_67890",
	// 		"user_id":    "admin",
	// 	},
	// )
}

// 示例6：记录验证错误
func ExampleValidationError() {
	// errors.LogValidationError(
	// 	"email",
	// 	"invalid-email-format",
	// 	map[string]interface{}{
	// 		"request_id": "req_11111",
	// 		"user_id":    "guest",
	// 	},
	// )
}

// 示例7：在控制器中使用
/*
func (c *UserController) Register(ctx *gin.Context) {
	// 获取请求上下文
	requestID := ctx.GetHeader("X-Request-ID")
	userIP := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	// 记录注册开始
	start := time.Now()

	// 业务逻辑...

	// 记录注册结果
	opLogger := audit.NewOperationLogger()
	opLogger.LogUserRegister(
		userID,
		username,
		userIP,
		audit.Success,
		map[string]interface{}{
			"duration_ms": time.Since(start).Milliseconds(),
			"user_agent":  userAgent,
		},
	)
}
*/

// 示例8：在错误处理中使用
/*
func (s *UserService) ConnectDatabase() error {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		errors.LogDatabaseConnectionError(
			config.Host,
			config.Port,
			config.Username,
			err,
			map[string]interface{}{
				"database": config.Database,
				"ssl_mode": config.SSLMode,
			},
		)
		return err
	}
	return nil
}
*/
