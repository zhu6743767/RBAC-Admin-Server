package email

// emailStore 邮箱验证码存储结构
// 包含邮箱地址和对应的验证码

type emailStore struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// EmailStoreMap 全局邮箱验证码存储映射
// 使用emailID作为键，存储邮箱和验证码信息
var EmailStoreMap = map[string]emailStore{}

// Set 存储邮箱验证码信息
// emailID: 唯一标识
// email: 邮箱地址
// code: 验证码

func Set(emailID, email, code string) {
	EmailStoreMap[emailID] = emailStore{
		Email: email,
		Code:  code,
	}
}

// Verify 验证邮箱验证码
// emailID: 唯一标识
// email: 邮箱地址
// code: 用户输入的验证码
// 返回验证结果
func Verify(emailID, email, code string) bool {
	info, ok := EmailStoreMap[emailID]
	if !ok {
		return false
	}
	if info.Email != email {
		return false
	}
	if info.Code != code {
		return false
	}
	return true
}

// Remove 删除邮箱验证码记录
// 验证成功或失败后调用，释放资源
func Remove(emailID string) {
	delete(EmailStoreMap, emailID)
}