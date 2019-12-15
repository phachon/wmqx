#!/bin/bash

APP_NAME=wmqx

ROOT_DIR=$(cd "$(dirname "$0")";pwd)
APP_FILE=./bin/${APP_NAME}

PID_FILE=${ROOT_DIR}/logs/${APP_NAME}.pid
LOG_FILE=${ROOT_DIR}/logs/${APP_NAME}.log

function check_pid() {
    if [[ -f ${PID_FILE} ]];then
        pid=`cat ${PID_FILE}`
        if [[ -n ${pid} ]]; then
            res=`ps -p ${pid}|grep -v "PID TTY" |wc -l`
            return `echo ${res}`
        fi
    fi
    return 0
}

function start() {
    check_pid
    run_res=$?
    if [[ ${run_res} -gt 0 ]];then
        echo -n "${APP_NAME} is running already, pid="
        cat ${PID_FILE}
        return 1
    fi
	chmod +x ${APP_FILE}
    nohup ${APP_FILE} &> ${LOG_FILE} &
    echo $! > ${PID_FILE}
    echo "${APP_NAME} start running, pid=$!"
}

function stop() {
    pid=`cat ${PID_FILE}`
    kill ${pid}
    echo "${APP_NAME} stop."
}

function restart() {
    pid=`cat ${PID_FILE}`
    stop
    start
}

function status() {
    check_pid
    run_res=$?
    if [[ ${run_res} -gt 0 ]];then
        echo "status: start"
    else
        echo "status: stop"
    fi
}

function help() {
    echo "$0 start|stop|restart|status|pid"
}

function pid() {
    cat ${PID_FILE}
}

if [[ "$1" == "" ]]; then
    help
elif [[ "$1" == "stop" ]];then
    stop
elif [[ "$1" == "start" ]];then
    start
elif [[ "$1" == "restart" ]];then
    restart
elif [[ "$1" == "status" ]];then
    status
elif [[ "$1" == "pid" ]];then
    pid
else
    help
fi

