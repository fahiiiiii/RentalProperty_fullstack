package conf

import (
    "github.com/beego/beego/v2/core/config"
    "log"
    "strconv"
)

type AppConfiguration struct {
    API struct {
        RapidAPIKey     string
        RequestTimeout  int
    }
    RateLimit struct {
        RequestsPerSecond float64
        BurstSize         int
    }
}

var AppConfig AppConfiguration

func LoadConfig() error {
    // Use Beego's configuration parser
    conf, err := config.NewConfig("ini", "conf/app.conf")
    if err != nil {
        return err
    }

    // Load API configuration
    AppConfig.API.RapidAPIKey = conf.String("rapidapi::key")
    
    // Parse request timeout
    timeoutStr := conf.String("api::request_timeout")
    if timeout, err := strconv.Atoi(timeoutStr); err == nil {
        AppConfig.API.RequestTimeout = timeout
    }

    // Load rate limit configuration
    rateStr := conf.String("ratelimit::requests_per_second")
    if rate, err := strconv.ParseFloat(rateStr, 64); err == nil {
        AppConfig.RateLimit.RequestsPerSecond = rate
    }

    burstStr := conf.String("ratelimit::burst_size")
    if burst, err := strconv.Atoi(burstStr); err == nil {
        AppConfig.RateLimit.BurstSize = burst
    }

    return nil
}