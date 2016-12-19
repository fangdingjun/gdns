gdns
====

gdns is a dns proxy server


features
=======

- support forward the query by rule,
     different domains use different upstream server
- support ip black list
- support google https dns

usage
=====
    
    go get github.com/fangdingjun/gdns
    cp $GOPATH/src/github.com/fangdingjun/gdns/config_example.yaml config.yaml
    vim config.yaml
    $GOPATH/bin/gdns -c config.yaml
