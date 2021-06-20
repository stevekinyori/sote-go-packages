trap "exit" INT TERM ERR
trap "kill 0" EXIT

export ORGID=10000
export USERID=soteuser
export DEVICE_ID=`date +%s`

echo $DEVICE_ID > coverage.out

nats sub $ORGID.$USERID &