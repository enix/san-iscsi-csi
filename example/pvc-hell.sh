#!/bin/bash

number_regex='^[0-9]+$'
if [[ ("$1" != "apply" && "$1" != "delete") || ! "$2" =~ $number_regex ]]; then
  echo "usage: $0 [apply|delete] quantity"
  exit 1
fi

ACTION=$1
QUANTITY=$2

echo ${ACTION} ${QUANTITY} pods and PVC

rm -f /tmp/pvc-hell.yaml
for ((i=1; i <= QUANTITY; i++)); do
  sed -e "s/{TEST_ID}/$i/" pvc-hell.yaml >> /tmp/pvc-hell.yaml
done

kubectl $ACTION -f /tmp/pvc-hell.yaml
rm /tmp/pvc-hell.yaml
