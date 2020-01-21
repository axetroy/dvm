[English](README.md) | 中文简体

[![Build Status](https://github.com/axetroy/dvm/workflows/test/badge.svg)](https://github.com/axetroy/dvm/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/axetroy/dvm)](https://goreportcard.com/report/github.com/axetroy/dvm)
![Latest Version](https://img.shields.io/github/v/release/axetroy/dvm.svg)
![License](https://img.shields.io/github/license/axetroy/dvm.svg)
![Repo Size](https://img.shields.io/github/repo-size/axetroy/dvm.svg)

## dvm

dvm是用于管理Deno版本的命令行工具。

专注于最简单的版本管理方法。

特性:

- [x] 跨平台支持
- [x] 简单易用
- [x] 没有运行时依赖
- [x] 零配置
- [x] 完全兼容已安装的 Deno

### 使用方法

无论你是已安装 Deno 否不影响 dvm 的使用

```bash
# 安装 Deno
$ dvm install v0.26.0
$ deno --version
deno v0.26.0

# 使用另一个版本的 Deno
$ dvm install v0.25.0
$ dvm use v0.25
$ deno --version
deno v0.25.0

# 卸载 Deno
$ dvm uninstall v0.25.0

# 帮助信息
$ dvm --help
```

### 安装

如果你使用的是 Linux/macOS 系统，你可以运行以下命令安装

```shell
# 安装最新版
curl -fsSL https://raw.githubusercontent.com/axetroy/dvm/master/install.sh | bash
# 安装指定版本
curl -fsSL https://raw.githubusercontent.com/axetroy/dvm/master/install.sh | bash -s v0.1.0
```

或者

在 [release page](https://github.com/axetroy/dvm/releases) 页面下载你平台相关的可执行文件

然后设置环境变量

例如, 可执行文件放在 `~/bin` 目录

```bash
# ~/.bash_profile
export PATH="$PATH:~/bin"
```

最后，试一下是否设置正确

```bash
dvm --help
```

### 升级

你可以重新下载可执行文件然后覆盖

或者输入以下命令进行升级到最新版

```bash
$ dvm upgrade # 升级到最新版
$ dvm upgrade v0.2.0 # 升级到指定版本
```

### 卸载

运行以下命令卸载或者手动移除 `dvm` 可执行文件和 `$HOME/.dvm` 目录

```shell
$ dvm destroy
```

### 从源码构建

```bash
> go get -v -u github.com/axetroy/dvm
> cd $GOPATH/src/github.com/axetroy/dvm
> make build
```

### 测试

```bash
$ make test
```

### 相关

[justjavac/dvm](https://github.com/justjavac/dvm) Node.js implement

[imbsky/dvm](https://github.com/imbsky/dvm) Reason implement

### 开源许可

The [MIT License](LICENSE)
