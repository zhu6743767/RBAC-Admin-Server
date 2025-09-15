package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"rbac.admin/config"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	config   *config.Config
	filename string
	mu       sync.RWMutex
	watcher  *fsnotify.Watcher
}

// NewConfigManager 创建新的配置管理器
func NewConfigManager(filename string) (*ConfigManager, error) {
	cm := &ConfigManager{
		filename: filename,
	}

	// 初始加载配置
	if err := cm.load(); err != nil {
		return nil, fmt.Errorf("初始加载配置失败: %w", err)
	}

	// 设置文件监听
	if err := cm.setupWatcher(); err != nil {
		log.Printf("文件监听设置失败: %v", err)
		// 继续运行，但不支持热重载
	}

	return cm, nil
}

// GetConfig 获取当前配置（线程安全）
func (cm *ConfigManager) GetConfig() *config.Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// load 重新加载配置
func (cm *ConfigManager) load() error {
	data, err := os.ReadFile(cm.filename)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	var newConfig config.Config
	if err := yaml.Unmarshal(data, &newConfig); err != nil {
		return fmt.Errorf("解析YAML配置失败: %w", err)
	}

	cm.mu.Lock()
	cm.config = &newConfig
	cm.mu.Unlock()

	log.Printf("✅ 配置已重新加载")
	return nil
}

// setupWatcher 设置文件监听
func (cm *ConfigManager) setupWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	cm.watcher = watcher

	go func() {
		defer watcher.Close()
		
		// 添加文件监听
		err := watcher.Add(cm.filename)
		if err != nil {
			log.Printf("添加文件监听失败: %v", err)
			return
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				
				// 检测文件写入完成
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("📝 检测到配置文件变化，重新加载...")
					
					// 延迟加载，避免文件写入过程中的读取问题
					time.Sleep(100 * time.Millisecond)
					
					if err := cm.load(); err != nil {
						log.Printf("❌ 配置重载失败: %v", err)
					} else {
						log.Printf("🔄 配置热重载成功")
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("文件监听错误: %v", err)
			}
		}
	}()

	return nil
}

// SaveConfig 保存配置到文件（配置回写）
func (cm *ConfigManager) SaveConfig(newConfig *config.Config) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := yaml.Marshal(newConfig)
	if err != nil {
		return fmt.Errorf("配置序列化失败: %w", err)
	}

	// 创建临时文件
	tempFile := cm.filename + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("写入临时文件失败: %w", err)
	}

	// 原子性替换文件
	if err := os.Rename(tempFile, cm.filename); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("替换配置文件失败: %w", err)
	}

	cm.config = newConfig
	log.Printf("💾 配置已保存到文件")
	return nil
}

// UpdateConfig 更新特定配置项
func (cm *ConfigManager) UpdateConfig(updater func(*config.Config)) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 创建配置的副本
	newConfig := *cm.config
	updater(&newConfig)

	return cm.SaveConfig(&newConfig)
}

// Close 关闭配置管理器
func (cm *ConfigManager) Close() error {
	if cm.watcher != nil {
		return cm.watcher.Close()
	}
	return nil
}