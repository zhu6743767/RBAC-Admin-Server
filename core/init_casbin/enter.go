package init_casbin

import (
	"fmt"
	"rbac_admin_server/global"

	casbin "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/sirupsen/logrus"
)

// InitCasbin 初始化Casbin权限管理器
func InitCasbin() (*casbin.CachedEnforcer, error) {
	// 创建GORM适配器
	adapter, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		return nil, fmt.Errorf("Casbin适配器初始化失败: %w", err)
	}

	// 加载模型配置文件
	modelPath := "config/casbin/model.conf"

	// 创建Casbin执行器
	enforcer, err := casbin.NewCachedEnforcer(modelPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("Casbin执行器初始化失败: %w", err)
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("Casbin策略加载失败: %w", err)
	}

	// 设置缓存过期时间（秒）
	enforcer.SetExpireTime(60 * 60) // 1小时

	logrus.Info("✅ Casbin权限管理器初始化成功")
	return enforcer, nil
}

// UpdatePolicy 更新Casbin策略
func UpdatePolicy(enforcer *casbin.CachedEnforcer, sub, obj, act string, newSub, newObj, newAct string) error {
	// 删除旧策略
	if _, err := enforcer.RemovePolicy(sub, obj, act); err != nil {
		return fmt.Errorf("删除旧策略失败: %w", err)
	}

	// 添加新策略
	if _, err := enforcer.AddPolicy(newSub, newObj, newAct); err != nil {
		return fmt.Errorf("添加新策略失败: %w", err)
	}

	// 保存策略
	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %w", err)
	}

	return nil
}

// AddRolePolicy 为角色添加策略
func AddRolePolicy(enforcer *casbin.CachedEnforcer, role, obj, act string) error {
	if _, err := enforcer.AddPolicy(role, obj, act); err != nil {
		return fmt.Errorf("添加角色策略失败: %w", err)
	}

	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %w", err)
	}

	return nil
}

// RemoveRolePolicy 移除角色策略
func RemoveRolePolicy(enforcer *casbin.CachedEnforcer, role, obj, act string) error {
	if _, err := enforcer.RemovePolicy(role, obj, act); err != nil {
		return fmt.Errorf("移除角色策略失败: %w", err)
	}

	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %w", err)
	}

	return nil
}