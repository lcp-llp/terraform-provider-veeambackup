package aws

type AWSAccessKeyAuth struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}