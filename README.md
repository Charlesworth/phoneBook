# phoneBook

A HTTP service implementing a phone/address book. Written in Go and using [BoltDB] as a storage engine and [httprouter] for route demux. Written by Charles Cochrane.

[BoltDB]: https://github.com/boltdb/bolt
[httprouter]: https://github.com/julienschmidt/httprouter

## Usage

There are 4 main things you can do with this application:

- add/update an entry
- delete an entry
- list all entries
- search entries by surname

An entry in the phone book contains a first name, surname, telephone number and series of address fields. All entries are input and output in JSON format.

## API

### - Add/Update an entry:
##### `PUT` appUrl/
Adds or updates a single entry in the phoneBook using JSON in the request body. Please input JSON in the format described in the [JSON Format](#JSON-Format) section of this README.

### - Delete an entry: 
##### `DELETE` appUrl/[Surname]
Delete all entries with under [Surname]
##### `DELETE` appUrl/[Surname]/[Firstname]
Delete single entry with [Firstname]  [Surname]

### - List all entries: 
##### `GET` appUrl/
List all of the entries present in the phoneBook

### - Search entries: 
##### `GET` appUrl/[Surname]
List all entries with [Surname]
##### `GET` appUrl/[Surname]/[Firstname]
Get single entry with [Firstname]  [Surname]

## JSON Format

Data is input into the phoneBook in the following format:

```json
{
  "Surname": "Smith",
  "Entries": [
    {
      "Firstname": "John",
      "TelNo": "1234567890",
      "Line1": "5a",
      "Line2": "Fake St",
      "TownCity": "Fakeville",
      "CountyState": "Fakeshire",
      "Country": "England",
      "ZipPostal": "AA1 2BB"
    }
  ]
}
```

Please have the JSON in this format in the HTTP request body when making a PUT request, not as part of a query string. 

The "Entries" element is an array, meaning that multiple elements can be stored, so there can be lots of different first names per surname. PUT only supports the input of one "Entries" element per single PUT request.

## Possible Future Features

- User aliases, so multiple phone books can be stored
- A generic storage interface, so multiple storage engines could be plugged in
- Integrate [Twilio] to optionally text any new phone entry to confirm its the correct number

[Twilio]: https://www.twilio.com/

## Limitations

- Listing all entries output JSON is not sorted into any order
- Doesn't account for repeat entries, i.e. if you know two people called John Smith
- Can only input one entry at a time
- No security of any kind, anyone could grab your phonebook records
