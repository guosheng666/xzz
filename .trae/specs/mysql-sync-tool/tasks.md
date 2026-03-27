# Tasks

- [x] Task 1: 创建项目结构和基础配置
  - [x] SubTask 1.1: 初始化Go module
  - [x] SubTask 1.2: 创建项目目录结构（cmd, pkg）
  - [x] SubTask 1.3: 添加依赖包（cobra, mysql driver）

- [x] Task 2: 实现命令行参数解析
  - [x] SubTask 2.1: 定义SQL文件路径参数flag（必需）
  - [x] SubTask 2.2: 定义数据库连接参数flag（host, port, user, password, database）
  - [x] SubTask 2.3: 定义字符集和超时参数flag
  - [x] SubTask 2.4: 定义事务模式参数flag

- [x] Task 3: 实现数据库连接模块
  - [x] SubTask 3.1: 创建数据库连接函数
  - [x] SubTask 3.2: 实现连接池管理
  - [x] SubTask 3.3: 添加连接测试功能

- [x] Task 4: 实现SQL文件解析和执行
  - [x] SubTask 4.1: 实现SQL文件读取功能
  - [x] SubTask 4.2: 实现SQL语句分割逻辑（按分号分割）
  - [x] SubTask 4.3: 实现事务执行模式
  - [x] SubTask 4.4: 实现非事务执行模式
  - [x] SubTask 4.5: 实现进度显示功能

- [x] Task 5: 添加帮助信息和错误处理
  - [x] SubTask 5.1: 实现详细的帮助信息
  - [x] SubTask 5.2: 添加文件不存在错误处理
  - [x] SubTask 5.3: 添加数据库连接错误处理
  - [x] SubTask 5.4: 添加SQL执行错误处理和回滚

- [x] Task 6: 配置跨平台编译
  - [x] SubTask 6.1: 创建编译脚本（build.sh/build.bat）
  - [x] SubTask 6.2: 测试各架构编译（windows-x86, windows-x64, linux-arm64, linux-x86, linux-loongarch64）

# Task Dependencies
- Task 2 依赖 Task 1
- Task 3 依赖 Task 1
- Task 4 依赖 Task 2 和 Task 3
- Task 5 依赖 Task 4
- Task 6 可独立进行
