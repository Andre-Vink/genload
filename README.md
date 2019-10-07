# genload

Simple command line tool to generate a load for a service.

This is a WORK IN PROGRESS. Use at your own risk.

## Usage

```$ genload <nr-of-calls> <url>```

`nr-of-call` total number of calls to make to the service listening at `<url>`
`url` the url to call

## Example

```$ genload 100 http://localhost:8080/api/users```

This command will call the service listening at url `http://localhost:8080/api/users`, 100 times. These calls will be done immediatly in parallel without waiting for each other. The output will show the call number, the returned HTTP status code or the returned error.

## Future

The tool needs to be extended with an option that indicates the number of parallel calls. For instance:

```genload <nr-of-tatal-calls> <nr-of-calls-in-parallel> <url>```

The tool shall call the `url` for a total of `<nr-of-total-calls>` but with a maximum of `<nr-of-calls-in-parallel>` calls simultaneous. Until the total number of calls is reached, a new call is spawned if the number of parallel calls is not reached.
