## GoConchRepublicBackend

This is a GO rewrite of a Java/Vertx backend application that collects, processes and inserts data into an 
AWS DynamoDB database supporting the Alexa Service - "The Conch Republic".

Whereas the Java/Vertx version uses a Vertx Message bus within one main Lambda application, this project 
is written to have separate go modules that make use of a Serverless Lambda, Event driven architecture 
(using a combination of AWS EventBridge, SNS and SQS).

Process wise, there are 3 main modules:
- initiate
- fetch
- database

The **initiate** module is invoked by way of an EventBridge CRON setup and, when invoked, simply
inserts 12 events back to EventBridge, each event having a detail string of {"month" : "1" } - "12" respectively.

The **fetch** module is invoked by way of a custom EventBridge event, (created
by the initiate module). Lambda will create as many instances of the fetch module
needed so requests can run as parallel as needed. It's job is to perform an
HTTP GET of Event Data for a given month (1 is the current month, 12 is the 12th month from the current), 
extract event data from the HTML using "**github.com/PuerkitoBio/goquery**", and then 
create individual event details that can be provided as a JSON string message to be sent into an SQS queue. 

The **database** module is invoked by way of an EventBridge CRON setup (to run a few minutes after the other two modules). 
It handles the insertion of the Event Data into a DynamoDB database by polling an SQS queue for event items sent by the fetch module 
 looping until all items have been processed. The 'timeout' for this lambda function is set to 20 seconds so as to add 
 some SQS Long polling capability.
- Note: When Fetch provides EventData for the Database module to process, one of the columns passed is 
'EventExpiry', which contains a int64/epoch value calculated based on the Event's EndDate. DynamoDB is configured to use 
the EventExpiry column data as TTL - ultimately allowing for the records to auto-purge after the EndDate has just passed.






