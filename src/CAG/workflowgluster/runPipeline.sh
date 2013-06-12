for i in `cat sampleid.txt` ; do
	echo "./workflowgluster -config gluster.json -sn $i | tee $i.log" | qsub -cwd -N j$i -pe orte 4
done
