#!/bin/bash
pw=$(htpasswd -nbBC 10 kumquat "$1")
IFS=":"
read -r -a strarr <<< "$pw"
echo "${strarr[1]}"