# Phileas ![logo](http://i.imgur.com/7eY6CUs.png) [![Build Status](https://img.shields.io/travis/avadhutp/phileas/master.svg?style=flat)](https://travis-ci.org/avadhutp/phileas) [![CodeCov](https://img.shields.io/codecov/c/github/avadhutp/phileas.svg?style=flat)](https://codecov.io/github/avadhutp/phileas) [![GoDoc](https://godoc.org/github.com/avadhutp/phileas?status.png)](https://godoc.org/github.com/avadhutp/phileas)
 
> “I see that it is by no means useless to travel, if a man wants to see something new”

Phileas can tell you where you should travel next, based on your instagram likes. Phileas ingests your likes and drops them as markers on a map. It also uses Yelp's API and enriches your liked media with some more information—like places of interest, eateries, hotels, etc.

# Installation
1. Create an `ini` file as shown in the [settings](#settings) section of this README
2. Put it in `/etc/phileas.ini`; optionally, you can pass the location of the ini file at run time
3. Download Phileas: `go get github.com/avadhutp/phileas`
4. First, run setup as `phileas setup`; this will set up the required table structure 
5. Next, initiate the backfill as `phileas backfill`; this will populate the database with all the historic data
6. Finally, run the actual Phileas server by issuing `phileas start`; you can issue this command in parallel with the `backfill` command

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
client_id = xxxxxxxxxxxxxx
secret = xxxxxxxxxxxxxxxx
access_token = xxxxxxxx
```
### Yelp
Philea's use Yelp's API to enrich instagram likes with additional categorization information. Sign up for a Yelp API (v2) account and generate a new auth token at [www.yelp.co.uk/developers/manage_api_keys](https://www.yelp.co.uk/developers/manage_api_keys). Next, copy the required configuration into the `yelp` section of phileas's config file:
```
[yelp]
consumer_key = xxxxxxxxxxxxxxxx
consumer_secret = xxxxxxxxxxxxxxxx
access_token = xxxxxxxxxxxxxxxx
access_token_secret = xxxxxxxxxxxxxxxx
```

# URLs
URL | Description | Expected HTTP code | Expected response |
----|-------------|--------------------|-------------------|
`/ping` | Healthcheck URL to see if the server is up and running | `200` | `pong` |
`/top` | Shows a google map with the locations overlayed | `200` | Google map markup |
`/top.json` | Locations from Phileas's databse encoded in the GeoJSON format | `200` | GeoJSON-encoded locations|