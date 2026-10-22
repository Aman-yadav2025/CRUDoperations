package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)
//HttpServer config info main address port etc
type HttpServer struct {
	Addr string `yaml:"address" env-required: "true"`
}

// main config info having env path and storage variables
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HttpServer  `yaml:"http_server"`
}


// function that returns the complete config info struct 
//throw fatal stop immediatly if function with naming Must fails
func MustLoad() *Config{

	//config path
	var configPath string

	//check in env
	configPath = os.Getenv("CONFIG_PATH")

	//if env doesnot have configPath then check for flags 
	if configPath == "" {
		//flags.String checks for flags its first parameter is 
		//name of the flag here it is config as in --config=
		//second parameter is default value here it is ""
		// third is the helper text
		flags := flag.String("config", "", "Path to the config file")
		//parse it into a string pointer
		flag.Parse()
		//assign the value to config path if nothing is given it will fall back to default
		configPath = *flags
	}
	//still nothing is found throw err
	if configPath == ""{
		log.Fatal("Config path is not fixed")
	}
	//if the path is given but the file doesn't exist
	if _,err := os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("Config File does not exist: %s: ", configPath)
	}
	//create a struct to store config info
	var cfg Config
	//Read the config info with the help of cleanenv.ReadConfig
	err :=cleanenv.ReadConfig(configPath, &cfg)
	//can't read throw error
	if err != nil{
		log.Fatalf("Can't read config files %s : ", err.Error())
	}
	//return the read info
	return &cfg
}
