# polybed
Polymer, go and mongodb

## How build server

To build the polybed server, in polybed root directory, execute:

```
go build
```

This command generates the polybed executable file.

## MongoDB connection string settings

The connection string used to connect to the Mongo DB is modifiable changing the
environment variable `POLYBED_DB_URL`.
