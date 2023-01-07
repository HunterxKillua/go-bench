/*
Copyright © 2023 Killua<captainchengjie@gmail.com>

*/
package cmd

import (
	"encoding/json"
	"ginBlog/api"
	"ginBlog/pkg/logger"
	"ginBlog/util/email"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	models "ginBlog/model"
	"ginBlog/server"
	"ginBlog/sql"
)

var cfgFile string
var port string
var emailP bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ginBlog",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init(handleConfigLog())
		defer logger.Sync()
		run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/main.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&port, "port", "p", "3100", "server port")
	rootCmd.Flags().BoolVarP(&emailP, "emailEnable", "e", true, "enable email server")
}

const (
	// recommendedHomeDir 定义放置 ginBlog 服务配置的默认目录.
	recommendedHomeDir = "configs"

	// defaultConfigName 指定了 ginBlog 服务的默认配置文件名.
	defaultConfigName = "main.yaml"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		dir := filepath.Join(recommendedHomeDir, defaultConfigName)
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)
		// Search config in home directory with name ".ginBlog" (without extension).
		// viper.AddConfigPath(filepath.Join(home, recommendedHomeDir, defaultConfigName))
		// viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigName)
		viper.SetConfigFile(dir)
	}

	viper.AutomaticEnv() // read in environment variables that match
	// 读取环境变量的前缀为 GINBLOG，如果是 ginBlog，将自动转变为大写。
	viper.SetEnvPrefix("ginBlog")
	// 将 viper.Get(key) key 字符串中 '.' 和 '-' 替换为 '_'
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Errorw("Using config file", "file", viper.ConfigFileUsed())
	}

	// 打印 viper 当前使用的配置文件，方便 Debug.
	logger.Infow("Using config file", "file", viper.ConfigFileUsed())
}

func run() error {
	// 打印所有的配置项及其值
	settings, _ := json.Marshal(viper.AllSettings())
	logger.Infow(string(settings))
	// 打印 db -> username 配置项的值
	logger.Infow(viper.GetString("db.username"))
	api.DB = sql.ConnectDB(handleConfigDB())
	email.LoginAuth(handleConfigEmail())
	server.Run(":" + handleConfigServer())
	return nil
}

// 处理email配置参数
func handleConfigEmail() (string, string) {
	return viper.GetString("email.username"), viper.GetString("email.secrets")
}

// 处理orm配置参数
func handleConfigDB() *models.ConfDB {
	return &models.ConfDB{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Database: viper.GetString("db.database"),
		MaxIdle:  viper.GetInt("db.max-idle-connections"),
		MaxOpen:  viper.GetInt("db.max-open-connections"),
		MaxLive:  viper.GetInt("db.max-connection-life-time"),
	}
}

// 处理服务配置参数
func handleConfigServer() string {
	return viper.GetString("server.port")
}

// 处理日志参数
func handleConfigLog() *logger.Options {
	return &logger.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}
