# CodeVault Operation Manual
## STEPS
```
# 1. 找个临时目录
mkdir ~/test-project
cd ~/test-project

# 2. 初始化 Git 并写点东西
git init
echo "Hello DevNexus" > README.md
git add .
git commit -m "First commit"

# 3. 关键一步：添加你的本地服务器作为远程仓库
# 注意：这里的 'demo.git' 可以随便起名，我们的代码会自动创建它
git remote add origin http://localhost:8080/demo.git

# 4. 推送！
git push --set-upstream origin master
```
## RESULT
```
Enumerating objects: 3, done.
Counting objects: 100% (3/3), done.
Writing objects: 100% (3/3), 218 bytes | 218.00 KiB/s, done.
Total 3 (delta 0), reused 0 (delta 0), pack-reused 0
To http://localhost:8080/demo.git
 * [new branch]      master -> master
Branch 'master' set up to track remote branch 'master' from 'origin'.
```

