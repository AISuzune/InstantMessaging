package main

import (
	"InstantMessaging/boot"
)

func main() {
	boot.ViperSetup("./manifest/config/config.yaml")
	boot.LoggerSetup()
	boot.MysqlDBSetup()

	boot.RedisSetup()
	boot.InitCasbin()
	boot.ServerSetup()
}
