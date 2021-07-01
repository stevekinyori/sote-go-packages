trap "exit" INT TERM ERR
trap "kill 0" EXIT

export ORGID=10000
export USERID=soteuser
export COGNITOID=5d5147e2-57fc-48a6-b493-1783931ae9c0
export DEVICE_ID=`date +%s`

echo $DEVICE_ID > .git/device.info

nats sub $ORGID.$USERID &