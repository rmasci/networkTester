Network Tester was written to test IPSEC tunnels by downloading 100mb of data and logging the progress. It runs both as a client and
as a server depending on how it's started.  I've tried to expand this a bit.  The first thing I normally do is setup the server portion. 
For this you can run from the command line:
	
	networkTester -s
	
Optionally you can specify:
	
	-p <port number> (default 8888)
	-L <log File>
	-b <size in mb> 

The -b switch sets the packet size to send, by default it's 100mb so on slower networks you might want to scale that down a bit to 50mb or 10mb:
	
	networkTester -s -p 8080 -b 10 -L /tmp/networkTest.log
	
NetworkTest will not daemonize so if you want to fire and forget and leave it running:
	
	nohup networkTester -s -p 8080 -b 10 -L /tmp/networkTest.log &

At this point you could use a web server and hit the networkTest test page from the command line it would look like this:
	
	curl http://<networkTester IP>/
	HTTP 1.1/ 200 OK
	Date: 1/10/2015 09:36:51pm
	From: 172.16.1.45:57327
	To: docker02 - :8888

To start the client you can run networkTest as follows:
	
	networkTester -c -a <networkTester server IP>

That will try to connect to the server on port 8888, Optionally you can specify:
	
	-p <port number> Default is 8888
	-L <log file>
	-t <timeout in seconds>

So to connect to the server example above:
	
	networkTester -c -a 10.100.1.15 -p 8080 

Program will output the following:

	Client side:
		2015/01/10 16:37:36 Connected to 10.100.1.15 on port 8888
		2015/01/10 16:38:18 200 OK, downloaded 100.00m in 42.33s, 18.90Mbps
		2015/01/10 16:38:18 Connected to 10.100.1.15 on port 8888
		2015/01/10 16:40:00 200 OK, downloaded 100.00m in 42.23s, 18.90Mbps

	Server /tmp/networkTest.log:
		2015/01/10 21:37:20 Connect From: 172.16.1.45:57346, to: 10.100.1.15:8888
		2015/01/10 21:38:03 Sent: 100.00m From: 10.100.1.15:8888 to 172.16.1.45:57346 in 42.10s 19.00Mbps
		2015/01/10 21:38:03 Connect From: 172.16.1.45:57363, to: 10.100.1.15:8888
		2015/01/10 21:38:45 Sent: 100.00m From: 10.100.1.15:8888 to 172.16.1.45:57346 in 42.10s 19.00Mbps

Again the program doesn't daemonize so if you want to fire and forget to leave it running in the background:
	
	nohup networkTester -c -a 10.100.1.15 -p 8080 -L /tmp/networkTest.log &

For some odd reason I needed the ability to generate HTML data for a database functionality so networkTester also has the ability of spitting out 
random HTML by accessing the url: http://<networkTest IP>/dbGen

	Data Testing docker02
	
	1/10/2015 10:10:57pm, HTTP 1.1/ 200 OK
	First Name	Last Name	Sold Today
	RAYMOND	ADAMS	85.36
	FRANK	WRIGHT	30.32
	...
	HARRIS	14.46
	ANDREW	GREEN	84.88
	ERIC	HILL	82.77
	