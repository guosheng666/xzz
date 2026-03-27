# MySQL 数据同步工具规格说明

## Why
需要一款轻量级的MySQL数据同步工具，通过运行SQL文件实现数据同步，支持不同数据库之间的数据迁移。

## What Changes
- 新建一个跨架构的MySQL数据同步CLI工具
- 支持Windows、arm64、x86、龙芯(loongarch64)架构
- 通过终端参数输入SQL文件路径和目标数据库连接信息
- 支持执行包含DDL和DML语句的SQL文件

## Impact
- 新增能力：执行SQL文件进行数据库同步
- 涉及系统：命令行工具

## ADDED Requirements

### Requirement: 跨架构支持
工具必须能够在以下架构上编译和运行：
- Windows x86
- Windows x64
- Linux arm64
- Linux x86
- Linux loongarch64

#### Scenario: 在不同架构上运行
- **WHEN** 用户在对应架构的机器上执行工具
- **THEN** 工具能够正常启动并显示帮助信息

### Requirement: SQL文件执行
工具通过终端参数接收配置信息，包括：
- SQL文件路径（必需）
- 目标数据库连接信息（host, port, user, password, database）
- 字符集设置
- 执行超时时间

#### Scenario: 正常执行SQL文件
- **WHEN** 用户提供SQL文件路径和目标数据库连接信息
- **THEN** 工具成功连接数据库并执行SQL文件中的语句

#### Scenario: SQL文件不存在
- **WHEN** 用户提供的SQL文件路径不存在
- **THEN** 工具显示错误信息并退出

#### Scenario: 参数不足
- **WHEN** 用户未提供必需的参数
- **THEN** 工具显示错误信息和帮助信息

### Requirement: 事务支持
工具支持事务执行模式：
- 默认使用事务执行，出错时自动回滚
- 可选禁用事务（autocommit模式）

#### Scenario: 执行成功
- **WHEN** SQL文件中所有语句执行成功
- **THEN** 工具提交事务并显示成功信息

#### Scenario: 执行失败
- **WHEN** SQL文件中某条语句执行失败
- **THEN** 工具回滚事务并显示错误信息

### Requirement: 批处理执行
工具支持分批执行SQL文件中的语句，避免内存溢出

#### Scenario: 大文件执行
- **WHEN** SQL文件包含大量语句
- **THEN** 工具分批读取并执行，不占用过多内存

### Requirement: 执行进度显示
工具显示执行进度信息

#### Scenario: 显示进度
- **WHEN** 正在执行SQL文件
- **THEN** 显示当前执行的语句和进度

## MODIFIED Requirements
无

## REMOVED Requirements
无

## 技术实现
- 编程语言：Go (跨平台编译支持好)
- 数据库驱动：go-sql-driver/mysql
- 命令行解析：cobra/pflag
- 文件处理：逐行读取SQL文件
- 编译方式：Go cross-compilation

## 命令行参数说明

```
mysql-sync-tool.exe --sql-file <path> --host <host> --port <port> --user <user> --password <password> --database <database>

参数说明：
  --sql-file    SQL文件路径（必需）
  --host        数据库主机地址（默认：127.0.0.1）
  --port        数据库端口（默认：3306）
  --user        数据库用户名（默认：root）
  --password    数据库密码（默认空）
  --database    数据库名称（必需）
  --charset     字符集（默认：utf8mb4）
  --timeout     执行超时时间秒（默认：300）
  --no-transaction  禁用事务模式（默认：false）
```
