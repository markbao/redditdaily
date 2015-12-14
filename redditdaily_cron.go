package main

import (
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"github.com/markbao/redditdaily/redditdaily"
	"fmt"
	"time"
)

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	c := cron.New();
	error := c.AddFunc(fmt.Sprintf("0 %v %v * * *", viper.GetInt("cron_min"), viper.GetInt("cron_hour")), redditdaily.Run)
	if error != nil {
		fmt.Println(error)
	}
	c.Start();

	fmt.Println("Cron started.")

	select{}
}
