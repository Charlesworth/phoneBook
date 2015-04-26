# phoneBook	<a href="https://travis-ci.org/Charlesworth/phoneBook"><img src="https://travis-ci.org/Charlesworth/phoneBook.svg?branch=master"></a>

A HTTP service implementing a phone/address book. Written in Go and using [BoltDB] as a storage engine and [httprouter] for route demux. Written by Charles Cochrane.

[BoltDB]: https://github.com/boltdb/bolt
[httprouter]: https://github.com/julienschmidt/httprouter

Content:
- [Usage](#usage)
- [API](#api)
- [JSON Format](#json-format)
- [How to Run](#how-to-run)
- [Possible Future Features](#possible-future-features)
- [Limitations](#limitations)
- [Testing](#testing)

## Usage

There are 4 main interactions that can be made to this application via the API:

- add/update an entry
- delete an entry
- list all entries
- search entries by surname

An entry in the phone book contains a first name, surname, telephone number and series of address fields. All entries are input and output in JSON format.

## API

### - Add/Update an entry:
##### `PUT` appUrl/
Adds or updates a single entry in the phoneBook using JSON in the request body. Please input JSON in the format described in the [JSON Format](#json-format) section of this README.

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

## How to Run

### Download binary executables

Binaries for Windows and Linux are available on the [Releases] section of this project. I cannot test Mac builds but Apple users will be able to follow the build instructions below instead.

After downlaoding the binaries, run the appropriate one for your system. When the phoneBook server starts it will be accessable on [http://localhost:2000/], ready and waiting for your HTTP API calls. Please note that phoneBook creates a phoneBook.db file in its current directory. If you want to reset the data, feel free to stop the server, delete the .db file and restart the server.

### Download source and building/running yourself

Building or Running a .go file will require you to have installed Go on your machine. Easy instructions for setting up Go can be found at [golang.org].

If you have the [Go] language installed on your machine, building the code yourself, using the "go install" command or simply running the code will also work. For example:

	$ go get github.com/Charlesworth/phoneBook
    	# cd into this phoneBook directory in your system
	$ go run phoneBook.go

As above, the server will start on [http://localhost:2000/] but you can easily change the port number in the source code.

[Releases]: https://github.com/Charlesworth/phoneBook/releases
[http://localhost:2000/]: http://localhost:2000/
[Go]: http://golang.org/
[golang.org]: http://golang.org/

## Possible Future Features

- User aliases, so multiple phone books can be stored
- A generic storage interface, so multiple storage engines could be plugged in
- Integrate [Twilio] to optionally text any new phone entry to confirm its the correct number
- Add other entry types: i.e. email address, FB profile, ect.

[Twilio]: https://www.twilio.com/

## Limitations

- Listing all entries output JSON is not sorted into any order
- Doesn't account for repeat entries, i.e. if you know two people called John Smith
- Can only input one entry at a time
- No security of any kind, anyone could grab your phonebook records

## Testing

Please see [Travis] for the most recent build test results. Test coverage currently at 84.8% Total.


[Travis]: https://travis-ci.org/Charlesworth/phoneBook
