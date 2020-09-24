go build

rm Schedule.log
touch Schedule.log

./poseidon --log_dir "./log" --v=4  --kubeConfig=/home/work/.kube/config  --schedulerName=poseidongpu

trap ' pkill -f poseidon ' SIGINT

grep SCHEDULE ./log/* > Schedule.log

#cat Schedule.log

