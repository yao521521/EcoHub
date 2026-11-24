package service

import (
	"net/url"
	"strings"
	"unicode"

	"server/internal/model"
)

// cloudDriveHostSuffixes 常见网盘分享域名。命中任一个即认为该播放地址是网盘地址。
var cloudDriveHostSuffixes = []string{
	"pan.quark.cn", "quark.cn",
	"aliyundrive.com", "alipan.com",
	"pan.baidu.com", "yun.baidu.com",
	"115.com", "115cdn.com",
	"cloud.189.cn",
	"caiyun.139.com", "yun.139.com",
	"drive.uc.cn",
	"pan.xunlei.com",
	"123pan.com", "123684.com", "123912.com",
	"lanzou.com", "lanzoui.com", "lanzoux.com",
	"weiyun.com",
	"pikpak.com", "mypikpak.com",
	"pan38.com",
}

// cloudDriveNameKeywords 线路名里的网盘关键词，用于兜住链接尚未入库时的展示名。
var cloudDriveNameKeywords = []string{
	"夸克", "quark",
	"阿里云盘", "aliyundrive", "alipan",
	"百度网盘", "百度云盘", "baidunetdisk", "pan.baidu",
	"115网盘", "天翼云盘", "189云盘", "移动云盘", "139云盘",
	"uc网盘", "123云盘", "123pan",
	"迅雷网盘", "lanzou", "蓝奏云",
	"微云", "weiyun", "pikpak",
	"alist", "openlist",
}

// cloudDriveNameTokens 采集源 flag 常见的短标识，要求完整匹配，避免误伤片名/站名。
var cloudDriveNameTokens = map[string]struct{}{
	"quark": {}, "baidu": {}, "115": {}, "uc": {}, "xunlei": {},
	"lanzou": {}, "pikpak": {}, "weiyun": {},
}

func isCloudDriveURL(raw string) bool {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false
	}
	lower := strings.ToLower(raw)
	if !strings.Contains(lower, "://") {
		if !strings.HasPrefix(lower, "//") {
			return false
		}
		lower = "https:" + lower
	}
	u, err := url.Parse(lower)
	if err != nil {
		return false
	}
	host := strings.TrimSuffix(strings.ToLower(u.Hostname()), ".")
	if host == "" {
		return false
	}
	for _, suffix := range cloudDriveHostSuffixes {
		if host == suffix || strings.HasSuffix(host, "."+suffix) {
			return true
		}
	}
	// AList / OpenList 自建网盘的接口路径
	if strings.Contains(strings.ToLower(u.Path), "/api/fs/") {
		return true
	}
	return false
}

func isCloudDriveName(name string) bool {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return false
	}
	for _, keyword := range cloudDriveNameKeywords {
		if strings.Contains(name, keyword) {
			return true
		}
	}
	for _, token := range strings.FieldsFunc(name, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if _, ok := cloudDriveNameTokens[token]; ok {
			return true
		}
	}
	return false
}

// filterCloudDriveSources 过滤网盘线路：线路名命中网盘关键词时整条移除；
// 否则逐条剔除网盘链接，混排线路只保留剩余直链。
func filterCloudDriveSources(sources []model.PlayLinkVo) []model.PlayLinkVo {
	if len(sources) == 0 {
		return nil
	}
	filtered := make([]model.PlayLinkVo, 0, len(sources))
	for _, source := range sources {
		if isCloudDriveName(source.Name) {
			continue
		}
		links := make([]model.MovieUrlInfo, 0, len(source.LinkList))
		for _, link := range source.LinkList {
			if isCloudDriveURL(link.Link) {
				continue
			}
			links = append(links, link)
		}
		if len(links) == 0 {
			continue
		}
		source.LinkList = links
		filtered = append(filtered, source)
	}
	return filtered
}

// filterCloudDriveFromSummary 过滤列表页的播放线路摘要。摘要只有展示名，按名称关键词剔除。
func filterCloudDriveFromSummary(summary string) string {
	summary = strings.TrimSpace(summary)
	if summary == "" {
		return ""
	}
	parts := strings.Split(summary, "$$$")
	kept := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || isCloudDriveName(part) {
			continue
		}
		kept = append(kept, part)
	}
	return strings.Join(kept, "$$$")
}
