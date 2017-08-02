# wikinow

![alt text](https://raw.githubusercontent.com/kartiksura/wikinow/master/wikinow.png)


#Aim:
1. Horizontally scalable (add nodes to scale-up)
2. Fault tolerant (killing nodes should not affect other tasks/data)
3. Performance
4. Simplicity
5. REST APIs

##Design philosophy:
For simplicity and performance, the central idea is: Avoid mapping Wikipedia, but store all the previously traversed data for future use.
We store the results and also the sub results which can be used for future use.


Whenever a new request is initiated by the user, a Job is created. 
A job has the following:
1. Source wiki title
2. Destination wiki title
3. Timeout in seconds

The user gets a RequestID which can be used to get the status of the job. 

##API Spec
###The API for creating a job:
http://localhost:8080/request?src=Mike%20Tyson&dst=Heavyweight&timeOut=50

The output is a  json object containing the requestID:
{"ID":"15069b34-23c7-42a0-ad32-22a36fccdece"}


###The API for getting the job status:
http://localhost:8080/solution?id=15069b34-23c7-42a0-ad32-22a36fccdece

The output is a json object containing the status of the JOB:
{"ID":"15069b34-23c7-42a0-ad32-22a36fccdece","Status":"SOLVED_FROM_CACHE","Path":["Mike Tyson","South by Southwest","The Daily Mirror","Heavyweight"],"ProcessingTime":"5.0327ms"}

The Path field contains the path from Source to Destination


##TODO:
Unit test cases.
Stats monitoring.
Dockerisation.
Storing advanced sub-results
