# gdns
a dns proxy server write by go

gdns much like dnsmasq or chinadns, but it can run on windows.

Features
========

support different domains use different upstream dns servers

support contact to the upstream dns server by tcp or udp

support blacklist list to block the fake ip

Install
=======

```bash
# get the depended library
go get github.com/miekg/dns
go get github.com/vharitonsky/iniflags
    
git clone https://github.com/fangdingjun/gdns
cd gdns
go build
    
# generate a sample config file
./gdns -dumpflags > dns.ini
    
# edit the dns.ini and run, need root privileges to bind on port 53
sudo ./gdns -config dns.ini
    
# test it
dig @localhost twitter.com
```

Arguments
===========

use `gdns -h` to show the command line arguments.

all arguments can specialed in config file or in command line.

there is a sample file in the `config_sample/` directory.

Third-part library
==================
use 
[dns](https://github.com/miekg/dns)
library to parse the dns message

use
[iniflags](https://github.com/vharitonsky/iniflags)
library to process the command line arguments and the config file
