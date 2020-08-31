# genload

Simple command line tool to generate a load for a service.

This is a **WORK IN PROGRESS**. Use at your own risk.

## Usage

```$ genload <nr-of-total-calls> <url>```

```$ genload <nr-of-total-calls> <nr-of-calls-in-parallel> <url>```

`nr-of-total-calls` total number of calls to make to the service listening at `url`.
When the value is `0`, the application will keep calling the `url` until stopped by CTRL+C.


`nr-of-calls-in-parallel` number of calls made concurrently.
If not specified `1` is used as value, meaning no concurrency.

`url` the url to call

The output will show the call number, thread nr,
elapsed time in milliseconds
and either the returned HTTP status code or the returned error.

## Example

```$ genload 100 http://localhost:8080/api/users```

This command will call the service listening at url `http://localhost:8080/api/users`,
100 times, one call after the other, without any concurrency.

```$ genload 0 5 http://lcoalhost:8080/api/users```

This command will keep sending calls to the url `http://localhost:8080/api/users`,
until the application is stopped. It will use 5 concurrent calls.

```$ genload 10000 1000 http://localhost:8080/api/users```

This command will send a total of 10000 calls to the url
`http://localhost:8080/api/users`, using 1000 concurrent calls.
