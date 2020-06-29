## GoConchRepublicBackend

This is a GO rewrite of a Java/Vertx backend application that
collects, processes and inserts data into an AWS DynamoDB database supporting
the Alexa Service - "The Conch Republic".

Whereas the Java/Vertx version uses a Vertx Message bus within one main Lambda application, this project 
is written to have separate modules that make use of a Serverless Lambda, Event driven architecture 
(using AWS EventBridge).

Process wise, there are 4 main modules:
- initiate
- fetch
- database
- cleanup

The **initiate** module is invoked by way of an EventBridge CRON setup and, when invoked, simply
inserts 12 events back to EventBridge, each event having a detail string
of {"month" : "1" } - "12" respectively.

The **fetch** module is invoked by way of a custom EventBridge event, (created
by the initiate module). Lambda will create as many instances of the fetch module
needed so requests can run as parallel as needed. Its job is to perform an
HTTP GET of Event Data for a given month (1 is the current month, 12 is the 12th month from the current), 
extract event data from the HTML using "**github.com/PuerkitoBio/goquery**", and then 
create individual event details that can be provided as JSON to be inserted
back into EventBridge, so that it can be picked up by the database module.

The **database** module is invoked by way of a custom EventBridge event (created by the fetch module) and
will handle the insertion of the Event Data into a DynamoDB database. 

Lastly, a separate **cleanup** module, invoked by way of an EventBridge CRON setup, occurs a few minutes after the 
above processes complete, and handles removing obsolete EventData from the DynamoDB table.





