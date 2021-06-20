# zearch

Zendesk search CLI (coding challenge)

## Overview

Using the provided data (tickets.json and users.json and organization.json)
write a simple command line application to search the data and returnt
he results in a human readable format.

## Running the app

The application was written in Go 1.16.4.

A [Makefile](./Makefile) is included to help you get started quickly. You can run `make help`
for a list of available targets. The following list requires Go installed locally.

To build and run the application:

  ```shell
  make run
  # or build and run manually an
  go build -o ./out/bin/zearch ./cmd/zearch/
  ./out/bin/zearch
  ```

To run unit tests:

  ```shell
  make test
  # or run manually (requires Go 1.16+)
  go test ./...
  ```

To run in a Docker container:

  ```shell
  make docker
  # or build image and run manually
  docker build -t my-zearch-image .
  docker run -ti --rm my-zearch-image -users data/users.json
  ```

To use your own set of data, you can specify the data filenames using flags, for example:

  ```shell
  ./out/bin/zearch -users my_users.json -organizations my_organizations.json -tickets my_tickets.json
  ```

## App design

This is a simple CLI app that uses [github.com/manifoldco/promptui](https://github.com/manifoldco/promptui)
to handle user input. On startup, the app will read the files `organizations.json`, `users.json` and `tickets.json`
and load them into memory using the [`internal/store/`](./internal/store) package.

### Store

The idea of the `store` package is to hold the data in different maps using the ID as the key.
It also holds a simple relationship between entities, simulating a database relationship. For example,
there is a list of user belonging to an organization, so that data aggregation can be done easily.

The data is read using Go's `encoding/json` package, which is converted into a slice of `map[string]interface{}`.
I decided to do this because it is simpler to use the `term` as a string and access the value of that `term`
in constant time from the maps.

I had previously considered using strong types per entity (see 
[types.go in the `prompts-and-storage` branch](https://github.com/jaimem88/zearch/blob/prompts-and-storage/internal/model/types.go)).
However, this would have led me to use Go's reflection to obtain the struct field to read te `value` and search for them.

### Search

Searching by ID of an entity is done in constant time thanks to the use of maps.
Searching by other terms requires iterating over all elements of that entity, for example all organizations,
and trying to find exact value matches per term.

### Trade-offs

- I chose Go because it's my strongest language. However, it's not the best tool for string processing and search.
Perhaps Python or Ruby would have made my life easier.
- The lack of generics in Go makes the code look very repetitive. Maybe I could have used higher order functions to try
to DRY the code somehow.
- Search by field other than `_id` is done in O(n). In future improvements, I could probably sort the values per field
and perform a more performant search.

### Assumptions

Data:
- JSON files do not contain duplicate IDs for an entity, if they do,
  only the latest one read will be available to be searched.

Relationships:
- An Organization has many users
- An Organization has many tickets
- A user belongs to one organization
- A ticket belongs to one organization

No relationships (for simplicity)
- A ticket does not have users (in a real life scenario that probably wouldn't be the case).
  I could assume that a ticket's `submitter_id` is the `user_id` but I chose to keep it simple.
  
Search:
- Exact match by string, including capitalization.


## Demo

![demo](data/demo.mov)
