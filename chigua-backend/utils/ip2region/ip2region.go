package ip2region

import (
	"chigua-backend/utils/logger"

	"github.com/lionsoul2014/ip2region/binding/golang/service"
)

var ip2Region *service.Ip2Region

func InitIp2Region() {
	// 配置创建 Ip2Region 查询服务
	v4Config, err := service.NewV4Config(service.VIndexCache, "resources/ip2region_v4.xdb", 20)

	if err != nil {
		logger.Fatalf("failed to create v4 config: %s", err)
	}

	v6Config, err := service.NewV6Config(service.VIndexCache, "resources/ip2region_v6.xdb", 20)
	if err != nil {
		logger.Fatalf("failed to create v6 config: %s", err)
	}

	ip2Region, err = service.NewIp2Region(v4Config, v6Config)
	if err != nil {
		logger.Fatalf("failed to create ip2region service: %s", err)
	}
}

func SearchArea(ip string) string {
	area, err := ip2Region.Search(ip)
	if err != nil {
		logger.Errorf("failed to search ip2region: %s", err)
		return ""
	}
	return area
}

func CloseIp2Region() {
	ip2Region.Close()
}
