package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "172.17.0.1"
	port     = 5434
	user     = "postgres"
	password = "example"
	dbname   = "mtg"
)

var db *sql.DB

type Karta struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Colors       []string `json:"colors"`
	Cmc          float64  `json:"cmc"`
	Type         string   `json:"type"`
	Types        []string `json:"types"`
	Supertypes   []string `json:"supertypes"`
	Subtypes     []string `json:"subtypes"`
	Rarity       string   `json:"rarity"`
	ImageUrl     string   `json:"imageUrl"`
	OriginalText string   `json:"originalText"`
}

type allCards struct {
	SveKarte []Karta `json:"cards"`
}

var cards allCards

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	router := gin.Default()
	router.GET("/import/:page", Import)
	router.GET("/card/:info", Info)
	router.GET("/list", Search)
	router.Run("localhost:9000")
}
func Search(c *gin.Context) {
	var id string
	var name string
	var colors string
	var cmc int
	var tip string
	var types string
	var supertypes string
	var subtypes string
	var rarity string
	var imageUrl string
	var originalText string
	var ret Karta

	var conditions struct {
		Condition []string
		Value     []interface{}
	}
	paramPairs := c.Request.URL.Query()

	for k, v := range paramPairs {
		conditions.Condition = append(conditions.Condition, k)
		conditions.Value = append(conditions.Value, v[0])
	}
	if len(conditions.Condition) != len(conditions.Value) {
		c.String(404, "wrong params")
		return
	}
	switch len(conditions.Condition) {
	case 1:
		stri := fmt.Sprintf(`SELECT * FROM "mtg" WHERE "%s"=$1`, conditions.Condition[0])
		//rows, _ := db.Query(`SELECT * FROM "mtg" WHERE col = "$1"  =$2`, conditions.Condition, conditions.Value[0])
		rows, _ := db.Query(stri, conditions.Value[0])
		for rows.Next() {
			_ = rows.Scan(&id, &name, &colors, &cmc, &tip, &types, &supertypes, &subtypes, &rarity, &imageUrl, &originalText)
			ret.Id = id
			ret.Name = name
			ret.Colors = append(ret.Colors, colors)
			ret.Cmc = float64(cmc)
			ret.Type = tip
			ret.Types = append(ret.Types, types)
			ret.Supertypes = append(ret.Supertypes, supertypes)
			ret.Subtypes = append(ret.Subtypes, subtypes)
			ret.Rarity = rarity
			ret.ImageUrl = imageUrl
			ret.OriginalText = originalText
			retur, _ := json.MarshalIndent(ret, "", "")
			c.String(200, string(retur))
			ret.Types = ret.Types[:0]
			ret.Colors = ret.Types[:0]
			ret.Supertypes = ret.Types[:0]
			ret.Subtypes = ret.Types[:0]
		}

	}

}

func Import(c *gin.Context) {
	var var3 string
	var var6 string
	var var7 string
	var var8 string
	page := c.Param("page")
	response, _ := http.Get("https://api.magicthegathering.io/v1/cards?page=" + page)
	json.NewDecoder(response.Body).Decode(&cards)
	for i := 0; i < len(cards.SveKarte); i++ {
		var1 := cards.SveKarte[i].Id
		var2 := cards.SveKarte[i].Name
		for j := range cards.SveKarte[i].Colors {
			var3 = fmt.Sprintf("%s%s", var3, cards.SveKarte[i].Colors[j])
		}
		var4 := cards.SveKarte[i].Cmc
		var5 := cards.SveKarte[i].Type
		for k := range cards.SveKarte[i].Types {
			var6 = fmt.Sprintf("%s%s", var6, cards.SveKarte[i].Types[k])
		}
		for l := range cards.SveKarte[i].Supertypes {
			var7 = fmt.Sprintf("%s%s", var7, cards.SveKarte[i].Supertypes[l])
		}
		for m := range cards.SveKarte[i].Subtypes {
			var8 = fmt.Sprintf("%s%s", var8, cards.SveKarte[i].Subtypes[m])
		}
		var9 := cards.SveKarte[i].Rarity
		var10 := cards.SveKarte[i].ImageUrl
		var11 := cards.SveKarte[i].OriginalText
		if _, err := db.Query("insert into mtg values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", var1, var2, var3, var4, var5, var6, var7, var8, var9, var10, var11); err != nil {
			fmt.Println(err)
			return
		}
		var3 = ""
		var6 = ""
		var7 = ""
		var8 = ""
	}
	c.String(200, "cards imported")
}

func Info(c *gin.Context) {
	var id string
	var name string
	var colors string
	var cmc int
	var tip string
	var types string
	var supertypes string
	var subtypes string
	var rarity string
	var imageUrl string
	var originalText string
	var ret Karta
	info := c.Param("info")
	row := db.QueryRow(`SELECT * FROM "mtg" WHERE "Id"=$1`, info)
	_ = row.Scan(&id, &name, &colors, &cmc, &tip, &types, &supertypes, &subtypes, &rarity, &imageUrl, &originalText)
	ret.Id = id
	ret.Name = name
	ret.Colors = append(ret.Colors, colors)
	ret.Cmc = float64(cmc)
	ret.Type = tip
	ret.Types = append(ret.Types, types)
	ret.Supertypes = append(ret.Supertypes, supertypes)
	ret.Subtypes = append(ret.Subtypes, subtypes)
	ret.Rarity = rarity
	ret.ImageUrl = imageUrl
	ret.OriginalText = originalText
	retur, _ := json.MarshalIndent(ret, "", "")

	c.String(200, string(retur))
}
