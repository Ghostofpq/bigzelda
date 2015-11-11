# bigzelda
BigZelda is a shortlink creator (get it? short link?)

To run BigZelda, install docker.

Then go in the BigZelda directory and use
./install
and
./start

This should run in your docker 2 containers
- redis 		: 	a simple redis-server
- bigzelda		: 	the bigzelda go application


BigZelda exposes a simple API on the port 6060. To change this port, modify the command line in start.sh (--publish >6060<:8000)
This API is accessible on localhost:6060
| First Header  | Second Header 						| Second Header 														|
| ------------- | ------------------------------------- |------------------------------------------------------------------ |
| Method		| GET 									|																	|
| URL  			| /shortlink/my.simpleURL.com 			| where	my.simpleURL.com is the URL you want to get a shortlink for	|
| param  		| custom  								| token under wich you want to save this URL 						|
| returns  		| a message indicating the redirection	|																	|
| example  		| /shortlink/www.google.com?custom=g  	| http://www.google.com is now accessible via /g 					|