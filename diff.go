package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/configor"
	"github.com/xormplus/xorm"
	"github.com/xormplus/core"
	"fmt"
	"github.com/tidwall/gjson"
)

var Config = struct {
	DB struct {
		   Host     string
		   User     string `default:"root"`
		   Password string
		   Port     string `default:"3306"`
		   Database string
	   }
}{}


type cat_diff struct {
	Author string
	Diff int64
	Href string
	Title string
	Time int64
}

var engine *xorm.Engine
var users []cat_diff

func main() {
	time := 1497686404;
	for {
		if time > 1500347801 {
			break
		}
		sql := "SELECT title,replynum,author,time,href FROM cat_time WHERE time BETWEEN ? AND ? "
		configor.Load(&Config, "./config.yml")
		engine, _ = xorm.NewEngine("mysql", Config.DB.User + ":" + Config.DB.Password + "@tcp(" + Config.DB.Host + ":" + Config.DB.Port + ")/" + Config.DB.Database)
		engine.Logger().SetLevel(core.LOG_DEBUG)
		results, _ := engine.Sql(sql, time-3700, time-100).Query().Json()
		result, _ := engine.Sql(sql, time -100, time + 3500).Query().Json()
		resjson := gjson.Parse(result)
		resjsons := gjson.Parse(results)
		resjson.ForEach(func(key, value gjson.Result) bool {
			title := value.Get("title").String()
			newnum := value.Get("replynum").Int()
			oldnum := resjsons.Get(`#[title=="` + title + `"].replynum`).Int()
			if oldnum != 0 {
				diff := newnum - oldnum
				if diff != 0 {
					user := cat_diff{Author:value.Get("author").String(), Diff:diff, Title:title, Time:value.Get("time").Int(), Href:value.Get("href").String()}
					users = append(users, user)
				}
			} else {
				user := cat_diff{Author:value.Get("author").String(), Diff:newnum, Title:title, Time:value.Get("time").Int(), Href:value.Get("href").String()}
				users = append(users, user)
			}
			return true
		})
		engine.Sync2(new(cat_diff))
		affected, err := engine.Insert(&users)
		users=make([]cat_diff,0)
		fmt.Println(affected, err)
		time = time + 3600
	}
}
