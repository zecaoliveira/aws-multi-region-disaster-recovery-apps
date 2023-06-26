package parameter_store

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/msfidelis/rest-api-demo/pkg/log"
	"github.com/msfidelis/rest-api-demo/pkg/memory_cache"
)

func GetParamValue(parameter string, cache_time int64) (string, error) {

	m := memory_cache.GetInstance()
	log := log.Instance()

	if cache_time > 0 {

		value, found := m.Get(parameter)

		if found {
			log.Info().
				Str("Parameter Store", parameter).
				Str("AWS_REGION", os.Getenv("AWS_REGION")).
				Msg("Returning parameter store value from cache")
			return fmt.Sprint(value), nil
		} else {
			log.Info().
				Str("Parameter Store", parameter).
				Str("AWS_REGION", os.Getenv("AWS_REGION")).
				Msg("Parameter value don't found in cache")
		}

	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return "", err
	}

	svc := ssm.New(sess)

	result, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(parameter),
		WithDecryption: aws.Bool(false),
	})

	fmt.Println(*result.Parameter.Value)

	if err != nil {
		return "", err
	}

	if cache_time > 0 {
		log.Info().
			Str("Parameter Store", parameter).
			Str("AWS_REGION", os.Getenv("AWS_REGION")).
			Int64("Cache_Time_Seconds", cache_time).
			Msg("Saving parameter store value on local cache")

		m.Set(parameter, *result.Parameter.Value, time.Second*time.Duration(cache_time))
	}

	return fmt.Sprint(*result.Parameter.Value), nil

}
