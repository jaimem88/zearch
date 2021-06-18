# zearch
Zendesk search CLI (coding challenge)

## Overview

Using the provided data (tickets.json and users.json and organization.json) write a simple command line application to search the data and return the results in a human readable format.

## Trade-offs

- I chose Go because it's my strongest language. However, it's not the best tool for string processing and search.
Perhaps Python or Ruby would have made my life easier.
- The lack of generics in Go makes the code look very repetitive. Maybe I could have used higher order functions to try
to DRY the code somehow.
-

## Assumptions

Data:
- JSON files do not contain duplicate data

Relationships:
- An Organization has many users
- An Organization has many tickets
- A user belongs to one organization
- A ticket belongs to one organization

No relationships (for simplicity)
- A ticket does not have users (in a real life scenario that probably wouldn't be the case).
  I could assume that a ticket's `submitter_id` is the `user._id` but it would complicate search more.


- Exact match by string, including capitalization.
