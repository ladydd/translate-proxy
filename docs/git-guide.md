# GitHub 新项目发布指南

## 前提条件（只需配置一次）

你已经完成了以下配置，不需要再做：

- `git config --global user.name "ladydd"` — 设置你的 Git 用户名
- `git config --global user.email "ladydd@163.com"` — 设置你的 Git 邮箱
- `ssh-keygen -t ed25519` — 生成 SSH 密钥对
- 公钥已添加到 GitHub（https://github.com/settings/keys）

---

## 方式一：网页 + 命令行

### 1. GitHub 上创建仓库

打开 https://github.com/new ，填仓库名，选 Public 或 Private，其他都不勾，点 Create repository。

### 2. 本地推送

在你的项目目录里打开终端，依次执行：

```powershell
git init
```
初始化 Git 仓库。在当前目录创建一个 `.git` 隐藏文件夹，让 Git 开始跟踪这个目录。

```powershell
git add .
```
把当前目录下所有文件添加到暂存区。`.` 表示所有文件。Git 不会自动跟踪文件，你得告诉它要管哪些。

```powershell
git commit -m "init"
```
提交暂存区的文件，生成一个版本快照。`-m "init"` 是提交说明，描述这次改了啥。

```powershell
git branch -M main
```
把当前分支重命名为 `main`。Git 默认分支名可能是 `master`，GitHub 现在用 `main`，统一一下。

```powershell
git remote add origin git@github.com:ladydd/仓库名.git
```
关联远程仓库。`origin` 是远程仓库的别名（约定俗成叫 origin），后面的地址就是你 GitHub 仓库的 SSH 地址。

```powershell
git push -u origin main
```
把本地代码推送到 GitHub。`-u` 表示设置上游关联，以后直接 `git push` 就行，不用再写 `origin main`。

---

## 方式二：纯命令行（GitHub CLI）

### 安装（只需一次）

```powershell
winget install GitHub.cli
```
安装 GitHub 命令行工具。

```powershell
gh auth login
```
登录 GitHub 账号，选 SSH 方式，跟着提示走。登录一次后永久有效。

### 新项目发布

在项目目录里执行：

```powershell
git init
```
初始化 Git 仓库。

```powershell
git add .
```
添加所有文件到暂存区。

```powershell
git commit -m "init"
```
提交第一个版本。

```powershell
gh repo create ladydd/仓库名 --private --source=. --push
```
一条命令完成三件事：在 GitHub 上创建仓库 + 关联本地目录 + 推送代码。
- `--private` 私有仓库（换成 `--public` 就是公开的）
- `--source=.` 指定当前目录为源码目录
- `--push` 创建完直接推送

---

## 后续更新代码

不管用哪种方式创建的，以后改了代码都是这三步：

```powershell
git add .
```
把所有改动添加到暂存区。

```powershell
git commit -m "改了啥写啥"
```
提交改动，写清楚这次改了什么。

```powershell
git push
```
推送到 GitHub。

---

## 常用场景速查

| 场景 | 命令 |
|------|------|
| 查看当前状态 | `git status` |
| 查看改了哪些内容 | `git diff` |
| 查看提交历史 | `git log --oneline` |
| 撤销还没 add 的修改 | `git checkout -- 文件名` |
| 撤销已经 add 但没 commit 的 | `git reset HEAD 文件名` |
