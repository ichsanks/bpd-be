package configs

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config is a struct that will receive configuration options via environment
// variables.
type Config struct {
	App struct {
		CORS struct {
			AllowCredentials bool     `mapstructure:"ALLOW_CREDENTIALS"`
			AllowedHeaders   []string `mapstructure:"ALLOWED_HEADERS"`
			AllowedMethods   []string `mapstructure:"ALLOWED_METHODS"`
			AllowedOrigins   []string `mapstructure:"ALLOWED_ORIGINS"`
			Enable           bool     `mapstructure:"ENABLE"`
			MaxAgeSeconds    int      `mapstructure:"MAX_AGE_SECONDS"`
		}
		File struct {
			Dir             string `mapstructure:"DIR"`
			FotoProfile     string `mapstructure:"FOTO_PROFILE"`
			Templates       string `mapstructure:"TEMPLATE_REPORT"`
			FotoDinas       string `mapstructure:"FOTO_DINAS"`
			DokumenBpd      string `mapstructure:"DOKUMEN_BPD"`
			Bpd             string `mapstructure:"BPD"`
			PerjalananDinas string `mapstructure:"PERJALANAN_DINAS"`
			TtdPegawai      string `mapstructure:"TTD_PEGAWAI"`
			ImageBranch     string `mapstructure:"IMAGE_BRANCH"`
			DokumenSppd     string `mapstructure:"DOKUMEN_SPPD"`
		}
		APIExternal struct {
			RecaptchaSecretKey string `mapstructure:"RECAPTCHA_SECRETKEY"`
			UsernameEsign      string `mapstructure:"USERNAME_ESIGN"`
			PasswordEsign      string `mapstructure:"PASSWORD_ESIGN"`
			IpEsign            string `mapstructure:"IP_ESIGN"`
			IpSuperman         string `mapstructure:"IP_SUPERMAN"`
			IpSupermanDev      string `mapstructure:"IP_SUPERMAN_DEV"`
			IpInoPtpn12        string `mapstructure:"IP_INO_PTPN12"`
			IpNadineSppd       string `mapstructure:"IP_NADINESPPD"`
			IpNadineBpd        string `mapstructure:"IP_NADINEBPD"`
		}
		Name     string `mapstructure:"NAME"`
		Revision string `mapstructure:"REVISION"`
		URL      string `mapstructure:"URL"`
	}

	Cache struct {
		Redis struct {
			Primary struct {
				Host     string `mapstructure:"HOST"`
				Port     string `mapstructure:"PORT"`
				Password string `mapstructure:"PASSWORD"`
			}
		}
	}

	DB struct {
		PostgreSQL struct {
			Read struct {
				Host     string `mapstructure:"HOST"`
				Port     string `mapstructure:"PORT"`
				Username string `mapstructure:"USER"`
				Password string `mapstructure:"PASSWORD"`
				Name     string `mapstructure:"NAME"`
				Timezone string `mapstructure:"TIMEZONE"`
			}
			Write struct {
				Host     string `mapstructure:"HOST"`
				Port     string `mapstructure:"PORT"`
				Username string `mapstructure:"USER"`
				Password string `mapstructure:"PASSWORD"`
				Name     string `mapstructure:"NAME"`
				Timezone string `mapstructure:"TIMEZONE"`
			}
		}
	}

	Event struct {
		Consumer struct {
			SQS struct {
				AccessKeyID       string `mapstructure:"ACCESS_KEY_ID"`
				BackoffSeconds    int    `mapstructure:"BACKOFF_SECONDS"`
				MaxMessage        int64  `mapstructure:"MAX_MESSAGE"`
				MaxRetries        int    `mapstructure:"MAX_RETRIES"`
				MaxRetriesConsume int    `mapstructure:"MAX_RETRIES_CONSUME"`
				Region            string `mapstructure:"REGION"`
				SecretAccessKey   string `mapstructure:"SECRET_ACCESS_KEY"`
				UUID              string `mapstructure:"UUID"`
				WaitTimeSeconds   int64  `mapstructure:"WAIT_TIME_SECONDS"`

				Topics struct {
					Order struct {
						Enabled bool   `mapstructure:"ENABLED"`
						URL     string `mapstructure:"URL"`
					} `mapstructure:"ORDER"`
				}
			}
		}

		Producer struct {
			SNS struct {
				AccessKeyID     string `mapstructure:"ACCESS_KEY_ID"`
				MaxRetries      int    `mapstructure:"MAX_RETRIES"`
				Region          string `mapstructure:"REGION"`
				SecretAccessKey string `mapstructure:"SECRET_ACCESS_KEY"`
				Topics          struct {
					FooCreated struct {
						ARN     string `mapstructure:"ARN"`
						Enabled bool   `mapstructure:"ENABLED"`
					} `mapstructure:"FOO_CREATED"`
				}
			}
		}
	}

	External struct {
		NotifPublisher struct {
			BaseURL  string `mapstructure:"BASE_URL"`
			Endpoint struct {
				Whatsapp string `mapstructure:"WHATSAPP"`
			}
			Source   string `mapstructure:"SOURCE"`
			Identity string `mapstructure:"IDENTITY"`
			Template string `mapstructure:"TEMPLATE"`
			Type     string `mapstructure:"TYPE"`
			Dial     struct {
				TimeoutSeconds   int64 `mapstructure:"TIMEOUT_SECONDS"`
				KeepAliveSeconds int64 `mapstructure:"KEEP_ALIVE_SECONDS"`
			}
			Client struct {
				TimeoutSeconds int64 `mapstructure:"TIMEOUT_SECONDS"`
			}
		} `mapstructure:"NOTIF_PUBLISHER"`

		Shipment struct {
			URL string `mapstructure:"URL"`
		} `mapstructure:"SHIPMENT"`
	} `mapstructure:"EXTERNAL"`

	Order struct {
		Tracking struct {
			LastMileKeywords []string `mapstructure:"LAST_MILE_KEYWORDS"`
		} `mapstructure:"TRACKING"`
	} `mapstructure:"ORDER"`

	Token struct {
		JWT struct {
			AccessToken   string `mapstructure:"ACCESS_TOKEN"`
			ExpiredInHour int    `mapstructure:"EXPIRED_IN_HOUR"`
		} `mapstructure:"JWT"`
	} `mapstructure:"TOKEN"`

	Server struct {
		Env      string `mapstructure:"ENV"`
		LogLevel string `mapstructure:"LOG_LEVEL"`
		Port     string `mapstructure:"PORT"`
		Shutdown struct {
			CleanupPeriodSeconds int64 `mapstructure:"CLEANUP_PERIOD_SECONDS"`
			GracePeriodSeconds   int64 `mapstructure:"GRACE_PERIOD_SECONDS"`
		}
	}

	Storage struct {
		S3 struct {
			BaseURL  string `mapstructure:"BASE_URL"`
			DataFile struct {
				Directory string `mapstructre:"DIRECTORY"`
			} `mapstructure:"DATAFILE"`
		} `mapstructure:"S3"`
	} `mapstructure:"STORAGE"`
}

var (
	conf Config
	once sync.Once
)

// Get are responsible to load env and get data an return the struct
func Get() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed reading config file")
	}

	once.Do(func() {
		log.Info().Msg("Service configuration initialized.")
		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	})

	return &conf
}
