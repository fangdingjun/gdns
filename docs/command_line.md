command line options
=====================

###bind###

   special the addres the dns server listen to
   
example: 
        `-bind 0.0.0.0:53`,
        `-bind :53`,
        `-bind 127.0.0.1:53`
        
###server###

special a filter file and the upstream dns server to use        
format

   **file_name**,**proto**:**addr**:**port**
   
   **file_name** is the file name contains the domain list
   
   **proto** is the upstream dns server protocol, `tcp` or `udp`
   
   **addr** is the ip address of upstream dns server
   
   **port** is the upstream dns server port
   
 this options can special multipe times
 
example:

`-server domain1.json:udp:8.8.8.8:53`,
    
`-server domain1.json:tcp:4.2.2.2:53`,
    
`-server domin2.json:udp:49.32.34.44:5353`
    
 see [example filter file](ex_config.md#filter-file)   
    
###upstream###

special the default upstream dns server   

format

**proto**:**addr**:**port**

   **proto** is the upstream dns server protocol, `tcp` or `udp`
   
   **addr** is the ip address of upstream dns server
   
   **port** is the upstream dns server port
   
example:

`-upstream udp:10.10.1.1:53`

###logfile###

special the file name the log save to

example:

`-logfile /var/log/gdns.log`

###debug###

output the debug log or not, default false

this options is only used for debugging

###blacklist###

special the blacklist file

if the reply dns message contains the ip in the blacklist, the message will be dropped

example:

`-blacklist fakeip.json`

see [example black list](ex_config.md#blacklist-file)