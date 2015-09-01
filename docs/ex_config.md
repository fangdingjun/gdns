Configure example
=============

Configure file
======
Use `gdns -h > config.ini` generate a example configure file

a configure file like this:

```conf
bind = :53  # the address bind to
blacklist =   # the blacklist file
configUpdateInterval = 0  # Update interval for re-reading config file set via -config flag. Zero disables config file re-reading.
debug = false  # output debug log, default false
logfile =  error.log # the logfile, default stdout
server = filter1.txt,udp:8.8.8.8:53  # special the filter and the upstream server to use when match
    #       format:
    #           FILTER_FILE_NAME,PROTOCOL:SERVER_NAME:PORT
    #       example:
    #           filter1.json,udp:8.8.8.8:53
    #               means the domains in the filter1.json will use the google dns  server by udp
    #       you can specail multiple filter and upstream server        
    #         
upstream = udp:114.114.114.114:53  # the default upstream server to use
```
comamnd
`gdns -config dns.ini`
use the dns.ini as a configure file

Filter file
===========
The filter file is a domains name list

command line
`--server domain1.json,udp:8.8.8.8:53`
means the domain name listed in domoin1.json will use 8.8.8.8 as the upstream server through udp

a filter file like this

```json
{
    "twitter.com":1,
    "facebook.com":1,
    "google.com":1
}
```

you can special multiple filter file and upstream dns server

Blacklist file
=============
The blacklist file contains the ip that the message will be dropped when the ip dispeared in the upstream server reply

the blacklist file like this

```json
{
    "113.123.21.43":1,
    "31.53.23.12":1
}
```
