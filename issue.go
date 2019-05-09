package score

import "fmt"

const (
	UnableToCreateHTTPRequest        string = "SCORE_E-0100"
	UnableToExecuteHTTPGETRequest    string = "SCORE_E-0200"
	InvalidAPIRespHTTPStatusCode     string = "SCORE_E-0300"
	UnableToReadBodyFromHTTPResponse string = "SCORE_E-0400"
	APIErrorResponseFound            string = "SCORE_E-0500"
	NoAPIKeyProvidedInCodeOrEnv      string = "SCORE_E-0600"
	SecretManagementError            string = "SCORE_E-0700"
)

// Issue is a structured problem identification with context information
type Issue interface {
	IssueContext() interface{} // this will be the scores object plus location (item index, etc.), it's kept generic so it doesn't require package dependency
	IssueCode() string         // useful to uniquely identify a particular code
	Issue() string             // the

	IsError() bool   // this issue is an error
	IsWarning() bool // this issue is a warning
}

// Issues packages multiple issues into a container
type Issues interface {
	ErrorsAndWarnings() []Issue
	IssueCounts() (uint, uint, uint)
	HandleIssues(errorHandler func(Issue), warningHandler func(Issue))
}

type issue struct {
	APIEndpoint    string `json:"context"`
	Code           string `json:"code"`
	Message        string `json:"message"`
	IsIssueAnError bool   `json:"isError"`
}

func NewIssue(apiEndpoint string, code string, message string, isError bool) Issue {
	result := new(issue)
	result.APIEndpoint = apiEndpoint
	result.Code = code
	result.Message = message
	result.IsIssueAnError = isError
	return result
}

func NewHTTPResponseIssue(apiEndpoint string, httpRespStatusCode int, message string, isError bool) Issue {
	result := new(issue)
	result.APIEndpoint = apiEndpoint
	result.Code = fmt.Sprintf("%s-HTTP-%d", InvalidAPIRespHTTPStatusCode, httpRespStatusCode)
	result.Message = message
	result.IsIssueAnError = isError
	return result
}

func (i issue) IssueContext() interface{} {
	return i.APIEndpoint
}

func (i issue) IssueCode() string {
	return i.Code
}

func (i issue) Issue() string {
	return i.Message
}

func (i issue) IsError() bool {
	return i.IsIssueAnError
}

func (i issue) IsWarning() bool {
	return !i.IsIssueAnError
}

// Error satisfies the Go error contract
func (i issue) Error() string {
	return i.Message
}
