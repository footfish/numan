#!  /bin/bash
#sample script for batching commands 
COMMAND='num add'
BATCHFILE="batchnumbers.txt"
ARG2="test.com"
ARG3="test carrier"

for LINE in $(cat $BATCHFILE) 
do
	echo -n "$LINE, " 
	$COMMAND "$LINE" "$ARG2" "$ARG3"  >/dev/null
	if [[ $? -eq 0 ]]
	then
		echo "TRUE"
	else
		echo "FALSE"
	fi
done
