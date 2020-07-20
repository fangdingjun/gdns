# gdns
a dns proxy server write by go


Features
========

support tcp, udp, tls, DoH(dns over https/http2)

Usage
=======

```bash
go get github.com/fangdingjun/gdns
cp $GOPATH/src/github.com/fangdingjun/gdns/config_example.yaml config.yaml
vim config.yaml
$GOPATH/bin/gdns -c config.yaml -log_level DEBUG
```

Third-part library
==================
use 
[dns](https://github.com/miekg/dns)
library to parse the dns message



### Support me ###

If you find my work useful, you consider support me for further develop.

Paypal:   [Support me on paypal.com](https://www.paypal.me/DingjunFang)

BTC: 3DSr6qNrrH4dQ3aTBhtkyYt94batcehmRf

