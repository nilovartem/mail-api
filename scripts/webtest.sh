#!/bin/sh
response=$(curl -X POST --user "beck.cierra@wb.ru:beck" localhost:8080/beck.cierra@wb.ru ... 2>/dev/null)
link=$(echo "$response")
if [[ ${link} == "" ]];then
    echo "test not ok, got empty link"
    exit
fi
response=$(curl --output - -X GET localhost:8080/get/$link)
if [[ ${response} != *"link expired"* ]];then
    echo "test ok"
else
    echo "test not ok"
fi