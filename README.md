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
Translate the `section` to another language(`:lang`)

	{
	     "section": "header", // API will ignore this value, because the section name is in the URL already
	     "rename_to": "", // default empty to keep this section name, if set, will rename the section name to this new value.
	     "items": {
	         "post": {
	             "rename_to": "", //default empty to keep this key "post"
	             "translate_to": "文章"
	         },
	         "search": {
	             "rename_to": "Search", // will rename "search" to "Search"
	             "translate_to": "搜索"
	         }
	     }
	}

## [GET] /translation/:lang
Get the translation of given langauge(`:lang`)

	[
		{
		     "section": "header", // the name of this section
		     "rename_to": "", // can ignore here
		     "items": {
                "post": {
                    "rename_to": "", // can ignore here
                    "translate_to": "文章"
                },
                "search": {
                    "rename_to": "",
                    "translate_to": "搜索"
                }
            }
       },
		...
	]

# Run API
You can build this API by running `build.sh` or just run the program from `bin` folder.

For development:
	
	./cni18n --dbfile=tmp/db.json

For production:

	export GIN_MODE=release && ./i18n-api --dbfile=data/db.json
