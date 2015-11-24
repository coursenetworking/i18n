# CN i18n

For collecting different languages translation

#Database structure
	{
	    "section": "header",
	    "items": {
	        "post": {
		         "zh-cn": "文章",
	            "<lang name>": "",
	        },
	        "poll": {
	            ...
	        },
	        ...
		 }
	}

# API list
## [POST] /translation/:lang/:section
	{
	     "items": {
	         "post":   "文章",
	         "search": "搜索",
	     },
	}

## [GET] /translation/:lang

	[
		{
	     	"section": "header",
	     	"to_lang": "zh-CN",
	     	"items": {
	         	"post":   "文章",
	         	"search": "搜索",
	     	},
		}, 
		...
	]
	
# Run API
You can build this API by running `build.sh` or just run the program from `bin` folder.

For development:
	
	./cni18n --dbhost=<host name> --dbname=<db name>

For production:

	export GIN_MODE=release && ./i18n-api --dbhost=<host name> --dbname=<db name>