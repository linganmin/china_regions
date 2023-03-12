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

func recursionGetRegion(ctx context.Context, parentPages []crawler.RegionPage, targetRegionPages *[]crawler.RegionPage) {

	for _, page := range parentPages {
		page := page

		if page.Level == crawler.LevelProvince {
			*targetRegionPages = append(*targetRegionPages, page)
		}
		regions := crawler.FetchPages(ctx, page.Url, page.Level+1, page.Code)

		*targetRegionPages = append(*targetRegionPages, regions...)

		if page.Level+1 < crawler.LevelVillage {
			recursionGetRegion(ctx, regions, targetRegionPages)
		}
	}

}

func main() {
	ctx := context.WithValue(context.Background(), "request_id", uuid.NewString())
	logger := zaplog.FromContext(ctx)
	rand.Seed(time.Now().UnixNano())

	// 所有省份
	provinces := crawler.FetchProvincePages(ctx, crawler.BaseUrl)
	allRegions := make([]crawler.RegionPage, 0)

	recursionGetRegion(ctx, provinces, &allRegions)

	wd, err := os.Getwd()
	if err != nil {
		logger.Errorf("get wd error:%+v", err)
		return
	}
	filename := wd + "/" + time.Now().Format("2006-01-02") + ".sql"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	writer := bufio.NewWriter(file)

	for _, region := range allRegions {
		region := region
		_, err = writer.WriteString(fmt.Sprintf("insert into regions (code,name, level, p_code) values (%v,'%v',%v,%v);\n", region.Code, region.Name, region.Level, region.PCode))
	}

	err = writer.Flush()

	if err != nil {
		logger.Errorf("writer.Flush error:%+v", err)
	}
}
