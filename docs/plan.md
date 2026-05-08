[X] add sliding window algorithm
[X] handle persistence (redis)
[X] writing test cases

[X] update GetRateLimiter, etc to match new design
[X] update factory.go (*actually removed since no longer necessary*)
[X] update main.go
[X] tests for handlers
[X] config rate limiters via env variables e.g. redis addr, limits, window sizes, etc.

### API
[X] choose algorithm via api
   [x] - metadata (e.g. api keys) *IN PROGRESS*
   [ ] - make a /config endpoint handler that will:
    - read X-API-KEY to determine user
    - decode the json body
    - validate the algorithm
    - call the KeyConfig