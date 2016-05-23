package main

import (
	"encoding/json"
	"os"
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
}

type tempFileType map[string][]string

type dbtempType map[string]TranslationItems

func inArray(item string, arr []string) bool {
	for _, v := range arr {
		if item == v {
			return true
		}
	}

	return false
}

func main() {

	var err error
	tempFile, err := os.OpenFile("t.json", os.O_RDONLY, 0400)
	if err != nil {
		panic(err)
	}
	defer tempFile.Close()

	dbfile, err := os.OpenFile("db.json", os.O_RDONLY, 0400)
	if err != nil {
		panic(err)
	}
	defer dbfile.Close()

	var tempf tempFileType
	err = json.NewDecoder(tempFile).Decode(&tempf)
	if err != nil {
		panic(err)
	}

	var dbf []Translation
	err = json.NewDecoder(dbfile).Decode(&dbf)
	if err != nil {
		panic(err)
	}

	dbtemp := make(dbtempType)
	for _, f := range dbf {
		dbtemp[f.Section] = f.Items
	}

	for temp_section, temp_items := range tempf {
		// old db has this section, update it
		if _, ok := dbtemp[temp_section]; ok {
			for db_item_text, _ := range dbtemp[temp_section] {
				if !inArray(db_item_text, temp_items) {
					delete(dbtemp[temp_section], db_item_text)
				}
			}

			for _, temp_item_text := range temp_items {
				if _, ok := dbtemp[temp_section][temp_item_text]; !ok {
					dbtemp[temp_section][temp_item_text] = make(TranslationItem)
				}
			}

		} else {
			dbtemp[temp_section] = make(TranslationItems)
			for _, temp_item := range temp_items {
				dbtemp[temp_section][temp_item] = make(TranslationItem)
			}
		}
	}

	newdb := make([]Translation, 0)
	for dbtemp_section_name, items := range dbtemp {
		newdb = append(newdb, Translation{
			Section: dbtemp_section_name,
			Items:   items,
		})
	}

	newdbFile, err := os.OpenFile("new_db.json", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0700)
	if err != nil {
		panic(err)
	}

	defer newdbFile.Close()
	err = json.NewEncoder(newdbFile).Encode(newdb)
	if err != nil {
		panic(err)
	}
}
