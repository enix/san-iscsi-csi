#!/bin/bash

# Copyright (c) 2021 Enix, SAS
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
# or implied. See the License for the specific language governing
# permissions and limitations under the License.
#
# Authors:
# Paul Laffitte <paul.laffitte@enix.fr>
# Arthur Chaloin <arthur.chaloin@enix.fr>
# Alexandre Buisine <alexandre.buisine@enix.fr>

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
