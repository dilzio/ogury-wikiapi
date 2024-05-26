# PelotechFun 

## Description

This project implements a wrapper API that aggregates data from the
Wikipedia [Pageviews API](https://wikitech.wikimedia.org/wiki/Analytics/AQS/Pageviews#Quick_start).
There are 3 endpoints:

1. **mostviewed**: given start and end dates, will return an aggregate ranking of top viewed articles
2. **viewcount**: given a start date, end date, and article name, will return the total views for that article in the
   date range
3. **mostviewedday**: given a 4-digit year, 2-digit month, and article name, will return the day the article had the
   highest number of views in that month

## Install and Run

A Dockerfile is provided for building, running tests, and running the app and is the suggested approach. The docker
commandline
can be installed from www.docker.com. To install and build the app:

1. Clone the repository to your local machine: `git clone https://github.com/dilzio/pelotechfun.git`
2. Cd to the top-level directory (where this README is located)
3. Build the docker image. This will also build the application: `docker build -t mtc-api .`

To run unit tests:
`docker run mtc-api go test ./storage ./indexer`

To run E2E integration test against live Wikipedia API:
`docker run mtc-api go test ./main`

To run the API (not necessary for tests):
`docker run -p 8080:8080 -it --rm --name mtc-api mtc-api`
## API Usage
The API is configured to run on localhost:8080. All calls are GET calls in keeping with REST norms and as such they can
be
called from a browser. Some example calls are below:

Find the day in July 2015 where the article "Albert_Einstein" had the most views:
`http://localhost:8080/mostviewedday/Albert_Einstein/2015/07`

reply:
```
{
 "startdate":"2015-07-01T00:00:00Z",
 "enddate":"2015-08-01T00:00:00Z",
 "articles":[
    {"name":"Albert_Einstein","views":17269,"time":"2015-07-23T00:00:00Z"}
  ]
 }
```

Find the set of most viewed articles from Jan 1-April 1 2021 (inclusive) in descending order
`http://localhost:8080/mostviewed/20210101/20210401`

reply:
```
{
 "startdate":"2021-01-01T00:00:00Z",
 "enddate":"2021-04-01T00:00:00Z",
 "articles":[
    {"name":"Main_Page","views":576633240,"time":"0001-01-01T00:00:00Z"},
    {"name":"Special:Search","views":121438577,"time":"0001-01-01T00:00:00Z"},
    {"name":"WandaVision","views":13407085,"time":"0001-01-01T00:00:00Z"},
    {"name":"Bible","views":12368446,"time":"0001-01-01T00:00:00Z"},
    {"name":"Deaths_in_2021","views":11305690,"time":"0001-01-01T00:00:00Z"},
    {"name":"Donald_Trump","views":10138579,"time":"0001-01-01T00:00:00Z"},
    {...},
}     
```

Find the total views for the article "Dua_Lipa" from Aug 15-Oct 31 2022 (inclusive)
`http://localhost:8080/viewcount/Dua_Lipa/20220815/20221031`

reply:
```
{
 "startdate":"2022-08-15T00:00:00Z",
 "enddate":"2022-10-31T00:00:00Z",
 "articles":[
    {"name":"Dua_Lipa","views":1378770,"time":"0001-01-01T00:00:00Z"}
    ]
}
```

## Notes:

- There is 100-day limit on the span between start and end dates for all api calls. This is essentially to guard
  against potential Wikipedia rate-limiting
- The API implements a basic local cache designed for demo and testing that stores the results of the API calls but
  never evicts
  and as such will eventually run out of memory if enough data is stored there.
- All results are aggregated in real time on every invocation from either cached values or values fetched from
  Wikipedia. If
  data from a particular date cannot be retrieved from one of these sources, the entire API invocation will fail to
  avoid
  returning incorrect values