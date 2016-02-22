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
```
Setting | Description |
--------|-------------|
`port`  | Ther port number on which to start the Phileas server | 
`mapquest_key` | Used for reverse geo coding. Can be gotten from [developer.mapquest.com](https://developer.mapquest.com) | 
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
### Google
Phileas needs to connec to google APIs for two things—

1. To be able to display maps in a browser (browser API key is required)
2. To query [Google Places API](https://developers.google.com/places/web-service/) for places information (server API key is required)

Both these keys can be gotten from [developers.google.com](https://developers.google.com/maps/signup?hl=en).
```
[google]
browser_key = xxxxxxxxxxxxxxxxxx
server_key = xxxxxxxxxxxxxxxxxx
```

# URLs
URL | Description | Expected HTTP code | Expected response |
----|-------------|--------------------|-------------------|
`/ping` | Healthcheck URL to see if the server is up and running | `200` | `pong` |
`/top` | Shows a google map with the locations overlayed | `200` | Google map markup |
`/countries.json` | Information about all the bookmarks in Phileas's database grouped by country; this information is encoded in the GeoJSON format | `200` | GeoJSON-encoded locations|
`/top.json` | Locations from Phileas's database encoded, again, in the GeoJSON format | `200` | GeoJSON-encoded locations|
`/loc/[location-id]` | Get's all images and their captions from Phileas's database | `200` | `HTML` with thumbnail images and captions |
