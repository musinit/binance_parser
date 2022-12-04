package config

type CloudfareCnf struct {
	API string
}

type Config struct {
	Port            string
	CloudfareConfig CloudfareCnf
	CloudflareAPI   string
	RedisConfig     RedisCnf
}

type RedisCnf struct {
	Addr   string
	Prefix string
}
