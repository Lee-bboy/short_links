package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"shortlinks/common"
	"shortlinks/config"
	"shortlinks/controllers"
	"shortlinks/router"
	"strconv"
	"strings"
	"time"
)

var (
	restore_path = flag.String("p", "", "用来恢复的文件路径")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	var command string
	if len(args) == 0 {
		command = "start"
	} else {
		command = args[0]
	}

	switch command {
	case "start":
		//监测服务区运行状态，如果shortlinks:id出现减小的情况，那么说明数据库数据被破坏，发送邮件提醒
		go func() {
			var id uint64

			//读取shortlinks:id的值
			conn := common.RedisPool.Get()
			reply, _ := conn.Do("GET", "shortlinks:id")
			idBytes, _ := reply.([]byte)
			id, _ = strconv.ParseUint(string(idBytes), 10, 64)

			ticker := time.NewTicker(1 * time.Second)

			for range ticker.C {
				reply, _ := conn.Do("GET", "shortlinks:id")
				idBytes, _ := reply.([]byte)
				newId, _ := strconv.ParseUint(string(idBytes), 10, 64)

				if newId < id {
					subject := "短网址生成服务异常"
					body := "短网址生成服务异常，请及时处理！"
					mail1 := common.MailTo("wangjian@linghit.com")
					mail2 := common.MailTo("zhengjiajia@linghit.com")
					mail3 := common.MailTo("lishiyuan@linghit.com")
					mail1.SubjectFromString(subject).BodyFromString(body)
					mail2.SubjectFromString(subject).BodyFromString(body)
					mail3.SubjectFromString(subject).BodyFromString(body)
					err1 := common.Mailer.Send(mail1)
					err2 := common.Mailer.Send(mail2)
					err3 := common.Mailer.Send(mail3)
					if err1 == nil {
						common.ERROR.Printf("发送邮件失败：%v\r\n", err1)
					}
					if err2 == nil {
						common.ERROR.Printf("发送邮件失败：%v\r\n", err2)
					}
					if err3 == nil {
						common.ERROR.Printf("发送邮件失败：%v\r\n", err3)
					}
				}

				id = newId
			}
		}()

		port := config.GetConf("port", "8080")

		router := router.NewRouter()

		log.Fatal(http.ListenAndServe(":"+port, router))
	case "restore":
		//读取待恢复的文件
		file, err := os.Open(*restore_path)
		if err != nil {
			panic(err)
		}

		csvReader := csv.NewReader(file)
		var records []string
		var serviceUrl *url.URL
		serviceUrl, err = url.Parse(config.GetConf("host", ""))
		if err != nil {
			panic(err)
		}
		//临时短链有效期
		expire, err := strconv.Atoi(config.GetConf("temp_shortlink_expire", "7"))
		if err != nil {
			panic(err)
		}

		for {
			records, err = csvReader.Read()
			if len(records) >= 2 {
				long := records[0]
				short := records[1]

				shortUrl, err := url.Parse(short)
				if err != nil {
					continue
				}

				if shortUrl.Scheme != serviceUrl.Scheme || shortUrl.Host != serviceUrl.Host {
					continue
				}

				//解码
				rawShortlink := shortUrl.Path
				segments := strings.Split(rawShortlink, "/")
				if len(segments) < 2 {
					continue
				}
				shortlink := segments[1]
				if strings.HasSuffix(shortlink, "z") {
					//持久短链接
					continue
				} else {
					id, err := controllers.DecodeShortlink(shortlink)
					if err != nil {
						continue
					}

					conn := common.RedisPool.Get()
					defer conn.Close()

					key := "shortlinks:" + strconv.FormatUint(id, 10)
					reply, err := conn.Do("SETNX", key, long)
					rows, _ := reply.(int64)
					if rows != 0 {
						conn.Do("EXPIRE", key, expire*86400)
					}
				}
			}

			if err != nil {
				if err == io.EOF {
					break
				} else {
					panic(err)
				}
			}
		}
	default:
		flag.Usage()
	}
}

func usage() {
	fmt.Printf(`commands: go run main.go [options] [start|store]
options:
p: the file path for restoring
	`)
}
