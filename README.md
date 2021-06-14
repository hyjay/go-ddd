# go-ddd
An example of [tactical domain-driven design](https://docs.microsoft.com/ko-kr/azure/architecture/microservices/model/tactical-ddd).

The project is about to demonstrate a few points:
- `internal/kit` provides a comprehensive and abstracted API for publishing and subscribing domain events via Google Pub/Sub.
- An example of a general user account service having some RESTful API endpoints - sign up user and get user.
- An example of publishing and handling a domain event, to implement a requirement `Send a welcome email to a user when the user signed up`.
- The code design follows the principles of [Clean Architecture(or Hexagonal architecture)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) focused in
separating each concern of pure domain/business logic, non-functional requirements such as monitoring/logging/etc, and implementations of 
database/RPC/etc at a low infrastructure level.
