package properties

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sylank/lavender-commons-go/utils"
)

// Secrets ...
type Secrets struct {
	ReCaptchaServerSecret string `json:"reCaptchaServerSecret"`
	EncriptionKey         string `json:"encriptionKey"`
}

// DynamoProperties ...
type DynamoProperties struct {
	Region    string               `json:"region"`
	TableInfo map[string]TableInfo `json:"tableInfo"`
}

// TableInfo ...
type TableInfo struct {
	TableName string `json:"tableName"`
}

// ReadSecretProperties ...
func ReadSecretProperties(fileName string) (*Secrets, error) {
	data := utils.ReadBytesFromFile(fileName)
	var obj Secrets
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		log.Println(fmt.Sprintf("Error while reading file, filename: %s",fileName), err)

		return nil, err
	}

	return &obj, nil
}

// ReadDynamoProperties ...
func ReadDynamoProperties(fileName string) (*DynamoProperties, error) {
	data := utils.ReadBytesFromFile(fileName)
	var obj DynamoProperties
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		log.Println(fmt.Sprintf("Error while reading file, filename: %s",fileName), err)

		return nil, err
	}

	return &obj, nil
}

// GetTableName ...
func (properties *DynamoProperties) GetTableName(customTableName string) string {
	tableName := properties.TableInfo[customTableName].TableName

	return fmt.Sprintf("lavender-%s-%s", GetEnvironmentName(), tableName)
}

// GetEnvironmentName ...
func GetEnvironmentName() string {
	return os.Getenv("environment_name")
}
