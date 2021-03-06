package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/contrib/gzip"
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
type dbfileHandler struct {
	dbfile     *os.File
	collection []Translation
	lock       sync.RWMutex
}

func (h *dbfileHandler) Collection() []Translation {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.collection
}

func (h *dbfileHandler) Section(name string, trans *Translation) error {
	h.lock.RLock()
	defer h.lock.RUnlock()

	for _, i := range h.collection {
		if i.Section == name {
			*trans = i
			return nil
		}
	}

	return errors.New("section does not exist")
}

func (h *dbfileHandler) Append(trans Translation) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if trans.Section == "" {
		return errors.New("section can't be blank")
	}

	h.collection = append(h.collection, trans)
	return h.sync()
}

func (h *dbfileHandler) Update(name string, trans Translation) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	for k, i := range h.collection {
		if i.Section == name {
			h.collection[k] = trans
			return h.sync()
		}
	}

	return errors.New("section does not exist")
}

func (h *dbfileHandler) sync() error {
	h.dbfile.Seek(0, 0)
	b, err := json.MarshalIndent(h.collection, "", " ")
	if err != nil {
		return err
	}

	err = h.dbfile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = h.dbfile.Write(b)
	return err
}

func (h *dbfileHandler) Close() {
	h.dbfile.Close()
}

func newDbfileHandler(f *os.File) *dbfileHandler {
	h := new(dbfileHandler)
	h.dbfile = f
	h.collection = make([]Translation, 0, 100)

	finfo, err := h.dbfile.Stat()
	if err != nil {
		panic(err.Error())
	}

	if finfo.Size() != 0 {
		err = json.NewDecoder(h.dbfile).Decode(&h.collection)
		if err != nil {
			panic(err.Error())
		}
	}

	return h
}

func newTranslation() *Translation {
	trans := new(Translation)
	trans.Items = make(TranslationItems)
	return trans
}

// convert the database one section item to API json section item structure
func toSectionStruct(trans *Translation, lang string) Section {
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

	return Section{
		Section: trans.Section,
		Items:   items,
	}
}

var debug = flag.Bool("debug", false, "debug mode or not")
var host = flag.String("host", ":8080", "Ex: localhost:8080")
var dbfile = flag.String("dbfile", "", "the file to store translation data")
var basePath = flag.String("basepath", ".", "the base runtime path")

var requestLocker sync.Mutex

func main() {
	flag.Parse()

	os.Chdir(*basePath)

	if *dbfile == "" {
		fmt.Println(errors.New("Err: dbfile can't be blank"))
		flag.Usage()
		return
	}

	dbf, err := os.OpenFile(*dbfile, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0700)
	if err != nil {
		panic(err.Error())
	}

	dbh := newDbfileHandler(dbf)
	defer dbh.Close()

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	apiHeader := func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Origin", "*")
	}

	// Dump data data output
	r.GET("/db", apiHeader, func(ctx *gin.Context) {
		ctx.JSON(200, dbh.Collection())
	})

	// Retrun th given language translation
	r.GET("/translation/:lang", apiHeader, func(ctx *gin.Context) {

		lang := ctx.Param("lang")
		secArr := make([]Section, 0, 10)

		for _, trans := range dbh.Collection() {
			secArr = append(secArr, toSectionStruct(&trans, lang))
		}

		ctx.JSON(200, gin.H{
			"result": true,
			"data":   secArr,
		})
	})

	r.GET("/translation/:lang/:section", apiHeader, func(ctx *gin.Context) {
		lang := ctx.Param("lang")
		sectionName := ctx.Param("section")

		trans := newTranslation()
		err := dbh.Section(sectionName, trans)
		if err != nil {
			ctx.JSON(404, gin.H{
				"result": false,
				"err":    err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"result": true,
			"data":   toSectionStruct(trans, lang),
		})
	})

	// Create/Update section translation in given language
	r.POST("/translation/:to_lang/:section", apiHeader, func(ctx *gin.Context) {

		var err error
		section := ctx.Param("section")
		toLang := ctx.Param("to_lang")
		isAdmin := ctx.Param("is_admin") == "1"

		requestLocker.Lock()
		defer requestLocker.Unlock()

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
		err = dbh.Section(section, trans)
		isNewSection := err != nil

		outdateErr := "Fail to save your change, because your page is expired, the language database has been updated, please fresh the page."

		if !isAdmin && isNewSection {
			ctx.JSON(200, gin.H{
				"result": false,
				"err":    outdateErr,
			})
			return
		}

		//new section name
		if sec.RenameTo != "" {
			trans.Section = sec.RenameTo
		}

		//check any removed item
		for o, _ := range trans.Items {
			if _, ok := sec.Items[o]; !ok {
				if isAdmin {
					delete(trans.Items, o)
				} else {
					ctx.JSON(200, gin.H{
						"result": false,
						"err":    outdateErr,
					})
					return
				}
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
				if isAdmin {
					if _, ok := trans.Items[o]; !ok {
						trans.Items[o] = make(TranslationItem)
					}

					trans.Items[o][toLang] = it.TranslateTo
				} else {
					ctx.JSON(200, gin.H{
						"result": false,
						"err":    outdateErr,
					})
					return
				}
			}
		}

		if isNewSection {
			trans.Section = section
			err = dbh.Append(*trans)
		} else {
			err = dbh.Update(section, *trans)
		}

		if err != nil {
			ctx.JSON(200, gin.H{
				"result": false,
				"err":    err.Error(),
			})
		}

		ctx.JSON(200, gin.H{
			"result": true,
		})
	})

	r.OPTIONS("*path", apiHeader, func(ctx *gin.Context) {
		//nothing
	})

	if *debug {
		r.StaticFile("/", "static/src/index.html")

		r.Static("/js", "static/src/js")
		r.Static("/img", "static/src/img")
		r.Static("/css", "static/src/css")
		r.Static("/bower_components", "static/src/bower_components")
		r.Static("/tpl", "static/src/tpl")
	} else {
		r.GET("/static/*path", getAsset)
		r.GET("/", getHome)
	}

	r.Run(*host)
}

func getHome(ctx *gin.Context) {
	serveStaticAsset("/index.html", ctx)
}

func getAsset(ctx *gin.Context) {
	serveStaticAsset(ctx.Params.ByName("path"), ctx)
}

func serveStaticAsset(path string, ctx *gin.Context) {
	data, err := Asset("static/dist" + path)
	if err != nil {
		ctx.String(400, err.Error())
		return
	}

	ctx.Data(200, assetContentType(path), data)
}

var extraMimeTypes = map[string]string{
	".icon": "image-x-icon",
	".ttf":  "application/x-font-ttf",
	".woff": "application/x-font-woff",
	".eot":  "application/vnd.ms-fontobject",
	".svg":  "image/svg+xml",
	".html": "text/html; charset-utf-8",
}

func assetContentType(name string) string {
	ext := filepath.Ext(name)
	result := mime.TypeByExtension(ext)

	if result == "" {
		result = extraMimeTypes[ext]
	}

	if result == "" {
		result = "text/plain; charset=utf-8"
	}

	return result
}
