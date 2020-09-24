s=$(kubectl get pods | awk '{print $1}')
    echo $s
    for podname in $s
    do  
        kubectl delete pod $podname
    done

