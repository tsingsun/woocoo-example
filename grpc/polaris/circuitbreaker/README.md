# 使用故障熔断

## 构建

```bash
cd ./provider
go build -o provider
cd ./consumer
go build -o consumer
```

## 进入控制台

预先通过北极星控制台创建对应的服务，如果是通过本地一键安装包的方式安装，直接在浏览器通过127.0.0.1:8080打开控制台

设置服务的熔断策略，错误率 1% ,让熔断容易发生.

## 启动服务

```bash
cd ./provide
./provide
# 熔断实例
./provide --cb=true
# 调用者
cd ./consumer
./consumer
```

## 验证

```bash
curl http://localhost:12000/echo
```

重复执行,可以看到在前面几次请求中,出现错误,后面的请求正常的越来越多.因为错误服务器熔断后,请求都转移到正常的服务器上了.