package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

type Card struct {
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
	CardSlice []Card `json:"cards"`
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
	_, err = db.Query(`SELECT * FROM "mtg"`)
	if err != nil {
		CreateTable()
	}

	router := gin.Default()
	router.GET("/import/:page", Import)
	router.GET("/card/:info", Info)
	router.GET("/list", Search)
	router.Run("localhost:9001")
}
func CreateTable() {
	db.Exec("CREATE TABLE $1", "mtg")
	db.Exec(`ALTER TABLE "mtg" ADD Id text`)
	db.Exec(`ALTER TABLE "mtg" ADD Name text`)
	db.Exec(`ALTER TABLE "mtg" ADD Colors text`)
	db.Exec(`ALTER TABLE "mtg" ADD Cmc integer`)
	db.Exec(`ALTER TABLE "mtg" ADD Type text`)
	db.Exec(`ALTER TABLE "mtg" ADD Types text`)
	db.Exec(`ALTER TABLE "mtg" ADD Supertypes text`)
	db.Exec(`ALTER TABLE "mtg" ADD Subtypes text`)
	db.Exec(`ALTER TABLE "mtg" ADD Rarity text`)
	db.Exec(`ALTER TABLE "mtg" ADD ImageUrl text`)
	db.Exec(`ALTER TABLE "mtg" ADD OriginalText text`)
}

func getData(c *gin.Context, rows *sql.Rows, total int, pageNumber int) {
	var id string
	var name string

	type Cards struct {
		Id   string `json:"Id"`
		Name string `json:"Name"`
	}

	type ret struct {
		Total     int     `json:"Total:"`
		Page      int     `json:"Page:"`
		Items     int     `json:"Items"`
		CardSlice []Cards `json:"Cards:"`
	}
	var retAll ret

	var counter int
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			c.String(404, "Error")
			return
		}
		retAll.CardSlice[counter].Id = id
		retAll.CardSlice[counter].Name = name
	}
	retAll.Total = total
	retAll.Page = pageNumber
	retAll.Items = counter
	retur, _ := json.MarshalIndent(retAll, "", "")
	c.String(200, string(retur))

}

func Search(c *gin.Context) {

	var conditions struct {
		Condition []string
		Value     []interface{}
	}
	paramPairs := c.Request.URL.Query()

	for k, v := range paramPairs {
		conditions.Condition = append(conditions.Condition, k)
		conditions.Value = append(conditions.Value, v[0])
	}
	if len(conditions.Condition) == 1 && conditions.Condition[0] == "page" {
		pageNumber := conditions.Value[0]
		pageNum := pageNumber.(int)
		pgNum := (pageNum - 1) * 10
		rows, err := db.Query(`SELECT "Id", "Name" FROM "mtg" LIMIT 10 OFFSET $1`, pgNum)
		if err != nil {
			c.String(404, "Please check if params are correct")
			return
		}
		defer rows.Close()
		getData(c, rows, 0, pageNum)
		return

	}
	pageNumber := conditions.Value[len(conditions.Value)-1]
	pageNum := pageNumber.(int)
	pgNum := (pageNum - 1) * 10

	conditions.Condition = conditions.Condition[:len(conditions.Condition)-1]
	conditions.Value = conditions.Value[:len(conditions.Value)-1]
	query := fmt.Sprintf(`SELECT "Id", "Name" FROM "mtg" WHERE "%s"=$1`, conditions.Condition[0])
	for i := range conditions.Condition {
		if i == 0 {
			continue
		}
		query = fmt.Sprintf(`%s AND "%s"=$%d`, query, conditions.Condition[i], i+1)
	}
	query = fmt.Sprintf(`%s LIMIT 10 OFFSET %s`, query, strconv.Itoa(pgNum))
	rows, err := db.Query(query, conditions.Value...)
	if err != nil {
		c.String(404, "Please check if params are correct")
		return
	}
	defer rows.Close()

	getData(c, rows, 0, pageNum)

}

func Import(c *gin.Context) {
	var var3 string
	var var6 string
	var var7 string
	var var8 string
	page := c.Param("page")
	response, err := http.Get("https://api.magicthegathering.io/v1/cards?page=" + page)
	if err != nil {
		c.String(404, "Please check if url and page number are correct")
		return
	}
	json.NewDecoder(response.Body).Decode(&cards)
	for i := 0; i < len(cards.CardSlice); i++ {
		var1 := cards.CardSlice[i].Id
		var2 := cards.CardSlice[i].Name
		for j := range cards.CardSlice[i].Colors {
			var3 = fmt.Sprintf("%s%s", var3, cards.CardSlice[i].Colors[j])
		}
		var4 := cards.CardSlice[i].Cmc
		var5 := cards.CardSlice[i].Type
		for k := range cards.CardSlice[i].Types {
			var6 = fmt.Sprintf("%s%s", var6, cards.CardSlice[i].Types[k])
		}
		for l := range cards.CardSlice[i].Supertypes {
			var7 = fmt.Sprintf("%s%s", var7, cards.CardSlice[i].Supertypes[l])
		}
		for m := range cards.CardSlice[i].Subtypes {
			var8 = fmt.Sprintf("%s%s", var8, cards.CardSlice[i].Subtypes[m])
		}
		var9 := cards.CardSlice[i].Rarity
		var10 := cards.CardSlice[i].ImageUrl
		var11 := cards.CardSlice[i].OriginalText
		if _, err := db.Query("insert into mtg values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", var1, var2, var3, var4, var5, var6, var7, var8, var9, var10, var11); err != nil {
			c.String(404, "Not found")
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
	var cmc float64
	var tip string
	var types string
	var supertypes string
	var subtypes string
	var rarity string
	var imageUrl string
	var originalText string

	var ret Card
	info := c.Param("info")
	row := db.QueryRow(`SELECT * FROM "mtg" WHERE "Id"=$1`, info)
	err := row.Scan(&id, &name, &imageUrl)
	if err != nil {
		c.String(404, "Card is not in database or id does not exist")
		return
	}
	ret.Id = id
	ret.Name = name
	ret.Colors[0] = colors
	ret.Cmc = cmc
	ret.Type = tip
	ret.Types[0] = types
	ret.Supertypes[0] = supertypes
	ret.Subtypes[0] = subtypes
	ret.Rarity = rarity
	ret.ImageUrl = imageUrl
	ret.OriginalText = originalText
	retur, err := json.MarshalIndent(ret, "", "")
	if err != nil {
		c.String(404, "Not found")
		return
	}

	c.String(200, string(retur))
}
