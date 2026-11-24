package service

import (
	"testing"

	"server/internal/model"
)

func TestFormatProvidePlayListsFiltersCloudDrive(t *testing.T) {
	sources := []model.PlayLinkVo{
		{
			SourceId: "s1",
			Name:     "线路A",
			LinkList: []model.MovieUrlInfo{{Episode: "第1集", Link: "https://a.com/1.m3u8"}},
		},
		{
			SourceId: "s2",
			Name:     "站-夸克网盘",
			LinkList: []model.MovieUrlInfo{{Episode: "正片", Link: "https://pan.quark.cn/s/abc123"}},
		},
		{
			SourceId: "s3",
			Name:     "站-线路B",
			LinkList: []model.MovieUrlInfo{{Episode: "正片", Link: "https://www.aliyundrive.com/s/xyz"}},
		},
		{
			SourceId: "manual",
			IsManual: true,
			Name:     "自建线路",
			LinkList: []model.MovieUrlInfo{{Episode: "正片", Link: "https://m.com/play.m3u8"}},
		},
	}

	gotFrom, gotURL := formatProvidePlayLists(sources)
	wantFrom := "自建线路$$$线路A"
	if gotFrom != wantFrom {
		t.Fatalf("vod_play_from = %q, want %q", gotFrom, wantFrom)
	}
	wantURL := "正片$https://m.com/play.m3u8$$$第1集$https://a.com/1.m3u8"
	if gotURL != wantURL {
		t.Fatalf("vod_play_url = %q, want %q", gotURL, wantURL)
	}
}

func TestFormatProvidePlayListsFiltersCloudLinksInsideMixedSource(t *testing.T) {
	sources := []model.PlayLinkVo{
		{
			SourceId: "s1",
			Name:     "站-线路A",
			LinkList: []model.MovieUrlInfo{
				{Episode: "第1集", Link: "https://a.com/1.m3u8"},
				{Episode: "网盘", Link: "https://pan.baidu.com/s/1abc"},
				{Episode: "第2集", Link: "http://a.com/2.mp4"},
			},
		},
	}

	gotFrom, gotURL := formatProvidePlayLists(sources)
	if gotFrom != "站-线路A" {
		t.Fatalf("vod_play_from = %q, want %q", gotFrom, "站-线路A")
	}
	wantURL := "第1集$https://a.com/1.m3u8#第2集$http://a.com/2.mp4"
	if gotURL != wantURL {
		t.Fatalf("vod_play_url = %q, want %q", gotURL, wantURL)
	}
}

func TestFilterCloudDriveFromSummary(t *testing.T) {
	summary := "站-线路A$$$站-夸克网盘$$$站-115网盘$$$站-115$$$站-线路B"
	got := filterCloudDriveFromSummary(summary)
	want := "站-线路A$$$站-线路B"
	if got != want {
		t.Fatalf("filterCloudDriveFromSummary = %q, want %q", got, want)
	}
}

func TestIsCloudDriveURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://pan.quark.cn/s/abcdef", true},
		{"https://pan.baidu.com/s/1abc", true},
		{"https://www.aliyundrive.com/s/xyz", true},
		{"https://cloud.189.cn/t/abc", true},
		{"https://caiyun.139.com/m/i?abc", true},
		{"https://drive.uc.cn/s/abc", true},
		{"https://www.123pan.com/s/abc", true},
		{"https://pan.example.com/api/fs/list", true},
		{"https://a.com/1.m3u8", false},
		{"https://www.baidu.com/s?wd=movie", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := isCloudDriveURL(tt.url); got != tt.want {
			t.Errorf("isCloudDriveURL(%q) = %v, want %v", tt.url, got, tt.want)
		}
	}
}
