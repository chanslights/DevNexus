# Step 1
    目标：实现CodeVault，适配Git Smarter HTTP

# Step 2
    目标：实现WebHook触发机制，让CodeVault主动通知OpsEngine。

    1. 在cmd/opsengine/main.go下实现handleWebhook函数
    2. 在internal/codevault/git/handler.go中支持获取最新的Commit ID，并发送请求

# Step 3
    目标：让OpsEngine在收到通知后，主动把代码clone下来，找到根目录下的.devnuxus.yaml并执行里面的指令

    1. ineteral/opsengine/pipeline/config.go下构造.devnexus.yaml结构体
    2. ineteral/opsengine/pipeline/parse.go实现FetchAndParse函数,拉取代码并解析
    3. cmd/opsengin/main.go修改handleWebHook函数，让它收到Webhook后调用FetchAndParse

# Step 4
    目标：容器化构建引擎,OpsEngine真正拉起一个Docker

    1. 安装Docker SDK
    2. internal/opsengine/docker/executor.go中初始化一个对象NewExecutor
    3. internal/opsengine/docker/executor.go封装一个RunStep，一个简易的Docker执行器
    4. 完善cmd/opsengine/main.go的handleWebhook函数

    ## 坑点
    1. 使用的Docker SDK版本有些不兼容。使用了replace强制替换依赖
    2. Cursor的Agent模式下会自动给我改一些奇怪的错误，还是CRAFT模式比较好

# Step 5
    目标：云原生部署

    1. 准备K8s，docker Desktop下开启K8s
    2. internal/opsengine/k8s/deployer.go下初始化部署器对象NewDeployer
    3. internal/opsengine/k8s/deployer.go增加UpdateImage 更新指定 Deployment 的镜像
    4. cmd/opsengine/main.go中的handleWebhook函数集成K8s能力

# Step 6
    目标：构建AI助手

    1. internal/ai/agent.go下定义一个标准的OpenAI格式客户端
    2. DeepSeek开放平台获取API key
    3. internal/opsengine/docker/executor.go把Dokcer Executor中的日志存到一个Buffer里面
    4. cmd/opsengine/main.go中添加逻辑：当捕获到错误是，呼叫AI

# Step 7
    缺点
    1. 数据目前都存在内存中，后续要存在DB中
    2. 并发控制还不完善
    3. 没有Web UI

    后续步骤
    1. 引入GORM和MySQL，对所有信息进行持久化
    2. 收到WebHook跑成PC模型，完成高并发调度
    3. 一个简单的Web界面
    4. 分布式与中间件，引入Redis做消息队列，引入MinIO存日志