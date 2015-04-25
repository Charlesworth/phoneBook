# phoneBook

A HTTP service implementing a phone/address book. Writen in Go and using [BoltDB] as a storage engine and [httprouter] for route demux.

[BoltDB]: https://github.com/boltdb/bolt
[httprouter]: https://github.com/julienschmidt/httprouter

## Usage

There are 4 main things you can do with this application:

- add/update an entry
- delete an entry
- list all entries
- search entries by surname

An entry in the phone book contains a fist name, surname, telephone number and series of address fields. All entries are input and output in JSON format.

## JSON Format


## API
All parameters must be sent in the request body as JSON, not as part of a query string. 

### Add/Update an entry: `PUT` [appUrl]/ 

### Delete an entry: `DELETE` [appUrl]/

### List all entries: `GET` [appUrl]/list

### Search entries by surname: `GET` [appUrl]/[Surname]

## Limitations