# reverseproxy
Reverse proxy for redirecting requests based on the selected load balancing policy 

The project contains a reverse proxy server that is configurable and based on the load balacing strategy set in the configuration file redirects the request to an available host chosed by the selected algorithm.

## Configuration file
In the configuration file, under the `listen` key, you should set the ip address and port where the server listens to. The default value is `0.0.0.0` and `8002`.
Under the `services` key you have to specify an array of services with the following configuration
  - `name` name of the service
  - `domain` that is used to match the upstream service where the request will be redirected. 
  - `lbPolicy` the load balancing policy used for getting the host that will process the request
  - `timeout` timeout for processing a request
  - `retries` number of retries for getting a host
  - `healthCheckTimeout` interval at which the healt check to be performed
  - `hosts` list of { `address`(string), `port`(number) } for all hosts where it should check
  
 
##Install
Project was written in golang and used go version 1.19

Clone the repository and run the following commands
``` go install ```
and
``` go build ```

To run the application run the command
``` go run main.go ```
