package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
	"github.com/xormplus/xorm"
	"github.com/xormplus/core"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"time"
)

var engine *xorm.Engine

var Config = struct {
	DB struct {
		   Host     string
		   User     string `default:"root"`
		   Password string
		   Port     string `default:"3306"`
		   Database string
		   Colum    string
	   }
}{}

type akb_diff struct {
	Author string
	Diff   int64
	Href   string
	Title  string
	Time   int64
}

var users []akb_diff

func main() {
	configor.Load(&Config, "./config.yml")
	sql := "SELECT DISTINCT title FROM "+Config.DB.Colum
	engine, _ = xorm.NewEngine("mysql", Config.DB.User + ":" + Config.DB.Password + "@tcp(" + Config.DB.Host + ":" + Config.DB.Port + ")/" + Config.DB.Database)
	engine.Logger().SetLevel(core.LOG_DEBUG)
	Alltitle, _ := engine.Sql(sql).Query().Json()
	gjson.Parse(Alltitle).ForEach(func(key, value gjson.Result) bool {
		replynum := make([]int64, 0)
		title := value.Get("title").String()
		sql = "SELECT replynum,time,title,href,author FROM "+Config.DB.Colum+" WHERE title=? "
		ReplyNumList, _ := engine.Sql(sql, title).Query().Json()
		gjson.Parse(ReplyNumList).ForEach(func(key, value gjson.Result) bool {
			replynum = append(replynum, value.Get("replynum").Int())
			return true
		})
		once := true
		for i := 1; i < len(replynum); i++ {
			if once == true {
				content, _ := goquery.NewDocument(gjson.Parse(ReplyNumList).Get("1").Get("href").String())
				data, _ := content.Find("#j_p_postlist > div.l_post.j_l_post.l_post_bright.noborder").Attr("data-field")
				tm, _ := time.Parse("2006-01-02 15:04", gjson.Parse(data).Get("content.date").String())
				replytime := gjson.Parse(ReplyNumList).Get(strconv.Itoa(i)).Get("time").Int()
				if replytime - tm.Unix() < int64(3600) {
					json := gjson.Parse(ReplyNumList).Get("0")
					diff := json.Get("replynum").Int()
					user := akb_diff{Author:json.Get("author").String(), Diff:diff, Href:json.Get("href").String(), Time:json.Get("time").Int(), Title:json.Get("title").String()}
					users = append(users, user)
				}
				once = false
			}
			if replynum[i] != replynum[i - 1] {
				diff := gjson.Parse(ReplyNumList).Get(strconv.Itoa(i)).Get("replynum").Int() - gjson.Parse(ReplyNumList).Get(strconv.Itoa(i - 1)).Get("replynum").Int()
				if diff > 0 {
					json := gjson.Parse(ReplyNumList).Get(strconv.Itoa(i - 1))
					user := akb_diff{Author:json.Get("author").String(), Diff:diff, Href:json.Get("href").String(), Time:json.Get("time").Int(), Title:json.Get("title").String()}
					users = append(users, user)
				}
			}
		}
		engine.Sync2(new(akb_diff))
		affected, err := engine.Insert(&users)
		users = make([]akb_diff, 0)
		fmt.Println(affected, err)
		return true
	})
}
