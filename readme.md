## GoConchRepublicBackend

This is a GO rewrite of a Java/Vertx backend application that
collects, processes and inserts data into an AWS DynamoDB database supporting
the Alexa Service called "The Conch Republic".

Whereas the Java/Vertx version uses a Vertx Message bus within one main Lambda application, 
This project is written to have separate isolated modules that make use of a Serverless Lambda, Event driven architecture (using AWS EventBridge)


Process wise, there are 3 main modules
- Initiate
- Fetch
- Database

The Initiate module is invoked by way of an EventBridge CRON setup and, when invoked, simply
inserts 12 events back to EventBridge, each event having a detail string
of "1" - "12" respectively.

The Fetch module is invoked by way of a custom EventBridge event, (created
by the Initiate module). Lambda will create as many instances of the Fetch module
needed so requests can run as parallel as needed. Its job is to perform an
HTTP GET of Event Data for a given month (1 is the current month, 12 is the 12th month from the current), 
extract event data from the HTML using "**github.com/PuerkitoBio/goquery**", and then 
create individual event details that can be provided as JSON to be inserted
back into EventBridge, so that it can be picked up by the Database module.

The Database module is invoked by way of a custom EventBridge event (created by the Fetch module) and
will handle the insertion of the Event Data into a DynamoDB database. 





