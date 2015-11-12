# bigzelda
__BigZelda__ is a shortlink creator (get it? short Link? gimme the rupees)

# Features
A shortlink is the association of a base URL (origin) and a more simple key (token) randomly generated or given by the user.
__BigZelda__ offers via a REST API a shortlink service to __register__ shortlinks, __use__ shortlinks and __inspect__ shortlinks (to get the number of time it has been used).
These shortlinks are stored in a Redis. 

# Conf
There is a configuration file in the conf directory.
With this configuration file, you can change
- the port on wich the __BigZelda__ application API is exposed
- the token max size 
- the shortlink Time To Live (Duration in seconds before redis removes an unused shortlink)
- the redis connections details and credentials (you only need to set them if you are running the application as a go app and not packaged in docker)
- the logfile name 	(final logfile name are : filenameYYYYMMDD.log)
- the log level 	(DEBUG/INFO/WARN/ERROR)

# Run

## with Docker
Go in the __BigZelda__ directory and use
./install.sh
and
./start.sh

This should run in your docker 2 containers
- redis 		: 	a simple redis-server
- bigzelda		: 	the bigzelda go application


# Api
__BigZelda__ exposes a simple API on the port 80. To change this port, modify the command line in start.sh (... --publish __80__:8000 ... )

This API is accessible on localhost (:80) (if you are on windows and use docker-toolbox, use the command >docker-machine ls to get the URL of you VM)

| | | |
| ------------- | ------------- | ------------- |
| __Method__		| GET 									|																	|
| __URL__  			| /shortlink/my.simpleURL.com 			| where	my.simpleURL.com is the URL you want to get a shortlink for	|
| __param__  		| custom  (optional) 					| value under wich you want to save this URL 						|
| __returns__  		| a Json structure containing the token and the origin	|																	|
| __example__  		| /shortlink/www.google.com?custom=g  	| {"Origin":"http://www.google.com","Shortlinkken":"http://localhost/g"}				 |

Since encoding a URL to create a shortlink would break the simplicity, for "complex" URL to encode (like https://github.com/tools/godep), please use the following POST method 

| | | |
| ------------- | ------------- | ------------- |
| __Method__		| POST 									|																	|
| __URL__  			| /shortlink				 			|  																	|
| __body__  		| a origin-token tuple  				| where origin is the target of the link and token the value under wich you want to save this URL |
| __returns__  		| a Json structure containing the token and the origin	|																	|
| __example__  		| /shortlink  body={"origin":"https://github.com/tools/godep","token":"godep"}| {"Origin":"https://github.com/tools/godepII","Shortlink":"http://localhost/godep"} |

| | | |
| ------------- | ------------- | ------------- |
| __Method__		| GET 										|											|
| __URL__  			| /shtlnk 									| where	shtlnk is the shortlink token		|
| __returns__  		| nothing, you are redirected to the target	|											|
| __example__  		| /g 										| you are now on  	http://www.google.com	|

| | | |
| ------------- | ------------- | ------------- |
| __Method__		| GET 										|											|
| __URL__  			| /admin/shtlnk 							| where	shtlnk is the shortlink token		|
| __returns__  		| the Shortlink object stored in redis		|											|
| __example__  		| /admin/shtlnk								| {"Id":"7e3b6d7d-dbeb-44ba-839b-80442accef55","Token":"g","Origin":"http://www.google.com","CreationTs":1447263537,"Count":1}	|

