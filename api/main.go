package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// for API json
type Section struct {
	RenameTo string       `json:"rename_to"`
	Section  string       `json:"section"`
	Items    SectionItems `json:"items"`
}

type SectionItems map[string]SectionItem

type SectionItem struct {
	RenameTo    string `json:"rename_to"`
	TranslateTo string `json:"translate_to"`
}

// for database
type TranslationItems map[string]TranslationItem
type TranslationItem map[string]string
type Translation struct {
	Section string           `bson:"section"`
	Items   TranslationItems `bson:"items"`
}

var dbhost = flag.String("dbhost", "", "Ex: localhost:27017")
var dbname = flag.String("dbname", "", "Ex: cnv3")
var host = flag.String("host", ":8080", "Ex: localhost:8080")

func main() {
	flag.Parse()

	if *dbhost == "" || *dbname == "" {
		fmt.Println(errors.New("Err: dbhost or dbname can't be blank"))
		flag.Usage()
		return
	}

	session, err := mgo.Dial(*dbhost)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	db := session.DB(*dbname)

	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Origin", "*")
	})

	r.GET("/translation/:lang", func(ctx *gin.Context) {
		iter := db.C("i18n").Find(bson.M{}).Iter()
		lang := ctx.Param("lang")

		secArr := make([]Section, 0, 10)
		trans := new(Translation)
		for iter.Next(trans) {
			items := make(SectionItems)
			for original, ts := range trans.Items {
				for langName, translatedTo := range ts {
					if langName == lang {
						items[original] = SectionItem{
							TranslateTo: translatedTo,
						}
						break
					}
				}
			}

			secArr = append(secArr, Section{
				Section: trans.Section,
				Items:   items,
			})
		}

		ctx.JSON(200, gin.H{
			"result": true,
			"data":   secArr,
		})
	})

	r.POST("/translation/:to_lang/:section", func(ctx *gin.Context) {
		var err error
		section := ctx.Param("section")
		toLang := ctx.Param("to_lang")
		sec := new(Section)
		err = ctx.BindJSON(sec)

		if err != nil {
			ctx.JSON(400, gin.H{
				"result": false,
				"err":    err.Error(),
			})

			return
		}

		trans := new(Translation)
		err = db.C("i18n").Find(bson.M{"section": section}).One(trans)
		if err != nil {
			ctx.JSON(400, gin.H{
				"result": false,
				"err":    "section does not exist",
			})

			return
		}

		// new section name
		if sec.RenameTo != "" {
			trans.Section = sec.RenameTo
		}

		//check any removed item
		for o, _ := range trans.Items {
			if _, ok := sec.Items[o]; !ok {
				delete(trans.Items, o)
			}
		}

		// check any new item or translated language
		for o, it := range sec.Items {
			newOriginalItem := true
			for o2, _ := range trans.Items {
				if o == o2 {

					// rename its key
					if it.RenameTo != "" {
						trans.Items[it.RenameTo] = trans.Items[o]
						delete(trans.Items, o)
						trans.Items[it.RenameTo][toLang] = it.TranslateTo
					} else {
						trans.Items[o][toLang] = it.TranslateTo
					}

					newOriginalItem = false
					break
				}
			}

			if newOriginalItem {
				if _, ok := trans.Items[o]; !ok {
					trans.Items[o] = make(TranslationItem)
				}

				trans.Items[o][toLang] = it.TranslateTo
			}
		}

		info, err := db.C("i18n").Upsert(bson.M{"section": section}, trans)
		if info.Updated == 0 || err != nil {
			ctx.JSON(400, gin.H{
				"result": false,
				"err":    err.Error(),
			})

			return
		}

		ctx.JSON(200, gin.H{
			"result": true,
		})
	})

	r.OPTIONS("*path", func(ctx *gin.Context) {
		//nothing
	})

	r.Run(*host)
}
