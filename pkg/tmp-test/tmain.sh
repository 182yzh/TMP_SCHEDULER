#while read line
#do
##
#	s=$(kubectl get pods | grep Completed | awk '{print $1}')
#	#echo $s
#	for podname in $s
#	do
#    	kubectl delete pod $podname
#	done

#	s=$(kubectl get pods | grep Completed | sed 's/-/ /g' | awk '{print $1}')
#	for filename in  ${s[*]}
#	do
#   	kubectl delete -f $filename
#	done



i=1
while read line
do	
#	echo $line
	
	s=$(kubectl get pods | grep Pending | awk '{print $1}'|wc -l)
    while [[ $s -gt 5 ]]
    do  
        sleep 2s
        s=$(kubectl get pods | grep Pending | awk '{print $1}'|wc -l)
    done


	if [[ $s -lt 6 ]] ;then
		./parse_create.sh $i $line
		# submit --
		name=$(echo $line |grep -Po '"jobid"\:".*?"')
	    name=${name:9}
    	name=${name%?}
		name=$(echo $name | sed 's/_//g')
		name=app$i
 		kubectl create -f $name.yaml
    fi


	#./parse_create.sh $line
	((i++))


done < JobLog100.log


:<<'COMMENT'
while [[ true ]]
do 

	s=$(kubectl get pods | grep Completed | sed 's/-/ /g' | awk '{print $1}')
    for filename in  ${s[*]}
    do
        kubectl delete -f $filename
    done


    ss=$(kubectl get pods | awk '{print $1}'|wc -l)
    if [[ ss -lt 1]]
	then 
		break
    fi
done
COMMENT

#rm app*.yaml
