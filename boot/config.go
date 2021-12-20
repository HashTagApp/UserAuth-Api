package boot

var (
	AppConf AppConfig
)

const configPath = "config.app.yamal"

type AppConfig struct {
	Port int `yaml:"port" json:"port"`
}

func parseAppConfig() {

	//cfg.DefaultConfigurator(``, &AppConf, func(config interface{}) {

	//	})

}
