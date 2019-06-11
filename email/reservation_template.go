package formatter

import (
	"strconv"
	"strings"

	"github.com/sylank/lavender-commons-go/utils"
)

// EmailTemplate ...
type EmailTemplate struct {
	rawText          string
	email            string
	name             string
	deletionURL      string
	reservationID    string
	fromDate         string
	toDate           string
	message          string
	costValue        int
	depositCostValue int
}

// InitEmail ...
func (template *EmailTemplate) InitEmail(fileName string) {
	template.rawText = string(utils.ReadBytesFromFile(fileName))
}

// SetEmail ...
func (template *EmailTemplate) SetEmail(email string) {
	template.email = email
}

// SetDeletionURL ...
func (template *EmailTemplate) SetDeletionURL(url string) {
	template.deletionURL = url
}

// SetName ...
func (template *EmailTemplate) SetName(name string) {
	template.name = name
}

// SetReservationID ...
func (template *EmailTemplate) SetReservationID(reservationID string) {
	template.reservationID = reservationID
}

// SetFromDate ...
func (template *EmailTemplate) SetFromDate(fromDate string) {
	template.fromDate = fromDate
}

// SetToDate ...
func (template *EmailTemplate) SetToDate(toDate string) {
	template.toDate = toDate
}

// SetMessage ...
func (template *EmailTemplate) SetMessage(message string) {
	template.message = message
}

// SetCostValue ...
func (template *EmailTemplate) SetCostValue(costValue int) {
	template.costValue = costValue
}

// SetDepositCostValue ...
func (template *EmailTemplate) SetDepositCostValue(depositCostValue int) {
	template.depositCostValue = depositCostValue
}

// GenerateEmailText ...
func (template *EmailTemplate) GenerateEmailText() string {
	var tmpText = template.rawText
	r := strings.NewReplacer(
		"<email>", template.email,
		"<url>", template.deletionURL,
		"<name>", template.name,
		"<reservationId>", template.reservationID,
		"<fromDate>", template.fromDate,
		"<toDate>", template.toDate,
		"<message>", template.message,
		"<costValue>", strconv.Itoa(template.costValue),
		"<depositCost>", strconv.Itoa(template.depositCostValue))

	return r.Replace(tmpText)
}
