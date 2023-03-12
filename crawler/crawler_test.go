package crawler

import (
	"context"
	"reflect"
	"testing"
)

func Test_fetch(t *testing.T) {
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: testing.CoverMode(), args: args{
			ctx: context.Background(),
			url: "https://qqlykm.cn/api/free/weather/get?city=%E8%8B%8F%E5%B7%9E",
		}, want: "", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetch(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fetch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchProvincePages(t *testing.T) {
	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: testing.CoverMode(), args: args{
			ctx: context.Background(),
			url: "http://www.stats.gov.cn/tjsj/tjbz/tjyqhdmhcxhfdm/2022/index.html",
		}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FetchProvincePages(tt.args.ctx, tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchProvincePages() = %v, want %v", got, tt.want)
			}
		})
	}
}
