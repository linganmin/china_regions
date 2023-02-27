package crawler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/linganmin/zaplog"
)

func fetch(ctx context.Context, url string) (string, error) {
	logger := zaplog.FromContext(ctx)
	logger.Debugf("fetch url:%+v", url)

	cli := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	resp, err := cli.Do(req)

	if err != nil {
		logger.Errorf("cli.Do error:%+v", err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Http resp status not ok:%+v", resp))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Errorf("io.ReadAll error:%+v", err)
		return "", err
	}

	return string(body), nil
}
