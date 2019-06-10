package messaging

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SendTransactionalEmail ...
func SendTransactionalEmail(message string, queueName string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	qURL, err := getQueueURL(queueName, svc)
	if err != nil {
		log.Println("Failed to fetch queue URL")
		return err
	}

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody: aws.String(message),
		QueueUrl:    &qURL,
	})

	if err != nil {
		log.Println("Failed to send message", err)
		return err
	}

	log.Println("Success", *result.MessageId)
	return nil
}

func getQueueURL(ququeName string, svc *sqs.SQS) (string, error) {
	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(ququeName),
	})

	if err != nil {
		log.Println("Error", err)
		return "", err
	}

	log.Println("Success", *result.QueueUrl)

	return *result.QueueUrl, nil
}
