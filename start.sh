#!/bin/bash

# 启动 consul dev
# echo "启动 consul"
# consul agent -dev &
mkdir -p /home/hubing/mygo/src/consul_test/nsq1/data
mkdir -p /home/hubing/mygo/src/consul_test/nsq2/data

# 启动 nsq
echo "启动 nsqlookupd"
nsqlookupd --http-address 0.0.0.0:4161 --tcp-address 0.0.0.0:4160 &

echo "启动 nsqd1"
nsqd    --lookupd-tcp-address 0.0.0.0:4160 \
        --http-address 0.0.0.0:4151 \
        --tcp-address 0.0.0.0:4150 \
        -data-path /home/hubing/mygo/src/consul_test/nsq1/data &


echo "启动 nsqd2"
nsqd    --lookupd-tcp-address 0.0.0.0:4160 \
        --http-address 0.0.0.0:5151 \
        --tcp-address 0.0.0.0:5150 \
        -data-path /home/hubing/mygo/src/consul_test/nsq2/data &

echo "启动 nsqadmin"
nsqadmin --lookupd-http-address=127.0.0.1:4161 &