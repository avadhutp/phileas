# Phileas ![logo](http://i.imgur.com/7eY6CUs.png) [![Build Status](https://img.shields.io/travis/avadhutp/phileas/master.svg?style=flat)](https://travis-ci.org/avadhutp/phileas) [![CodeCov](https://img.shields.io/codecov/c/github/avadhutp/phileas.svg?style=flat)](https://codecov.io/github/avadhutp/phileas) [![GoDoc](https://godoc.org/github.com/avadhutp/phileas?status.png)](https://godoc.org/github.com/avadhutp/phileas)
 
> “I see that it is by no means useless to travel, if a man wants to see something new”

Phileas can tell you where you should travel next, based on your instagram likes. Phileas ingests your likes and drops them as markers on a map. It also uses Yelp's API and enriches your liked media with some more information—like places of interest, eateries, hotels, etc.

# Settings
Phileas works off of an `ini` file. This supports the following configs:
### Common
Example: 
```
[common]
port = 8081
mapquest_key = XXXXX
google_maps_key = XXXXX
```
Setting | Description |
--------|-------------|
`port`  | Ther port number on which to start the Phileas server | 
`mapquest_key` | Used for reverse geo coding. Can be gotten from [developer.mapquest.com](https://developer.mapquest.com) | 
`google_maps_key` | Google Maps API key. Can be gotten from [developers.google.com](https://developers.google.com/maps/signup?hl=en) |
### Mysql
Example:
```
[mysql]
host = localhost
port = 3306
username = user
password = password
schema = phileas
```
### Instagram
Phileas needs to connect to the Instagram API to fetch your likes. For this, it needs to connect to your Intagram API Account, the details for which can be gotten from [www.instagram.com/developer/clients/manage](https://www.instagram.com/developer/clients/manage/). You need to 

1. Create a new client in the _Manage Clients_ section
2. And then copy the required config into the `instagram` section of phileas's config file as show below:
```
[instagram]
client_id = xxxxxx
secret = xxxxxxxx
access_token = xxxxxxxx
```
