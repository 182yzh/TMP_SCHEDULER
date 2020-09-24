#!/bin/bash
	order=$1
	echo "order =" $order
	line=`echo ${@:2}`
#	echo $line
	gpuinfo=$(echo $line | grep -Po  '"detail".*?\]\}\]\}' | head -1)
	#echo $gpuinfo
	gpus=$(echo $gpuinfo | grep -Po  '\[".*?"\]')
	#echo $gpus
	tasknum=$(echo $gpus|grep -o '\["'|wc -l)
	#echo $tasknum
	need=$(echo $gpus|grep -Po '\[".*?"\]' |head -1 |grep -o '"' |wc -l)
	need=`expr ${need} / 2`
	#echo $need
	name=$(echo $line |grep -Po '"jobid"\:".*?"')
	name=${name:9}
	name=${name%?}
	name=$(echo $name | sed 's/_//g')
	name="app"$order
	#echo $name
	echo $name $need $tasknum
	if [[ ${need} -gt 8 || `expr ${need} \* ${tasknum}` -gt 24 ]]
	then
    	need=2
		tasknum=8
	fi
	./create_file.sh $name $need $tasknum
	
