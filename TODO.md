# Actions that I want to achieve

For learning purposes, This is my Todo list

!! 
Please ask chatGPT to generate a use case where you can practice 
- grRPC
- API proto,
- mircroservice,
- API gateway (authen with JWT, loadbalancer, keyclock ...)
- MongoDB or PostgreSQL (pgx, sqlc and migrate) with horizontal scalability (DB replication)
- caches with redis
- event-driven with NATS
- monitoring with victoria metrics + victoria Los
!!

## Todo

- [x] Handle errors
- [x] gracefull shutdown
- [x] PostgresSql
- [x] use go-migrate for DB migration
- [x] add mockery
- [ ] Add github actions (tests, release-please)
- [ ] Add an ORM with [sqlc](https://github.com/sqlc-dev/sqlc)
- [ ] Nginx or traefik as reverse proxy and load balancer
- [ ] Authen system (JWT or basic Authen or OAuth2)
- [ ] HTTPS
- [ ] manage users
- [ ] finish product
- [ ] complete database schema
- [ ] [cleanenv](https://pkg.go.dev/github.com/ilyakaznacheev/cleanenv), [viper](https://github.com/spf13/viper), [cobra](https://github.com/spf13/cobra),  

## Work In progress

- [x] Use mock in unit test
- [ ] correct FIXMEs
- [ ] Update doc

## Stack to build

- gRPC
- API protobuf
- microservices archi
- Traefik as load balancer (extremely easy, great for microservices)
- MongoDB for NoSQL or (PostgreSQL with `pgx` + `sqlc` + `migrate`)
- postman for API
- Redis for caching
- Victoria Metrics + Victoria Logs for monitoring
- Event-driven with Kafka (Best choice) or NATS
- websocket ??
- JWT bearer token for auth
- GKS (Google K8 Engine)
- keycloak for IdP/SSO
- GitHub Actions for CI/CD
- Terraform for IAS
- docker for container
- vault for secrets
