package conf

import (
	"github.com/spf13/viper"
	"user-service/internal/util"
)

func Init() {

	listDir := []string{".", "../", "../../", "../../../", "../../../../"}

	for _, dir := range listDir {
		viper.SetConfigName("env")
		viper.SetConfigType("json")
		viper.AddConfigPath(dir)
		err := viper.ReadInConfig()
		if err == nil {
			viper.SetConfigName("env.override")
			err = viper.MergeInConfig()
			util.Panic(err)
			return
		}
	}

	panic("cannot load env")
}

func GetDatabaseDSN() string {
	return viper.GetString("DATABASE_DSN")
}
