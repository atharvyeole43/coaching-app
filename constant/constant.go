package constant

const DateFormat string = "2006-01-02"

var ErrorConstants = struct {
	InvalidRequestPayload        string
	ConversionErrorStringToFloat string
	UnauthorizedAccess           string
	ErrorFetchingExistingRecord  string
	ErrorParsingDate             string
}{
	InvalidRequestPayload:        "Invalid request payload",
	ConversionErrorStringToFloat: "Error while conversion of string to float ",
	UnauthorizedAccess:           "Unauthorized Access",
	ErrorFetchingExistingRecord:  "error while getting the Existing record for revenue :",
	ErrorParsingDate:             "Getting Error while Parsing the Date :",
}

const (
	ErrInvalidUserID = "Invalid user ID"
	ErrUserIDMissing = "Authorization failed: User ID is missing"
)
