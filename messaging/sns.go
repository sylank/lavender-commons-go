package messaging

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// PublishMessage ..
func PublishMessage(message string, subject string) error {
	svc := sns.New(session.New())

	params := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(os.Getenv("EMAIL_SNS_TOPIC_ARN")),
		Subject:  aws.String(subject),
	}

	resp, err := svc.Publish(params)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println(resp)
	return nil
}