package printer

import (
	"encoding/json"
	"fmt"
	"github.com/rozdolsky33/ocloud/internal/app"
)

// MarshalToJSON marshals any data to JSON and prints it to stdout
// It takes an app context for logging errors and any data to be marshaled
func MarshalToJSON(data interface{}, appCtx *app.AppContext) {
	// Marshal the data to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		appCtx.Logger.Error(err, "Failed to marshal data to JSON")
		return
	}

	// Print the JSON data
	fmt.Println(string(jsonData))
}
