package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config 系统配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Template TemplateConfig `yaml:"template"`
	Upload   UploadConfig   `yaml:"upload"`
	Log      LogConfig      `yaml:"log"`
	Cache    CacheConfig    `yaml:"cache"`
	Site     SiteConfig     `yaml:"site"`
	API      APIConfig      `yaml:"api"`
	Plugin   PluginConfig   `yaml:"plugin"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host              string `yaml:"host"`
	Port              int    `yaml:"port"`
	ReadTimeout       int    `yaml:"readTimeout"`
	WriteTimeout      int    `yaml:"writeTimeout"`
	MaxHeaderBytes    int    `yaml:"maxHeaderBytes"`
	JWTSecret         string `yaml:"jwtSecret"`
	SessionSecret     string `yaml:"sessionSecret"`
	Secure            bool   `yaml:"secure"`
	EnableRateLimit   bool   `yaml:"enableRateLimit"`
	RateLimitRequests int    `yaml:"rateLimitRequests"`
	RateLimitWindow   int    `yaml:"rateLimitWindow"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Prefix   string `yaml:"prefix"`
	Charset  string `yaml:"charset"`
	MaxIdle  int    `yaml:"maxIdle"`
	MaxOpen  int    `yaml:"maxOpen"`
}

// TemplateConfig 模板配置
type TemplateConfig struct {
	Dir        string `yaml:"dir"`
	Cache      bool   `yaml:"cache"`
	DefaultTpl string `yaml:"defaultTpl"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	Dir            string  `yaml:"dir"`
	MaxSize        int     `yaml:"maxSize"`
	AllowedExts    string  `yaml:"allowedExts"`
	DenyExts       string  `yaml:"denyExts"`
	Watermark      bool    `yaml:"watermark"`
	WatermarkImg   string  `yaml:"watermarkImg"`
	WatermarkPos   string  `yaml:"watermarkPos"`
	WatermarkText  string  `yaml:"watermarkText"`
	WatermarkColor string  `yaml:"watermarkColor"`
	WatermarkAlpha float64 `yaml:"watermarkAlpha"`
	ImageMaxWidth  int     `yaml:"imageMaxWidth"`
	ImageMaxHeight int     `yaml:"imageMaxHeight"`
}

// LogConfig 日志配置
type LogConfig struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Type               string `yaml:"type"`
	Path               string `yaml:"path"`
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	Password           string `yaml:"password"`
	Expire             int    `yaml:"expire"`
	EnableArticleCache bool   `yaml:"enableArticleCache"`
	ArticleCacheTime   int    `yaml:"articleCacheTime"`
	EnableListCache    bool   `yaml:"enableListCache"`
	ListCacheTime      int    `yaml:"listCacheTime"`
	EnableTagCache     bool   `yaml:"enableTagCache"`
	TagCacheTime       int    `yaml:"tagCacheTime"`
	EnableSpecialCache bool   `yaml:"enableSpecialCache"`
	SpecialCacheTime   int    `yaml:"specialCacheTime"`
	EnableSearchCache  bool   `yaml:"enableSearchCache"`
	SearchCacheTime    int    `yaml:"searchCacheTime"`
	EnableStatsCache   bool   `yaml:"enableStatsCache"`
	StatsCacheTime     int    `yaml:"statsCacheTime"`
	EnableSitemapCache bool   `yaml:"enableSitemapCache"`
	SitemapCacheTime   int    `yaml:"sitemapCacheTime"`
	EnableRSSCache     bool   `yaml:"enableRSSCache"`
	RSSCacheTime       int    `yaml:"rssCacheTime"`
}

// SiteConfig 站点配置
type SiteConfig struct {
	Name             string `yaml:"name"`
	URL              string `yaml:"url"`
	Keywords         string `yaml:"keywords"`
	Description      string `yaml:"description"`
	Email            string `yaml:"email"`
	ICP              string `yaml:"icp"`
	StatCode         string `yaml:"statCode"`
	CopyRight        string `yaml:"copyRight"`
	DefaultLang      string `yaml:"defaultLang"`
	StaticIndex      bool   `yaml:"staticIndex"`
	StaticList       bool   `yaml:"staticList"`
	StaticArticle    bool   `yaml:"staticArticle"`
	StaticSpecial    bool   `yaml:"staticSpecial"`
	StaticTag        bool   `yaml:"staticTag"`
	StaticMobile     bool   `yaml:"staticMobile"`
	StaticDir        string `yaml:"staticDir"`
	StaticSuffix     string `yaml:"staticSuffix"`
	StaticMobileDir  string `yaml:"staticMobileDir"`
	StaticMobileSuffix string `yaml:"staticMobileSuffix"`
	Close            bool   `yaml:"close"`
	CloseReason      string `yaml:"closeReason"`
	CommentAutoCheck bool   `yaml:"commentAutoCheck"`
	SessionSecret    string `yaml:"sessionSecret"`
}

// APIConfig API配置
type APIConfig struct {
	Enabled bool   `yaml:"enabled"`
	Key     string `yaml:"key"`
	Secret  string `yaml:"secret"`
	Cors    bool   `yaml:"cors"`
}

// PluginConfig 插件配置
type PluginConfig struct {
	Dir        string `yaml:"dir"`
	ConfigFile string `yaml:"configFile"`
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Save 保存配置到文件
func (c *Config) Save() error {
	// 将配置转换为YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// 写入文件
	err = ioutil.WriteFile("config.yaml", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
