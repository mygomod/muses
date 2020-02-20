# 贡献代码指南

- Go >= 1.12
- MySQL >= 5.7

## Git 工作流

### fork 代码

1. 访问 https://github.com/i2eco/shop-api
2. 点击 "Fork" 按钮 (位于页面的右上方)

### Clone 代码

```bash
cd $GOPATH/src/github.com/i2eco/
git clone https://github.com/<your-github-account>/muses
cd caller
git remote add upstream 'https://github.com/i2eco/muses'
git config --global --add http.followRedirects 1
```

### 创建 feature 分支

```bash
git checkout -b feature/my-feature 
```

### 同步代码

```bash
git fetch upstream
git rebase upstream/master
```

### 提交 commit

```bash
git add .
git commit
git push origin my-feature
```
### 提交 PR

```bash
访问 https://github.com/i2eco/muses, 

点击 "Compare" 比较变更并点击 "Pull request" 提交 PR。
```
