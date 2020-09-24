
while [[ true ]]
do

    s=$(kubectl get pods | grep Completed | sed 's/-/ /g' | awk '{print $1}')
    for filename in  ${s[*]}
    do
#       	 echo "delete job " $filename
		 kubectl delete -f $filename.yaml
    done

#	dpods=$(kubectl get pods | grep Terminating | awk '{print $1}')
#    echo ${dpods[*]}
#	for podname in  ${dpods[*]}
 #   do 
#		echo "delete terminating pod" $podname
 #       kubectl delete pod $podname
#		jobname=$(echo $podname | sed 's/-/ /g' |awk '{print $1}')
#		kubectl delete job $jobname
#		echo "delete terminating job" $jobname
#		kubectl delete pod $podname	
  #  done



    ss=$(kubectl get pods | awk '{print $1}'|wc -l)
    if [[ ss -lt 1 ]] ; then
        break
    fi
done

rm app*.yaml
