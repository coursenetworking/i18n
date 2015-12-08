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
type SectionTranslation struct {
	Section string            `json:"section"`
	ToLang  string            `json:"to_lang"`
	Items   map[string]string `json:"items"`
}

// for database
type TranslationItem map[string]map[string]string

type Translation struct {
	Section string          `bson:"section"`
	Items   TranslationItem `bson:"items"`
}

// Create new English item
type Items struct {
	Items []string `json:"items" bson:"items"`
}

var dbhost = flag.String("dbhost", "", "Ex: localhost:27017")
var dbname = flag.String("dbname", "", "Ex: cnv3")
var host = flag.String("host", ":8080", "Ex: localhost:8080")

func main() {
	flag.Parse()

	if *dbhost == "" || *dbname == "" {
		fmt.Println(errors.New("Err: dbhost or dbname can't be blank"))
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
		trans := new(Translation)
		lang := ctx.Param("lang")

		secArr := make([]SectionTranslation, 0, 10)
		for iter.Next(trans) {
			items := make(map[string]string)
			for original, ts := range trans.Items {
				items[original] = ""
				for langName, translatedLang := range ts {
					if langName == lang {
						items[original] = translatedLang
						break
					}
				}
			}

			secArr = append(secArr, SectionTranslation{
				Section: trans.Section,
				ToLang:  lang,
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
		sec := new(SectionTranslation)
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

		for original, _ := range trans.Items {
			for original2, translated := range sec.Items {
				if original == original2 {
					trans.Items[original][toLang] = translated
					break
				}
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

	r.POST("/source/:section", func(ctx *gin.Context) {
		section := ctx.Param("section")
		n, err := db.C("i18n").Find(bson.M{"section": section}).Count()
		if n == 1 {
			ctx.JSON(400, gin.H{
				"result": false,
				"err":    "section exists already",
			})
			return
		} else if err != nil {
			ctx.JSON(400, gin.H{
				"result": false,
				"err":    err.Error(),
			})
			return
		}

		sec := new(Items)
		err = ctx.BindJSON(sec)

		if err != nil {
			ctx.JSON(400, gin.H{
				"result": false,
				"err":    err.Error(),
			})

			return
		}

		if 0 == len(sec.Items) {
			ctx.JSON(400, gin.H{
				"result": false,
				"err":    "items can't be empty",
			})

			return
		}

		trans := new(Translation)
		trans.Items = TranslationItem{}
		trans.Section = section
		for _, v := range sec.Items {
			trans.Items[v] = make(map[string]string)
		}

		err = db.C("i18n").Insert(trans)
		if err != nil {
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
