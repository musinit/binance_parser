package config

type CloudfareETHCnf struct {
	API string
}

type Config struct {
	Port               string
	CloudfareEHTConfig CloudfareETHCnf
	BlockscautAPI      string
	RedisConfig        RedisCnf
}

type RedisCnf struct {
	Addr   string
	Prefix string
}
