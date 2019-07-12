#!/bin/bash

server_bin_name="gowebbenchmark"

. ./libs.sh

lenth=${#web_frameworks[@]}

test_result=()

# mac
cpu_cores=`sysctl -a | grep machdep.cpu.core_count | wc -w`
# linux
# cpu_cores = `cat /proc/cpuinfo | grep processor | wc -l`
if [ $cpu_cores -eq 0 ]
then 
    cpu_cores=1
fi

test_web_framework()
{
    echo "testing web framework: $2"
    # run server
    ./$server_bin_name '-wf' $2 '-s' $3 &
    # time sleep 2
    sleep 2

    # 获得并发数
    throughput=`wrk -t$cpu_cores -c$4 -d30s http://127.0.0.1:8080/hello -s pipeline.lua --latency -- / 16 | grep Requests/sec | awk '{ print $2 }'`
    echo "throughput: $throughput requests/second"
    test_result[$1]=$throughput

    pkill -9 $server_bin_name
    sleep 2
    echo "finished testing $2"
    echo
}

test_all()
{
    echo "###################################"
    echo "                                   "
    echo "      ProcessingTime  $1ms         "
    echo "      Concurrency     $2           "
    echo "                                   "
    echo "###################################"
    for ((i=0; i<$lenth; i++))
    do
        test_web_framework $i ${web_frameworks[$i]} $1 $2
    done
}

pkill -9 $server_bin_name

# time sleep测试模拟 io读写频繁 非cpu 耗时操作
# 并发相同情况下 time sleep 增加
echo ","$(IFS=$','; echo "${web_frameworks[*]}" ) > processtime-pipeline.csv
test_all 0 100
# test_all 0 5000
echo "0 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime-pipeline.csv
test_all 10 100
# test_all 10 5000
echo "10 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime-pipeline.csv
test_all 100 100
# test_all 100 5000
echo "100 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime-pipeline.csv
test_all 500 100
# test_all 500 5000
echo "500 ms,"$(IFS=$','; echo "${test_result[*]}" ) >> processtime-pipeline.csv

# time sleep测试模拟 io读写频繁 非cpu 耗时操作
# time sleep 相同的情况下 并发增加
echo ","$(IFS=$','; echo "${web_frameworks[*]}" ) > concurrency-pipeline.csv
test_all 30 100
echo "100,"$(IFS=$','; echo "${test_result[*]}" ) >> concurrency-pipeline.csv
test_all 30 1000
echo "1000,"$(IFS=$','; echo "${test_result[*]}" ) >> concurrency-pipeline.csv
test_all 30 5000
echo "5000,"$(IFS=$','; echo "${test_result[*]}" ) >> concurrency-pipeline.csv

mv -f processtime-pipeline.csv ./testresults
mv -f concurrency-pipeline.csv ./testresults

