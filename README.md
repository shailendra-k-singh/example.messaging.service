# example.messaging.service
A sample messaging service. User can create,read all created messages,or read or delete a specified message (using message id) using provided REST APIs. Additionally user can also perform a palindrome check on the message text while retrieving the message. 

## REST API documentation
Refer swagger.yaml in Project's root directory.

## Architecture
The application starts a HTTP server and listens on the specified port (default is 8090) for incoming requests and performs REST based CRUD operations. Additionally, it also performs request tracing for each request using Jaeger utilities.

An example flow for a REST POST call is shown below:
<br>
<br>
![Example Message create flow (POST)](https://github.com/shailendra-k-singh/example.messaging.service/blob/master/images/POST.png?raw=true)

## Deploy
To start the project (includes build): `docker-compose up`
<br>
To shutdown the project: `docker-compose down`

For CLI testing:
To build the project: `make build`
<br>
To run tests: `make test`

## Observability and metrics
To access Metrics: `http://localhost:14269/metrics`
<br>
To access Jaeger UI for request tracing: `http://localhost:16686/search`

Sample tracing snippet from Jaeger UI for a POST flow:
<br>
<br>
![Example tracing for POST call](https://github.com/shailendra-k-singh/example.messaging.service/blob/master/images/trace.png?raw=true)

## Usage

1. Bring up the project using command `docker-compose up`
2. Invoke the below endpoints for respective operations:
	- Create Message: POST `http://localhost:8090/v1/messages` ( with json body e.g. {"text": "sample"})
	- Retrieve all messages: GET `http://localhost:8090/v1/messages`
	- Retrieve a specific message: GET `http://localhost:8090/v1/messages/{id}` ( a valid positive integer id, e.g. `http://localhost:8090/v1/messages/1`)
	- Retrieve a specific message and check if the message text is palindrome: GET `http://localhost:8090/v1/messages/{id}?is-palindrome` ( a valid positive integer id, e.g. `http://localhost:8090/v1/messages/1?is-palindrome`)
	- Delete a specfic message: DELETE `http://localhost:8090/v1/messages/{id}` ( a valid positive integer id, e.g. `http://localhost:8090/v1/messages/1`)
