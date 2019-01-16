# Nerf-Controller
This code interfaces and run the web app GUI for an automated nerf gun that makes live stream donations into shots at the streamer. Designed to run on a Raspberry Pi Zero W connected to a nerf gun or other method of firing projectiles with a user interface via a web app, the program interfaces with StreamLabs to record donations that will add rounds to the hopper to be fired. The only projectiles to be used are safe to fire at humans such as the Nerf Rival Rounds or similar. Future iterations may have a 3D printable mechanism that can be used to fire the rounds, see the deployment section for more information.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

[Golang v1.11.x](https://golang.org/doc/install "Golang")
[StreamLabs account](https://streamlabs.com/ "StreamLabs")
[StreamLabs API App](https://streamlabs.com/dashboard#/apisettings "StreamLabs API App")
make

### Installing

A step by step series of examples that explains how to get a development environment running. To build the system that will fire rounds see the Deployment section below.

1.  Install Golang
2. Create a SteamLabs account
3. Create a StreamLabs Application
4. Create a StreamLabsAPI.json file in the project home directory to have the correct information for the application you created.
	1. Example StreamLabsAPI.json:
```json
{
  "ClientID": "7FDyzIU5NPbDLJ0kvB5C5CYSay6VYxNoNmza0RW1",
  "ClientSecret": "wleRBri2UFhUFYBCdnDhOASgBm2uQ7H60vkC34hB",
  "RedirectURI": "http://localhost:8080/live"
}
```
5. (Optional) If you would like to change what port is used you can set the PORT environment veriable. by default the port used is 8080.
	1. Example of setting the port to 3000 instead of 8080 using a bash comand
	```
	PORT="3000"
	```
6. The setup should now be complete the program can be tested and run with the following make command from the project home directory.
```
make all run
```

If correctly setup you should see the following in the console.
```
go fmt
go test -v ./...
=== RUN   TestGetPort
--- PASS: TestGetPort (0.00s)
=== RUN   TestHomeHandler
--- PASS: TestHomeHandler (0.00s)
=== RUN   TestFireHandler
{Fri, 21 Dec 2018 22:54:35 EST} FIRE!!!
--- PASS: TestFireHandler (0.00s)
=== RUN   TestTokenHandler
{Fri, 21 Dec 2018 22:54:35 EST} https://streamlabs.com/api/v1.0/authorize?client_id=7FDyzIU5NPbDLJ0kvB5C5CYSay6VYxNoNmza0RW1&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Flive&response_type=code&scope=donations.read
{Fri, 21 Dec 2018 22:54:35 EST} https://streamlabs.com/api/v1.0/authorize?client_id=7FDyzIU5NPbDLJ0kvB5C5CYSay6VYxNoNmza0RW1&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Flive&response_type=code&scope=donations.read
--- PASS: TestTokenHandler (0.00s)
=== RUN   TestRandomPageHandler
Sorry but it seems this page does not exist...
Sorry but it seems this page does not exist...
Sorry but it seems this page does not exist...
--- PASS: TestRandomPageHandler (0.00s)
=== RUN   TestFire
{Fri, 21 Dec 2018 22:54:35 EST} FIRE!!!
--- PASS: TestFire (0.00s)
=== RUN   TestRandomValue
6
14
8
7
9
12
11
5
15
10
13
--- PASS: TestRandomValue (0.00s)
PASS
ok  	_/Users/rdufrene/work/nerf-contorller	(cached)
go build -o nerf-controller -v 
_/Users/rdufrene/work/nerf-contorller
go build -o nerf-controller -v ./...
./nerf-controller
Now listening to port :8080
```
Now visiting http://localhost:8080/ should display the user interface.

## Running the tests

The tests are run using the make command `make test` in the project home directory.

## Deployment

Setup raspberry pi. TODO: Still in progress.

## Built With

* [Golang](https://golang.org/) - Code Backend and Framework
* [StreamLabs](https://streamlabs.com/) - Donations Management
* Make - Code and Compilation Rules Management

## Contributing

Please read [CONTRIBUTING.md](https://github.com/crabbymonkey/nerf-contorller/blob/master/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/crabbymonkey/nerf-contorller/tags). 

## Authors

* **Ryan Dufrene** - *Initial work* - [crabbymonkey](https://github.com/crabbymonkey)

See also the list of [contributors](https://github.com/crabbymonkey/nerf-contorller/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* *Billie Thompson* ([PurpleBooth](https://github.com/PurpleBooth)) - README.md and CONTRIBUTING.md templates used
