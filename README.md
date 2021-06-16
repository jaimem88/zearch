# zearch
Zendesk search CLI (coding challenge)

## Overview

Using the provided data (tickets.json and users.json and organization.json) write a simple command line application to search the data and return the results in a human readable format.

## Assumptions

Data:
- JSON files do not contain duplicate data

Relationships:
- An Organization has many users
- An Organization has many tickets
- A user has one organization
- A ticket has one organization

No relationships (for simplicity)
- A ticket does not have users (in a real life scenario that probably wouldn't be the case).
  I could assume that a ticket's `submitter_id` is the `user._id` but it would complicate search more.
