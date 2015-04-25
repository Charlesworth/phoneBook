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

### Add/Update an entry:
##### `PUT` appUrl/
Adds or updates a single entry in the phoneBook using JSON in the request body. Please put the JSON in the following format for each entry:

{"Surname":"Smith","Entries":{[{"Firstname":"John","TelNo":"1234567890","Line1":"5a","Line2":"Fake St","TownCity":"Fakeville","CountyState":"Fakeshire","Country":"England","ZipPostal":"AA1 2BB"}]}

### Delete an entry: 
##### `DELETE` appUrl/[Surname]
Delete all entries with under [Surname]
##### `DELETE` appUrl/[Surname]/[Firstname]
Delete single entry with [Firstname]  [Surname]

### List all entries: 
##### `GET` appUrl/
List all of the entries present in the phoneBook

### Search entries: 
##### `GET` appUrl/[Surname]
List all entries with [Surname]
##### `GET` appUrl/[Surname]/[Firstname]
Get single entry with [Firstname]  [Surname]

## Limitations

### 
