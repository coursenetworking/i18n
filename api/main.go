package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
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
	Section string           `json:"section"`
	Items   TranslationItems `json:"items"`
}

// the db Object, it reflect to the dbfile
type dbfileHandler []Translation

func (d *dbfileHandler) Section(name string, trans *Translation) error {
	for _, i := range *d {
		if i.Section == name {
			*trans = i
			return nil
		}
	}

	return errors.New("section does not exist")
}

func (d *dbfileHandler) Append(trans Translation) error {
	if trans.Section == "" {
		return errors.New("section can't be blank")
	}

	*d = append(*d, trans)
	return nil
}

func (d *dbfileHandler) Update(name string, trans Translation) error {
	for k, i := range *d {
		if i.Section == name {
			(*d)[k] = trans
			return nil
		}
	}

	return errors.New("section does not exist")
}

func newTranslation() *Translation {
	trans := new(Translation)
	trans.Items = make(TranslationItems)
	return trans
}

var host = flag.String("host", ":8080", "Ex: localhost:8080")
var dbfile = flag.String("dbfile", "tmp/db.json", "the file to store translation data")

func main() {
	flag.Parse()

	if *dbfile == "" {
		fmt.Println(errors.New("Err: dbfile can't be blank"))
		flag.Usage()
		return
	}

	dbf, err := os.OpenFile(*dbfile, os.O_RDWR|os.O_CREATE, 0700)
	if err != nil {
		panic(err.Error())
	}
	defer dbf.Close()

	dbTranslations := new(dbfileHandler)
	finfo, err := dbf.Stat()
	if err != nil {
		panic(err.Error())
	}

	if finfo.Size() != 0 {
		err = json.NewDecoder(dbf).Decode(dbTranslations)
		if err != nil {
			panic(err.Error())
		}
	}

	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Origin", "*")
	})

	// Dump data data output
	r.GET("/translation", func(ctx *gin.Context) {
		ctx.JSON(200, dbTranslations)
	})

	// Retrun th given language translation
	r.GET("/translation/:lang", func(ctx *gin.Context) {
		lang := ctx.Param("lang")

		secArr := make([]Section, 0, 10)

		for _, trans := range *dbTranslations {
			items := make(SectionItems)
			for original, ts := range trans.Items {
				find := false
				for langName, translatedTo := range ts {
					if langName == lang {
						find = true
						items[original] = SectionItem{
							TranslateTo: translatedTo,
						}
						break
					}
				}

				if !find {
					items[original] = SectionItem{
						TranslateTo: "",
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

	// Create/Update section translation in given language
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

		trans := newTranslation()
		err = dbTranslations.Section(section, trans)
		isNewSection := err != nil

		//new section name
		if sec.RenameTo != "" {
			trans.Section = sec.RenameTo
		}

		//check any removed item
		for o, _ := range trans.Items {
			if _, ok := sec.Items[o]; !ok {
				delete(trans.Items, o)
			}
		}

		//check any new item or translated language
		for o, it := range sec.Items {
			newOriginalItem := true
			for o2, _ := range trans.Items {
				if o == o2 {
					//rename its key
					if it.RenameTo != "" && it.RenameTo != o {
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

		if isNewSection {
			trans.Section = section
			err = dbTranslations.Append(*trans)
		} else {
			err = dbTranslations.Update(section, *trans)
		}

		if err != nil {
			ctx.JSON(300, gin.H{
				"result": false,
				"err":    err.Error(),
			})
		}

		dbf.Seek(0, 0)
		err = json.NewEncoder(dbf).Encode(dbTranslations)
		if err != nil {
			ctx.JSON(300, gin.H{
				"result": false,
				"err":    err.Error(),
			})
		}

		ctx.JSON(200, gin.H{
			"data":   dbTranslations,
			"result": true,
		})
	})

	r.OPTIONS("*path", func(ctx *gin.Context) {
		//nothing
	})

	r.Run(*host)
}
