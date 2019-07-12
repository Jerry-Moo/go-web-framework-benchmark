#!/bin/bash

for file in `ls *.csv`
do
    awk -f test.awk $file > t_$file
done
