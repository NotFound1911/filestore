#!/bin/bash
# 检查service进程
check_process() {
    sleep 1
    res=`ps aux | grep -v grep | grep "cmd/bin" | grep $1`
    if [ -n "$res" ]; then
        echo -e "\033[32m 已启动 \033[0m" "$1"
        return 1
    else
        echo -e "\033[31m 启动失败 \033[0m" "$1"
        return 0
    fi
}

# 编译service可执行文件
build_service() {
    go build -o cmd/bin/$1 service/$1/main.go
    resbin=`ls cmd/bin/ | grep $1`
    echo -e "\033[32m 编译完成: \033[0m cmd/bin/$resbin"
}

# 启动service
run_service() {
    nohup ./cmd/bin/$1 >> $logpath/$1.log 2>&1 &
    sleep 1
    check_process $1
}
# 创建运行日志目录
logpath=./tmp/log
mkdir -p $logpath
services="
upload
file_manager
account
apigw
transfer
"
# 执行编译service
mkdir -p cmd/bin/ && rm -f service/cmd/*
for service in $services
do
    build_service $service
done

# 执行启动service
for service in $services
do
    run_service $service
done

echo '微服务启动完毕.'