#!/bin/bash
# script to run the sampling of cliquetree structure and learning parameters with hedden variables

TIMESTAMP=$(date +"%F_%T")
LOG="${TIMESTAMP}.log"

echo "Starting test script..."

for FILE in *.csv
	do
	NVAR="$(head -n1 $FILE | tr -cd ',' | wc -c)"
	NVAR=$(( $NVAR + 1 ))
	STOP=$(( ($NVAR+1)/2 ))
	STEP=$(( ($NVAR+9)/10 ))
	echo "Processing: ${FILE}, ${NVAR} variables"

	# number of times to repeat the whole experiment
	for I in {1..5}
		do
		# variate tree-width
		for K in {3,5,7,11}
			do
			# check if this network supports this tree-width
			if [ $NVAR -ge $(( $K + 2 )) ]
				then
				# variate amount of hidden variables
				for (( H=0; H<$STOP; H+=$STEP ))
					do
					echo "k=${K}, h=${H}, i=${I}"
					learn -f $FILE -s $FILE.fg$I -e 1e-2 -iterem 10 -h $H -k $K -check=true >> $LOG
					# ./learn -f $FILE -s "FILE.fg$I" -e 1e-4 -iterem 10 -h $H -k $K >> $LOG
				done
			fi
		done
	done
done
