It just works. If not, please create **an Issue**.

# Simple as hell!

Written in `net/http`, using `go-redis` and `pgx`

# What does it do?

Stores healthcheck results in **postgresql** and
has a webhook token feature (stores tokens in **Redis**)

Tokens are required for log receiving, but not for
anything else

# API

| Method | URL                 | Params                                            | Definition                                                         |
|--------|---------------------|---------------------------------------------------|--------------------------------------------------------------------|
| POST   | /                   | _Authorizaion_ (header) - **"Token <token_str>"** | receive logs in standard format (token required with prefix Token) |
| GET    | /logs               |                                                   | list logs                                                          |
| GET    | /logs?name_filter=  | _name_filter_ - **"my service"**                  | list logs that were sent with given name ("name" == name_filter)   |
| DELETE | /logs               |                                                   | delete all logs                                                    |
| DELETE | /logs?clear_before= | _clear_before_ - **2006-01-02T15:04:05Z07:00**    | delete logs that were sent before given datetime                   |
| POST   | /tokens             |                                                   | retrieve a new token                                               |
| DELETE | /tokens?token=      | _token_ - **<token_str>**                         | delete an existing token                                           |

# Launch 
```
docker compose up -d --build
```