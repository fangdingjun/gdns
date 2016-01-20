#!/usr/bin/env python
import math
import netaddr
import urllib2
url = "http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest"
data = urllib2.urlopen(url).read()

with open("delegated-apnic-lastest", "w") as fp:
    fp.write(data)

results = []
with open("delegated-apnic-latest") as fp:
    for line in fp:
        # format
        # apnic|CN|ipv4|1.0.0.0|256|2012232|allocated
        l = line.strip()
        if not l:
            continue
        if l[0] == '#':
            continue
        lns = l.split("|")
        if len(lns) != 7:
            continue
        if lns[2] != 'ipv4':
            continue
        if lns[1] != 'CN':
            continue
        ip = lns[3]
        mask = 32-int(math.log(int(lns[4]), 2))
        cidr = "%s/%d" % (ip, mask)
        results.append(cidr)
r = netaddr.cidr_merge(results)

with open("region_cn.txt", "w") as fp:
    for a in r:
        fp.write("%s,%s,%d\n" % (a.ip,a.netmask, a.prefixlen))
