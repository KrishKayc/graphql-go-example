# Go GraphQL Application Using Echo, graphql-go.
  Sample app which depicts a Cab Booking platform. The app supports graphQL queries to book a nearby cab and view booking history for a user.

## Setup sample Data
   ``` Run db/install.sql db/dump.sql ```

## Running the app
```
  go build  (once done executable will be generated in the same directory)
  ./bookCab
  
 ```
## Running Tests
```
    cd booker/tests
    go test -coverpkg=..
    
    cd graphql/fields/tests
    go test -coverpkg=..
```

## Explore graphQL endpoints
   Open **GraphiQL** App and use 'http://locahost:3000/graphql' to explore the types and work with queries and mutations


