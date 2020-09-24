
while [[ true ]]
do

    s=$(kubectl get job | awk 'NR>1{print $1 " " $2}'| sed 's/\// /g' | awk '{if ($2==$3 ){print $1}}') #awk '{if ( $2 == $3 ) {print $1}}')
#	echo $s	
    for filename in  ${s[*]}
    do
      	 echo "delete job " $filename
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
