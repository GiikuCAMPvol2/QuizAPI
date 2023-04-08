package controller

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"server/config"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/murasame29/rakutenapi"
)

type ReturnData struct {
	Quiz          string             `json:"quiz"`
	Answer        int                `json:"answer"`
	Images        []rakutenapi.Image `json:"images"`
	AffiliateLink string             `json:"affiliatelink"`
}

func readCsv() [][]string {
	file, err := os.Open("rakuten_genre.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return rows
}

// 文字列をUint64に変換する
func convertUint(x string) (uint64, error) {
	if len(x) != 0 {
		data, err := strconv.ParseUint(x, 10, 64)
		if err != nil {
			return data, fmt.Errorf("符号なしの整数値のみの入力 :%s", err)
		}
		return data, nil
	} else {
		return 1, nil
	}
}

// ランダムなジャンルIDを返す
func randomGenre() []string {
	rand.Seed(time.Now().UnixNano())
	RakutenGenre := readCsv()
	return RakutenGenre[rand.Intn(len(RakutenGenre)-1)]
}

// GetQuizされた時の関数
func GetQuiz(c *gin.Context) {
	hits, err := convertUint(c.Query("hits"))

	if err != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("wrong argument type :%s", err),
		})
		return
	}

	page, err := convertUint(c.Query("page"))

	if err != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("wrong argument type :%s", err),
		})
		return
	}

	sort, err := convertUint(c.Query("sort"))

	if err != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("wrong argument type :%s", err),
		})
		return
	}

	genreId := c.Query("genreid")

	if len(c.Query("keyword")) == 0 && len(c.Query("genreid")) == 0 {
		genreId = randomGenre()[0]
	}

	requestData := rakutenapi.RakutenSearch{
		S:        config.InitRequest(),
		Keyword_: c.Query("keyword"),
		GenreId_: genreId,
		Hits_:    hits,
		Page_:    page,
		Sort_:    sort,
	}

	responseData, err := requestData.Do()

	if err != nil {
		c.JSON(400, gin.H{
			"message": fmt.Sprintf("wrong argument type :%s", err),
		})
		return
	}

	for i := range responseData.Items[0].Item.MediumImageUrls {
		responseData.Items[0].Item.MediumImageUrls[i].ImageURL = strings.ReplaceAll(responseData.Items[0].Item.MediumImageUrls[i].ImageURL, "?_ex=128x128", "")
	}
	ReturnData := ReturnData{
		Quiz:          responseData.Items[0].Item.ItemName,
		Answer:        responseData.Items[0].Item.ItemPrice,
		Images:        responseData.Items[0].Item.MediumImageUrls,
		AffiliateLink: responseData.Items[0].Item.AffiliateURL,
	}
	c.JSON(200, ReturnData)
}
