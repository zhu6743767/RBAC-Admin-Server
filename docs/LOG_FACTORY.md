# 日志工厂模式文档

## 概述

本项目实现了完整的日志工厂模式，支持多种日志实现（logrus、zap、zerolog）的动态切换和配置。通过工厂模式，系统可以在运行时根据配置选择不同的日志实现，无需修改代码。

## 架构设计

### 核心组件

1. **Logger接口** - 统一的日志操作接口
2. **LoggerFactory接口** - 日志工厂接口，用于创建日志实例
3. **FactoryRegistry** - 工厂注册中心，管理所有日志工厂
4. **LoggerManager** - 日志管理器，负责日志实例的生命周期管理

### 文件结构

```
core/logger/
├── logger_factory.go      # 工厂模式核心接口和注册中心
├── logrus_factory.go      # logrus日志工厂实现
├── file_output.go         # 文件输出支持
└── logger.go             # 原始日志实现（向后兼容）
```

## 使用方法

### 1. 基本使用

```go
// 使用工厂模式创建日志实例
factory, err := logger.GetFactory("logrus")
if err != nil {
    // 处理错误
}

logInstance, err := factory.Create(logConfig)
if err != nil {
    // 处理错误
}

logInstance.Info("日志消息")
logInstance.Close()
```

### 2. 配置文件

在配置文件中指定日志类型：

```yaml
log:
  type: "logrus"  # 支持的类型: logrus, zap, zerolog
  level: "info"
  format: "json"
  output: "both"
  log_dir: "logs"
  max_size: 100
  max_age: 7
  max_backups: 10
```

### 3. 支持的日志类型

| 类型   | 描述           | 状态   |
|--------|----------------|--------|
| logrus | 当前默认实现   | ✅ 已支持 |
| zap    | 高性能日志库   | 🔲 待实现 |
| zerolog| 零分配日志库   | 🔲 待实现 |

### 4. 扩展新的日志实现

实现新的日志工厂：

```go
type ZapFactory struct{}

func (f *ZapFactory) Type() string {
    return "zap"
}

func (f *ZapFactory) Create(config config.LogConfig) (Logger, error) {
    // 实现zap日志的创建逻辑
    return &ZapAdapter{...}, nil
}

// 注册工厂
func init() {
    logger.RegisterFactory("zap", &ZapFactory{})
}
```

## 配置选项

### 日志级别
- `debug`: 调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

### 日志格式
- `json`: 结构化JSON格式
- `text`: 可读性文本格式

### 输出位置
- `stdout`: 控制台输出
- `file`: 文件输出
- `both`: 控制台和文件同时输出

### 文件轮转
- `max_size`: 单个文件最大大小(MB)
- `max_age`: 文件最大保留天数
- `max_backups`: 最大备份文件数
- `compress`: 是否压缩旧日志

## 向后兼容性

原有的日志系统仍然可用，但为了保持向后兼容性，建议使用工厂模式：

```go
// 旧方式（仍支持，但不推荐）
logger := logger.New(config)

// 新方式（推荐）
factory, _ := logger.GetFactory("logrus")
logger, _ := factory.Create(config)
```

## 性能考虑

- 工厂模式在创建日志实例时有一次性开销
- 日志实例创建后，性能与原有实现相同
- 支持日志实例的复用和缓存

## 错误处理

```go
// 获取不支持的日志类型
factory, err := logger.GetFactory("unsupported")
if err != nil {
    // 错误处理：不支持的日志类型
}

// 创建日志实例失败
logger, err := factory.Create(invalidConfig)
if err != nil {
    // 错误处理：配置错误或创建失败
}
```

## 测试

运行使用示例：

```bash
go run examples/factory_usage_example.go
```

## 迁移指南

从原有日志系统迁移到工厂模式：

1. **配置文件更新**: 添加 `type: "logrus"` 配置项
2. **代码更新**: 使用工厂模式创建日志实例
3. **测试验证**: 确保所有日志功能正常工作
4. **逐步迁移**: 可以分模块逐步迁移，保持向后兼容