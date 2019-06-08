package test

import (
	"log"
	"testing"

	props "github.com/sylank/lavender-commons-go/properties"
)

const (
	EXPECTED_TABLE_NAME = "user_data"
)

func Test_ReadAndLoadTableName(t *testing.T) {
	properties, err := props.ReadDynamoProperties("../config/database_properties.json")
	if err != nil {
		t.Error("Failed to load file")
	}

	tableName := properties.GetTableName("user_data")
	log.Println("Table name: " + tableName)

	if tableName != EXPECTED_TABLE_NAME {
		t.Error("Table name does not match")
	}
}
