# gdns
a dns proxy server write by go

gdns much like dnsmasq or chinadns, but it can run on windows.

Features
========

support different domains use different upstream servers

support contact to the upstream server by tcp or udp

support blacklist list to block the fake ip

Install
=======

    # get the depended library
    go get github.com/miekg/dns
    go get github.com/vharitonsky/iniflags
    
    git clone https://github.com/fangdingjun/gdns
    cd gdns
    go build
    
    # generate a sample config file
    ./gdns -dumpflags > dns.ini
    
    # edit the dns.ini
    sudo ./gdns -config dns.ini
    
    # test it
    dig @localhost twitter.com
    
Arguments
===========

use `gdns -h` to show the command line arguments.

all arguments can specialed in config file or in command line.

this is a sample file in the config_sample directory.
