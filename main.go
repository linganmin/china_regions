package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/linganmin/zaplog"

	"github.com/linganmin/china_regions/crawler"
)

func main() {
	ctx := context.WithValue(context.Background(), "request_id", uuid.NewString())
	logger := zaplog.FromContext(ctx)
	rand.Seed(time.Now().UnixNano())

	wd, err := os.Getwd()
	if err != nil {
		logger.Errorf("get wd error:%+v", err)
		return
	}
	filename := wd + "/" + time.Now().Format("2006-01-02") + ".sql"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	writer := bufio.NewWriter(file)

	// 所有省份
	provinces := crawler.FetchProvincePages(ctx, crawler.BaseUrl)
	for _, province := range provinces {
		_, err = writer.WriteString(fmt.Sprintf("------------------ %s\n", province.Name))
		_, err = writer.WriteString(fmt.Sprintf("insert into regions (code,name, level, p_code) values (%v,'%v',%v,%v);\n", province.Code, province.Name, province.Level, 0))

		if err != nil {
			logger.Errorf("write file error:%+v", err)
			return
		}
		// 地级市
		cities := crawler.FetchPages(ctx, province.Url, crawler.LevelCity)
		for _, city := range cities {
			_, err = writer.WriteString(fmt.Sprintf("--- %s\n", city.Name))

			_, err = writer.WriteString(fmt.Sprintf(fmt.Sprintf("insert into regions (code,name, level, p_code) values (%v,'%v',%v,%v);\n", city.Code, city.Name, city.Level, province.Code)))
			if err != nil {
				logger.Errorf("write file error:%+v", err)
				return
			}

			// 县/区
			counties := crawler.FetchPages(ctx, city.Url, crawler.LevelCounty)
			for _, county := range counties {
				_, err = writer.WriteString(fmt.Sprintf("--- %s\n", county.Name))

				_, err = writer.WriteString(fmt.Sprintf(fmt.Sprintf("insert into regions (code,name, level, p_code) values (%v,'%v',%v,%v);\n", county.Code, county.Name, county.Level, city.Code)))
				if err != nil {
					logger.Errorf("write file error:%+v", err)
					return
				}

				// 镇/街道
				towns := crawler.FetchPages(ctx, county.Url, crawler.LevelTown)
				for _, town := range towns {
					_, err = writer.WriteString(fmt.Sprintf("--- %s\n", town.Name))
					_, err = writer.WriteString(fmt.Sprintf(fmt.Sprintf("insert into regions (code,name, level, p_code) values (%v,'%v',%v,%v);\n", town.Code, town.Name, town.Level, county.Code)))
					if err != nil {
						logger.Errorf("write file error:%+v", err)
						return
					}

					// 村/居委会
					villages := crawler.FetchPages(ctx, town.Url, crawler.LevelVillage)
					for _, village := range villages {
						_, err = writer.WriteString(fmt.Sprintf(fmt.Sprintf("insert into regions (code,name, level, p_code) values (%v,'%v',%v,%v);\n", village.Code, village.Name, village.Level, town.Code)))
						if err != nil {
							logger.Errorf("write file error:%+v", err)
							return
						}
					}
				}

			}

		}

	}

	err = writer.Flush()

	if err != nil {
		logger.Errorf("writer.Flush error:%+v", err)
	}
}
