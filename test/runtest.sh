#!/bin/sh

export MIHARI_DEBUG=YES

(cd ../; make)
if [ $? != 0 ];then exit 1; fi

../bin/mihari sh -c "seq 10 20 > generated_1.txt"
../bin/mihari sh -c "seq 20 30 > generated_2.txt"
../bin/mihari sh -c "paste generated_1.txt generated_2.txt > generated_3.txt"
../bin/mihari sh -c "paste generated_3.txt input_1.txt > generated_4.txt"
../bin/mihari sh -c "sort -n generated_4.txt > generated_5.txt"
../bin/mihari sh -c "cat -n generated_4.txt > generated_6.txt"

../bin/mihari --makefile
