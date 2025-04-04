# handler

Contains code that are responsible for receiving requests from clients and sending responses back to clients. 
This layer is responsible for handling the transport layer protocol i.e. REST, GraphQL, gRPC.
This allows transport specific details to not leak into the business logic layer (controller) thus enabling a single 
business logic to be exposed via multiple protocols.
This layer also handles the authentication & authorization logic thus keeping the controller layer purely for business 
logic.
