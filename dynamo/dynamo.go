package dynamo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pkg/errors"

	props "github.com/sylank/lavender-commons-go/properties"
)

// UserModel ...
type UserModel struct {
	UserID   string
	FullName string
	Email    string
	Phone    string
	Inserted int
}

// ReservationModel ...
type ReservationModel struct {
	ReservationID string
	Email         string
}

// TODO: refactor the Deletion model and the DeletionInsertionModel struct
type DeletionModel struct {
	UserId        string `json:"userId"`
	ReservationId string `json:"reservationId"`
	Type          string `json:"type"`
	Message       string `json:"message"`
}

var client *dynamodb.DynamoDB
var properties *props.DynamoProperties

// CreateConnection ...
func CreateConnection(dynamoProperties *props.DynamoProperties) *dynamodb.DynamoDB {
	properties = dynamoProperties
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client = dynamodb.New(sess)

	return client
}

// GetDynamoClient ..
func GetDynamoClient() *dynamodb.DynamoDB {
	return client
}

// IsUserStored ...
func IsUserStored(email string) (*UserModel, error) {
	userTableName := properties.GetTableName("userData")
	log.Println(userTableName)

	proj := expression.NamesList(expression.Name("FullName"), expression.Name("Email"), expression.Name("Phone"), expression.Name("UserId"))
	result, err := queryItems(email, userTableName, proj)
	if err != nil {
		log.Println("Query API call failed:")
		log.Println((err.Error()))

		return nil, err
	}

	for _, i := range result.Items {
		item := UserModel{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Println("Got error unmarshalling:")
			log.Println(err.Error())
		}

		if item.Email == email {
			log.Println("Record found!")
			return &item, nil
		}
	}

	log.Println("Record not found!")
	return nil, nil
}

//TODO: refactor
func queryItems(email string, table string, proj expression.ProjectionBuilder) (*dynamodb.ScanOutput, error) {
	return customQuery("Email", email, table, proj)
}

func queryItemsByUserId(userId string, table string, proj expression.ProjectionBuilder) (*dynamodb.ScanOutput, error) {
	return customQuery("UserId", userId, table, proj)
}

func customQuery(clumnName string, value string, table string, proj expression.ProjectionBuilder) (*dynamodb.ScanOutput, error) {
	filt := expression.Name(clumnName).Equal(expression.Value(value))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Println("Got error building expression:")
		log.Println(err.Error())

		return nil, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(table),
	}

	// Make the DynamoDB Query API call
	result, err := client.Scan(params)
	if err != nil {
		log.Println("Query API call failed:")
		log.Println((err.Error()))
	}

	return result, err
}

// ClearUserData ...
func ClearUserData(userID string) error {
	userTableName := properties.GetTableName("userData")

	proj := expression.NamesList(expression.Name("FullName"), expression.Name("Email"), expression.Name("Phone"), expression.Name("UserId"))
	result, err := queryItemsByUserId(userID, userTableName, proj)

	userEmail := ""

	for _, i := range result.Items {
		item := UserModel{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Println("Got error unmarshalling:")
			log.Println(err.Error())

			return err
		}

		if item.UserID == userID {
			log.Println("Record found!")
			userEmail = item.Email
		}
	}

	if len(userEmail) > 0 {
		log.Println("User id: " + userID)
		log.Println("Email: " + userEmail)
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":r": {
					S: aws.String("#CLEARED#"),
				},
			},
			TableName: aws.String(userTableName),
			Key: map[string]*dynamodb.AttributeValue{
				"UserId": {
					S: aws.String(userID),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set Email = :r, FullName = :r, Phone = :r"),
		}

		_, updateError := client.UpdateItem(input)
		if updateError != nil {
			log.Println(updateError.Error())
			return errors.New("Updating error")
		}

		log.Println("Record updated")
	} else {
		return errors.New("Email not found")
	}

	return nil
}

// ClearReservationData ...
func ClearReservationData(email string) error {
	reservationTableName := properties.GetTableName("reservations")

	proj := expression.NamesList(expression.Name("ReservationId"), expression.Name("Email"))
	result, err := queryItems(email, reservationTableName, proj)

	for _, i := range result.Items {
		item := ReservationModel{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Println("Got error unmarshalling:")
			log.Println(err.Error())

			return err
		}

		if item.Email == email {
			log.Println("Record found!")
			reservationID := item.ReservationID

			log.Println("Reservation id for mail: " + reservationID)

			input := &dynamodb.UpdateItemInput{
				ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
					":r": {
						S: aws.String("#CLEARED#"),
					},
				},
				TableName: aws.String(reservationTableName),
				Key: map[string]*dynamodb.AttributeValue{
					"ReservationId": {
						S: aws.String(reservationID),
					},
				},
				ReturnValues:     aws.String("UPDATED_NEW"),
				UpdateExpression: aws.String("set Email = :r"),
			}

			_, updateError := client.UpdateItem(input)
			if updateError != nil {
				log.Println(updateError.Error())
				return updateError
			}
		}
	}

	log.Println("Record updated")

	return nil
}

// InsertDeletionTypeTable ...
func InsertDeletionTypeTable(deletionModel *DeletionModel, tableName string) error {
	log.Print("Insert deletion data: ")
	log.Println(deletionModel)

	type DeletionInsertionModel struct {
		UserID        string
		ReservationID string
		Type          string
		Message       string
	}

	av, err := dynamodbattribute.MarshalMap(&DeletionInsertionModel{UserID: deletionModel.UserId, ReservationID: deletionModel.ReservationId, Type: deletionModel.Type, Message: deletionModel.Message})
	if err != nil {
		log.Println("Got error marshalling new reservationModel item:")
		log.Println(err.Error())

		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = GetDynamoClient().PutItem(input)
	if err != nil {
		log.Println("Got error calling PutItem:")
		log.Println(err.Error())

		return err
	}

	return nil
}
