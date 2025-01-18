# rsscan

Simple command line interface to maintain a list of podcasts and download the latest episode.

### Motivation
I'm sick of podcast applications in every single device and I prefer to have the latest episodes in a server available in my home network to any of my devices.

### Usage

*Building*
```
ninja bin/rsscan
```

*Usage*
```
# add podcast
./bin/rsscan -add https://seradio.libsyn.com/rss

# list saved podcasts
./bin/rsscan -list
Software Engineering Radio - the podcast for professional software developers
  Episode: SE Radio 651: Paul Frazee on Bluesky and the AT Protocol
  Published:Fri, 17 Jan 2025 02:13:00 +0000

# remove podcast
./bin/rsscan -remove https://seradio.libsyn.com/rss

# update saved podcasts metadata and episodes (if new episode exists)
./bin/rsscan -update

# custom data directory (where metadata and audio is saved)
export RSSCAN_DATA_DIR=/var/lib/rsscan # should exists and have proper permissions
./bin/rsscan -add https://seradio.libsyn.com/rss

```
