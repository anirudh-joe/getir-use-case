# REST API Endpoints

- Fetches the data in the provided MongoDB collection and returns the results
in the requested format.
- create(POST) and fetch(GET) data from an in-memory database.

# Mongo DB endpoint 
## Request URI
- httpMethod = POST
- http://localhost:8080/mongo

## Request Payload
The request payload of the first endpoint will include a JSON with 4 fields.
- “startDate” and “endDate” fields should contain the date in the following format “YYYY-MM-DD”.

-  “minCount” and “maxCount” are of type float and used for filtering the data. The documents should be between “minCount” and “maxCount”.

## Response Payload
Response payload should have 3 main fields.
> “code” is for status of the request. 
> “msg” is for description of the code
 + 0 means success.
 + 404 means 
    + http Method not supported. (http Header InternalServerError)
    + body is nil. (http Header InternalServerError)
    + unable to read the body. (http Header InternalServerError)
    + unable to map body to request object defined above. (http Header InternalServerError)
 > “records” will include all the filtered items according to the request. Response object contains 
+ key
+ createdAt
+ totalCount

# In-Memory DB endpoint 

## GET
### Request URI
> http://localhost:8080/in-memory?key=active-tabs
- The request uri of GET endpoint will have 1 query parameter. That is “key”
param holds the key (any key in string type)

### Response Payload
> Response payload of GET endpoint will return a JSON with 2 fields or error.
- “key” fields holds the key
- “value” fields holds the value


## POST
### Request URI
> http://localhost:8080/in-memory
### Request Payload
> The request payload of POST endpoint will include a JSON with 2 fields

- “key” fields holds the key (any key in string type)
- “value” fields holds the value (any value in string type)

### Response Payload
> The Response payload of POST endpoint will be same as request
