#!/bin/bash

echo "停止现有容器..."
sudo docker-compose down

echo "重新构建应用..."
sudo docker-compose build --no-cache aq3cms

echo "启动服务..."
sudo docker-compose up -d

echo "等待服务启动..."
sleep 10

echo "查看应用日志..."
sudo docker-compose logs aq3cms

echo "服务状态:"
sudo docker-compose ps
