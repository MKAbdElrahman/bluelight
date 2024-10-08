# Bluelight 

## Contract Based Architecture


```mermaid
graph LR
 
    handler[HTTP handler] -- uses --> service[Service]

    service[Service] -- uses --> infrastructure[Other Infrastructure]
    service -- uses --> db[Repositories]
    db -- uses --> postgres[(PostgreQSL)]
    handler -- uses --> request_contracts[Request Contracts]

     handler -- uses --> response_contracts[Response Contracts]
```



- Contracts are explicit and put in their own  layer.


- The Handlers will use the contracts layer to read http requests and validate they are not malformed and also for sending responses.


- To prepare responses, handlers talk to the application services layer which should be representative of the usecases.


- The domain objects should have no clue that they are bing used in a web application or stored using a sql database. I'm strongly againest putting annotaions on core domain entities.




## Features

- Sending JSON Respnses and Parsing JSON Requests
- Runtime app configutation
- Database setup and confguration
- SQL migrations
- CRUD operations
- Optimistic Concurrency Control
- Filtering, Sorting, and Pagination
- Rate limiting
- Sending Emails and Calling external APIs
- Graceful shutodown
- Authentication And Authorization
- Background operations
- CORS
- Metrics and Monitoring



##  Design Guidelines 

### Context Usage
 **Context is for Cancellation:** Use `context.Context` only for cancellation and timeouts, not for passing dependencies.

 **Avoid `context.Value` for Dependencies:** Passing loggers or other data via `context.Value` hides dependencies and risks runtime errors due to lack of type safety.

 **Prefer Explicit Injection:** Always pass dependencies like loggers directly in constructors or function parameters for clarity and safety.

Read more at [Context is for cancelation
](https://dave.cheney.net/2017/01/26/context-is-for-cancelation)

## Acknowledgments 

- I learned alot from Alex Edwards books (I think its the best resource to learn Go).

- I also learned alot from Kent Beck books and his ideas about software design. I recommend alot reading his
  recent Book Tidy First. 

- I also can't recommend more Mark Richards and Neals Ford writings about software architecture.

Thank you All

Mohamed Kamal
