#!/bin/bash
# name gpuneed tasknums
#printf "%s %d %d\n" $1 $2 $3  

cp /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/gpu_spin.yaml /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/$1.yaml
# Task Name
sed -i '4s/replace-value/'$1'/' /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/$1.yaml
sed -i '10s/replace-value/'$1'/' /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/$1.yaml
sed -i '16s/replace-value/'$1'/' /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/$1.yaml

# Task Resource Need
sed -i '23s/replace-value/'$2'/' /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/$1.yaml
sed -i '27s/replace-value/'$2'/' /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/$1.yaml

# Task Number
sed -i '6s/replace-value/'$3'/' /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/$1.yaml
sed -i '7s/replace-value/'$3'/' /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/$1.yaml

# Task Time
sed -i '28s/replace-value/'$4'/' /home/work/gowork/src/github.com/kubernetes-sigs/poseidon/pkg/tmp-test/$1.yaml


