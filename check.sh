#!/bin/bash
attempts=0
ROUTER_IP=$1
if [ -z "$ROUTER_IP"]; then ROUTER_IP=192.168.1.1; fi
echo Router: $ROUTER_IP
log_result () {
    echo $(date +"%D %R ")$@ >> /var/log/check.log
}

while true; do
    sleep 90
    ping -c 1 $ROUTER_IP & wait $!
    if [ $? != 0 ]; then
        log_result  "Ping $ROUTER_IP failed. Strike $attempts"
        attempts=$((attempts+1))
        if [ $attempts -gt 10 ]; then
            log_result "Rebooting the system in 30 seconds"
            sleep 30
            reboot
        fi
    else
        attempts=0
    fi
done
