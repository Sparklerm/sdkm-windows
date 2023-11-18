# sdkm-windows
一个适用于Windows的JDK管理工具

## 1. 介绍
sdkm-windows是一个适用于Windows的JDK管理工具，可以方便的管理多个JDK版本，支持JDK的安装、卸载、切换、查看、配置等功能。

目前支持的JDK类型有： `Oracle JDK`,`Zulu JDK`,`GraalVM JDK`

## 2. 指令
```pwsh
sdkm ls # 列出已安装的JDK版本
sdkm use <version> # 切换JDK版本
sdkm install <version> # 安装指定版本的JDK
sdkm remove <version> # 卸载指定版本的JDK
sdkm available <jdk type> # 列出可用的JDK版本 
```

## 3. 配置
`conf/config.json`为sdkm-windows的配置文件
- `JDK_DIR`：JDK安装目录
- `JDK_ENV_NAME`：JDK环境变量名

**` 避免使用Admin权限的情况下，采用用户环境变量，如已配置系统级环境变量会产生冲突 `**

`conf/jdk_version.json`为JDK版本配置文件，可自行拓展