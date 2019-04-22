package score

import (
	"fmt"
	"os"
)

func lookupAPIKeyInEnv(context string, apiKeyEnvVarName string) (string, bool, Issue) {
	apiKey, apiKeyOK := os.LookupEnv(apiKeyEnvVarName)
	if !apiKeyOK {
		issue := newIssue(context, NoAPIKeyProvidedInCodeOrEnv, fmt.Sprintf("Environment variable %q does not contain API key", apiKeyEnvVarName), true)
		return apiKey, false, issue
	}
	return apiKey, true, nil
}
