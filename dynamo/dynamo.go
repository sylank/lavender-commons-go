package dynamo

import (
	"log"
	"strconv"

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

// DeletionInsertModel ...
type DeletionInsertModel struct {
	UserID        string `json:"UserId"`
	ReservationID string `json:"ReservationId"`
	Type          string `json:"Type"`
	Message       string `json:"Message"`
}

// ReservationModel ...
type ReservationModel struct {
	ReservationID string
	FromDate      string
	ToDate        string
	UserID        string
	Deleted       bool
}

// ReservationDynamoModel ...
type ReservationDynamoModel struct {
	ReservationID string
	FromDate      string
	ToDate        string
	UserID        string
	Deleted       string
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
	result, err := CustomQuery("Email", email, userTableName, proj)
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

			return nil, err
		}

		if item.Email == email {
			log.Println("Record found!")
			return &item, nil
		}
	}

	log.Println("Record not found!")
	return nil, nil
}

// CustomQuery ...
func CustomQuery(clumnName string, value string, table string, proj expression.ProjectionBuilder) (*dynamodb.ScanOutput, error) {
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

// QueryUserByUserID ...
func QueryUserByUserID(userID string) (*UserModel, error) {
	usersTableName := properties.GetTableName("userData")

	proj := expression.NamesList(expression.Name("FullName"), expression.Name("Email"), expression.Name("Phone"), expression.Name("UserId"))
	result, err := CustomQuery("UserId", userID, usersTableName, proj)

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

			return nil, err
		}

		return &item, nil
	}

	log.Println("Record not found!")
	return nil, nil
}

// ClearUserData ...
func ClearUserData(userID string) error {
	userTableName := properties.GetTableName("userData")

	proj := expression.NamesList(expression.Name("FullName"), expression.Name("Email"), expression.Name("Phone"), expression.Name("UserId"))
	result, err := CustomQuery("UserId", userID, userTableName, proj)
	if err != nil {
		return err
	}

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

// InsertDeletionTypeTable ...
func InsertDeletionTypeTable(deletionModel *DeletionInsertModel, tableName string) error {
	log.Print("Insert deletion data: ")
	log.Println(deletionModel)

	av, err := dynamodbattribute.MarshalMap(deletionModel)
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

// QueryReservationTypeTable ...
func QueryReservationTypeTable(reservationID string, table string) ([]ReservationModel, error) {
	log.Println("Query data with reservationId: " + reservationID)
	var retData []ReservationModel
	proj := expression.NamesList(expression.Name("ReservationId"), expression.Name("FromDate"), expression.Name("ToDate"), expression.Name("UserId"), expression.Name("Deleted"))
	result, err := CustomQuery("ReservationId", reservationID, table, proj)
	if err != nil {
		log.Println("Query API call failed:", err)
		return nil, err
	}

	log.Println("Result array:")
	log.Println(result.Items)
	for _, i := range result.Items {
		log.Println("Marshalling:")
		log.Println(i)
		item := ReservationDynamoModel{}
		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Println("Got error unmarshalling:", err)
			return nil, err
		}

		deleted, err := strconv.ParseBool(item.Deleted)
		if err != nil {
			return nil, err
		}

		retData = append(retData, ReservationModel{ReservationID: item.ReservationID, FromDate: item.FromDate, ToDate: item.ToDate, UserID: item.UserID, Deleted: deleted})
	}

	return retData, err
}

// InsertReservationTypeTable ...
func InsertReservationTypeTable(reservationModel *ReservationModel, table string) {
	av, err := dynamodbattribute.MarshalMap(reservationModel)
	if err != nil {
		log.Println("Got error marshalling new reservationModel item:")
		log.Println(err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	_, err = client.PutItem(input)
	if err != nil {
		log.Println("Got error calling PutItem:")
		log.Println(err.Error())
	}

	log.Println("Item inserted with reservationId: " + reservationModel.ReservationID)
}

// DeleteReservationType ...
func DeleteReservationType(reservationID string, table string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ReservationId": {
				S: aws.String(reservationID),
			},
		},
		TableName: aws.String(table),
	}

	_, err := client.DeleteItem(input)
	if err != nil {
		log.Println("Got error calling DeleteItem", err)

		return err
	}

	log.Println("Item deleted with reservationId: " + reservationID)

	return nil
}

// UpdateDeletedReservationStatus ...
func UpdateDeletedReservationStatus(reservationID string, userID string, table string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				BOOL: aws.Bool(true),
			},
		},
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"ReservationId": {
				S: aws.String(reservationID),
			},
			"UserId": {
				S: aws.String(userID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Deleted = :r"),
	}

	_, updateError := client.UpdateItem(input)
	if updateError != nil {
		log.Println(updateError.Error())

		return updateError
	}

	log.Println("Record updated")

	return nil
}
