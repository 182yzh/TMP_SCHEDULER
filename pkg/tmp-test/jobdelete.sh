s=$(kubectl get jobs | awk '{print $1}')
    echo $s
    for podname in $s
    do  
        kubectl delete job $podname
    done

