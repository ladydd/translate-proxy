# Git + GitHub 首次配置指南（Windows）

新电脑或重装系统后，按这个文档配置一次就行。

---

## 1. 安装 Git

打开 https://git-scm.com/download/win ，下载安装，一路默认下一步就行。

装完后打开 PowerShell，验证一下：

```powershell
git --version
```

能看到版本号就说明装好了。

---

## 2. 配置身份信息

告诉 Git 你是谁，每次提交代码会带上这个信息：

```powershell
git config --global user.name "你的GitHub用户名"
```
设置用户名，会显示在每次 commit 记录里。

```powershell
git config --global user.email "你的GitHub邮箱"
```
设置邮箱，GitHub 通过这个邮箱关联你的账号。

验证配置：

```powershell
git config --global --list
```

---

## 3. 生成 SSH 密钥对

SSH 密钥用来让你的电脑和 GitHub 之间安全通信，不用每次输密码。

```powershell
ssh-keygen -t ed25519 -C "你的GitHub邮箱"
```

- `-t ed25519` 指定加密算法（目前最推荐的）
- `-C "邮箱"` 是备注，方便你知道这个 key 是谁的

执行后会问三个问题，全部直接回车就行：
- 保存位置 → 默认 `C:\Users\你的用户名\.ssh\id_ed25519`
- 密码 → 留空（不设密码，方便）
- 确认密码 → 留空

完成后会生成两个文件：
- `id_ed25519` — 私钥（绝对不能给别人，相当于你的身份证）
- `id_ed25519.pub` — 公钥（要添加到 GitHub 上）

---

## 4. 查看公钥

```powershell
cat ~/.ssh/id_ed25519.pub
```

输出类似这样一行：

```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI...... 你的邮箱
```

复制这一整行。

---

## 5. 把公钥添加到 GitHub

1. 打开 https://github.com/settings/keys
2. 点 `New SSH key`
3. Title：随便填，比如"我的电脑"、"公司电脑"（方便你区分是哪台机器）
4. Key type：保持默认 `Authentication Key`
5. Key：粘贴刚才复制的公钥
6. 点 `Add SSH key`

---

## 6. 验证连接

```powershell
ssh -T git@github.com
```

第一次会问你是否信任，输入 `yes`。

看到这个就说明成功了：

```
Hi ladydd! You've successfully authenticated, but GitHub does not provide shell access.
```

---

## 常见问题

### 换了电脑怎么办？

重新走一遍这个文档就行。每台电脑生成自己的密钥对，公钥添加到 GitHub。GitHub 上可以添加多个 SSH Key，互不影响。

### 私钥丢了/泄露了怎么办？

1. 去 https://github.com/settings/keys 删掉对应的公钥
2. 本地重新生成一对：`ssh-keygen -t ed25519 -C "邮箱"`
3. 把新公钥添加到 GitHub

### Permission denied (publickey) 怎么办？

说明 SSH Key 没配对。检查：
- `cat ~/.ssh/id_ed25519.pub` 看看有没有公钥
- GitHub 上有没有添加这个公钥
- 没有的话重新走第 3-5 步
