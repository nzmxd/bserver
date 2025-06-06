package initialize

import (
	"bufio"
	"github.com/nzmxd/bserver/global"
	"github.com/nzmxd/bserver/utils"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"os"
	"strings"
)

func OtherInit() {
	dr, err := utils.ParseDuration(global.CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(global.CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err)
	}

	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
	file, err := os.Open("go.mod")
	if err == nil && global.CONFIG.AutoCode.Module == "" {
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		global.CONFIG.AutoCode.Module = strings.TrimPrefix(scanner.Text(), "module ")
	}
}
