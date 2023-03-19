package crawler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/linganmin/zaplog"
)

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"

const BaseUrl = "http://www.stats.gov.cn/sj/tjbz/tjyqhdmhcxhfdm/2022/"

func fetch(ctx context.Context, url string) (string, error) {
	logger := zaplog.FromContext(ctx)
	logger.Debugf("fetch url:%+v", url)

	// 页面有并发限制，随机暂停
	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(100)+60))

	cli := &http.Client{}

	cli.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", ua)

	resp, err := cli.Do(req)
	if err != nil {
		logger.Errorf("cli.Do resp:%v error:%+v", resp, err)

		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Http resp status not ok:%+v", resp))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Errorf("io.ReadAll body:%v error:%+v", string(body), err)
		return "", err
	}

	return string(body), nil
}

type RegionPage struct {
	PCode string `json:"p_code"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Url   string `json:"url"`
	Level Level  `json:"level"`
}

type Level int

const (
	LevelProvince Level = iota + 1
	LevelCity
	LevelCounty
	LevelTown
	LevelVillage
)

func FetchProvincePages(ctx context.Context, url string) []RegionPage {
	logger := zaplog.FromContext(ctx)
	content, err := fetch(ctx, url)

	if err != nil {
		logger.Errorf("FetchProvincePages error:%+v", err)
		return nil
	}
	rg := regexp.MustCompile("href=\"(\\d*\\.html)\">([\u4e00-\u9fa5]*)")

	items := rg.FindAllStringSubmatch(content, -1)

	var list []RegionPage

	for _, item := range items {

		if err != nil {
			logger.Errorf("strconv.Atoi error:%+v", err)
			continue
		}

		list = append(list, RegionPage{
			Code:  strings.Replace(item[1], ".html", "", -1) + strings.Repeat("0", 10),
			Url:   BaseUrl + item[1],
			Name:  item[2],
			Level: LevelProvince,
			PCode: "0",
		})
	}
	return list
}

func FetchPages(ctx context.Context, url string, level Level, pCode string) []RegionPage {
	if url == BaseUrl {
		return nil
	}
	logger := zaplog.FromContext(ctx)

	content, err := fetch(ctx, url)

	if err != nil { // 进重试
		for i := 3; i > 0; i-- {
			time.Sleep(time.Second * 2)
			content, err = fetch(ctx, url)
			if err != nil {
				continue
			} else {
				break
			}
		}
	}

	if err != nil {
		logger.Errorf("FetchCityPages error:%+v", err)
		return nil
	}

	content = strings.Replace(content, "\r\n", "", -1)

	//rg := regexp.MustCompile("(\\d+/\\d+\\.html)*[\">td<]*(\\d{12})[<A-z/> =\\d.\"]*([\u2e80-\ufffdh]{0,30})")
	rg := regexp.MustCompile("(\\d+/\\d+\\.html)*[\">td<]*(\\d{12})[<A-z/> =\\d.\"]*([⺀-龥0-9（）()〇\uE170𡐓𡌶]{0,30})")
	items := rg.FindAllStringSubmatch(content, -1)

	var list []RegionPage

	for _, item := range items {

		url := BaseUrl

		if item[1] != "" {
			switch level {
			case LevelCity:
				url = BaseUrl + item[1]
			case LevelCounty:
				url = BaseUrl + item[2][:2] + "/" + item[1]
			case LevelTown:
				url = BaseUrl + item[2][:2] + "/" + item[2][2:4] + "/" + item[1]

			}
		}

		list = append(list, RegionPage{
			PCode: pCode,
			Code:  item[2],
			Url:   url,
			Name:  item[3],
			Level: getLevelByCode(item[2]),
		})
	}
	return list
}

func getLevelByCode(code string) Level {

	if len(code) != 12 {
		return Level(0)
	}

	//fmt.Println(code[len(code)-3:], code[len(code)-6:len(code)-3], code[len(code)-8:len(code)-6], code[len(code)-10:len(code)-8], code[len(code)-12:len(code)-10])
	if code[len(code)-3:] != "000" {
		return LevelVillage
	}

	if code[len(code)-6:len(code)-3] != "000" {
		return LevelTown
	}
	if code[len(code)-8:len(code)-6] != "00" {
		return LevelCounty
	}
	if code[len(code)-10:len(code)-8] != "00" {
		return LevelCity
	}

	return LevelProvince
}
