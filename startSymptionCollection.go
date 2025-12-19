package main

import "fmt"



func main(){
// Process request
	/* Command Syntax :startProcessTraces.go -n <namsesapce> -c <"podnames">
		Example:
			startProcessTraces.go -n miniudm -c "uecm"
			startProcessTraces.go -n dracvnf -c "uecm nim"
	**** Concurrency validaiton check = 1  **** 
		1. Get the namespace from the arguement.
    	2. Get the pod from the arguement to capture trace. 
    	3. Check the namespace exists in the cluster.
    	4. Check the deployments of the pod exists in the namespace.
    	5. Check the readiness and liveness 
    	6. Calculate the number of Process inside the mcc platform for the required process.
    
    Record the start timestamp of execution also record the  currnet timestamp of the pod 
    Routine :: 1 
    	Parallely enable the trace of each process.
    Routine :: 2 
    	Parallely enable the pcap.
    Routine :: 3 
    	Exceute the pytbot command on testclient pod and wait the command's output status
	    Routine :: 4
	    	watch the /cmconfig.log --> Parallely, read the file for any run time error [mapping should be done using keywords]
	    							--> create a channel and inform to the user if any starts/error was logging.
	    Routine :: 5
	    	watch the /logstore/TspCore for every nth Sec --> Parallely,  check for any cores.
	    									              --> create a channel and inform to the user if any coring happens
	    Routine :: 6
	    	watch the /RTPTraceError --> Parallely, read the file for any run time error [mapping should be done using keywords]
	    							--> create a channel and inform to the user if any starts/error was logging.
	    Routine :: 7
	    	watch the /Envoy --> Parallely, read the file for any run time error [mapping should be done using keywords]
	    					--> create a channel and inform to the user if any starts/error was logging.
	    Routine :: 8
	        watch the /dumplog  --> Parallely, read the file for any run time error [mapping should be done using keywords]
	        					--> create a channel and inform to the user if any starts/error was logging.
    Routine ::  9
    	Once the pybot cmd status is set to complete
    		Routine :: 10
      			store the process traces into single file , can be used for analysis and also for symton collecitons
		    Routine ::  
		    	Parallely Disables the trace of each process.
		    Routine ::
		    	Parallely Disables the pcaps.
		    Routine :: 
		        Stops the process for to watch the /cmconfig.log
		        Store the traces for analysis or symptom collection
		    Routine :: 
		        Stops the process for to watch the /logstore/TspCore
		        Store the traces for analysis or symptom collection
		    Routine :: 
		        Stops the process for to watch the /RTPTraceError
		        Store the traces for analysis or symptom collection
		    Routine :: 
		        Stops the process for to watch the /Envoy
		        Store the traces for analysis or symptom collection
		    Routine :: 
		        Stops the process for to watch the /dumplog
		        Store the traces for analysis or symptom collection
	Copy all the informatio to local path (All traces, pcap, log.html)*/
}