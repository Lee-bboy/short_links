package controllers

import (
	"errors"
	"math/rand"
	"net/http"
	"short_links/common"
	"short_links/config"
	"short_links/models"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var ControllerApi = new(ApiController)

var (
	ErrLongLinkNotFound = errors.New("long link not found")
)

type LinkResult struct {
	Links map[string]string `json:"links"`
}

type ApiController struct {
	Controller
}

func (c *ApiController) ConvertLinks(w http.ResponseWriter, r *http.Request) {
	data, err := c.parseParams(r)
	if err != nil {
		w.Write([]byte(`{"errcode": 400, "errmsg": "请求不正确"}`))
		return
	}

	//参数验证
	appId, err := common.GetAppId(r.FormValue("appKey"))
	if err != nil {
		w.Write([]byte(`{"errcode": 400, "errmsg": "请求不正确"}`))
		return
	}

	var links = make([]*models.Link, 0)
	l, ok := data["links"]
	if !ok {
		w.Write([]byte(`{"errcode": 400, "errmsg": "请求不正确"}`))
		return
	}
	tmpLinks, ok := l.([]interface{})
	if !ok {
		w.Write([]byte(`{"errcode": 400, "errmsg": "请求不正确"}`))
		return
	}
	for _, tmp := range tmpLinks {
		tmpLongLink, ok := tmp.(string)
		if !ok {
			w.Write([]byte(`{"errcode": 400, "errmsg": "请求不正确"}`))
			return
		}

		tmpLink := new(models.Link)
		tmpLink.LongLink = tmpLongLink
		tmpLink.AppId = appId
		links = append(links, tmpLink)
	}

	var persist = false
	p, ok := data["persist"]
	if ok {
		persist, ok = p.(bool)
		if !ok {
			w.Write([]byte(`{"errcode": 400, "errmsg": "请求不正确"}`))
			return
		}
	}

	var result = new(LinkResult)
	result.Links = make(map[string]string)
	for _, link := range links {
		var short string
		var err error
		if persist {
			short, err = c.convertPersistLink(link)
		} else {
			short, err = c.convertTempLink(link)
		}
		if err != nil {
			result.Links[link.LongLink] = "unkown"
		} else {
			result.Links[link.LongLink] = short
		}
	}

	c.jsonResponse(w, result)
}

func (c *ApiController) JumpToLongLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortlink, ok := vars["shortlink"]
	if !ok {
		w.Write([]byte("页面不存在或已过期"))
		return
	}

	var longLink string
	var err error
	if strings.HasSuffix(shortlink, "z") {
		longLink, err = c.getPersistLongLink(strings.TrimRight(shortlink, "z"))
	} else {
		longLink, err = c.getTempLongLink(shortlink)
	}
	if err != nil {
		if err == ErrLongLinkNotFound {
			w.Write([]byte("页面不存在或已过期"))
			return
		} else {
			w.Write([]byte("请求失败，请重试"))
			return
		}
	}

	http.Redirect(w, r, longLink, http.StatusFound)
}

func (c *ApiController) convertPersistLink(l *models.Link) (string, error) {
	var short string
	var err error

	engine := common.RegisterMysql()
	_, err = engine.Insert(l)
	var log = &models.Log{AppId: l.AppId, LongLink: l.LongLink, Persist: 1}
	if err != nil {
		log.Status = 1
		log.StatusText = err.Error()

		short = ""
	} else {
		log.Status = 0
		log.StatusText = ""

		//等号结尾的说明是永久性链接
		short = config.GetConf("host", "") + "/" + EncodeShortlink(l.Id) + "z"
	}

	engine.Insert(log)

	return short, err
}

func (c *ApiController) getPersistLongLink(shortlink string) (string, error) {
	linkId, err := DecodeShortlink(shortlink)
	if err != nil {
		return "", ErrLongLinkNotFound
	}

	engine := common.RegisterMysql()
	link := new(models.Link)
	has, err := engine.Id(linkId).Get(link)
	if err != nil {
		return "", err
	}
	if !has {
		return "", ErrLongLinkNotFound
	}

	return link.LongLink, nil
}

func (c *ApiController) convertTempLink(l *models.Link) (string, error) {
	autoIncrementId := common.NewAutoIncrementId("short_links:id")

	//临时短链有效期
	expire, err := strconv.Atoi(config.GetConf("temp_shortlink_expire", "7"))
	if err != nil {
		return "", err
	}

	engine := common.RegisterMysql()
	var log = &models.Log{AppId: l.AppId, LongLink: l.LongLink, Persist: 0}

	var id uint64
	conn := common.RedisPool.Get()
	defer conn.Close()
	for {
		id = autoIncrementId.Get()
		key := "short_links:" + strconv.FormatUint(id, 10)
		reply, err := conn.Do("SETNX", key, l.LongLink)
		if err != nil {
			log.Status = 1
			log.StatusText = err.Error()
			engine.Insert(log)

			return "", err
		}

		rows, _ := reply.(int64)
		if rows != 0 {
			conn.Do("EXPIRE", key, expire*86400)
			break
		}
	}

	log.Status = 0
	log.StatusText = ""
	engine.Insert(log)

	//等号结尾的说明是永久性链接
	short := config.GetConf("host", "") + "/" + EncodeShortlink(id)

	return short, nil
}

func (c *ApiController) getTempLongLink(shortlink string) (string, error) {
	linkId, err := DecodeShortlink(shortlink)
	if err != nil {
		return "", ErrLongLinkNotFound
	}

	conn := common.RedisPool.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", "short_links:"+strconv.FormatUint(linkId, 10))
	if err != nil {
		return "", err
	}

	longLinkBytes, _ := reply.([]byte)
	if len(longLinkBytes) == 0 {
		return "", ErrLongLinkNotFound
	}

	return string(longLinkBytes), nil
}

func EncodeShortlink(num uint64) string {
	plain := []byte(strconv.FormatUint(num, 35))

	randByte := getRandomByte()

	for i, b := range plain {
		plain[i] = encodeByte(b, int(randByte))
	}
	plain = append(plain, randByte)

	return string(plain)
}

func DecodeShortlink(shortlink string) (uint64, error) {
	bytes := []byte(shortlink)
	l := len(bytes)
	cipher := bytes[0 : l-1]
	randByte := bytes[l-1]

	for i, b := range cipher {
		cipher[i] = decodeByte(b, int(randByte))
	}

	return strconv.ParseUint(string(cipher), 35, 64)
}

func encodeByte(b byte, step int) byte {
	if step < 0 {
		step = -1 * step
	}

	if step > 9 {
		step = step % 10
	}

	if b >= 48 && b <= 57 {
		//如果是数字
		b += byte(step)
		if b > 57 {
			b -= 10
		}
	} else if b >= 97 && b <= 121 {
		//如果是小写字母
		b += byte(step)
		if b > 121 {
			b -= 25
		}
	} else if b >= 65 && b <= 89 {
		//如果是大写字母
		b += byte(step)
		if b > 89 {
			b -= 25
		}
	}

	return b
}

func decodeByte(b byte, step int) byte {
	if step < 0 {
		step = -1 * step
	}

	if step > 9 {
		step = step % 10
	}

	if b >= 48 && b <= 57 {
		//如果是数字
		b -= byte(step)
		if b < 48 {
			b += 10
		}
	} else if b >= 97 && b <= 121 {
		//如果是小写字母
		b -= byte(step)
		if b < 97 {
			b += 25
		}
	} else if b >= 65 && b <= 89 {
		//如果是大写字母
		b -= byte(step)
		if b < 65 {
			b += 25
		}
	}

	return b
}

func getRandomByte() byte {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxy")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return bytes[r.Intn(len(bytes))]
}
