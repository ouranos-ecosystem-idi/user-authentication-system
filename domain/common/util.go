package common

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

// StringPtr
// Summary: This is function which returns the pointer to the string value.
// input: s(string): string value
// output: (*string) pointer to the string value
func StringPtr(s string) *string {
	return &s
}

// IsOutputDump
// Summary: This is function which checks whether the output dump is enabled.
// output: (bool) true if the output dump is enabled, false otherwise
func IsOutputDump() bool {
	executeEnv := os.Getenv("GO_ENV")
	outputDumpEnvs := []string{"local", "dev", "sbx", "qa", "stg", "uat"}
	for _, v := range outputDumpEnvs {
		if executeEnv == v {
			return true
		}
	}
	return false
}

// GenerateUUIDString
// Summary: This is function which generates the UUID string.
// input: n(int): number of UUID strings to generate
// output: (string) generated UUID string with comma separator
func GenerateUUIDString(n int) string {
	UUIDs := make([]string, 0, n)
	for i := 0; i < n; i++ {
		newUUID, _ := uuid.NewUUID()
		UUIDs = append(UUIDs, newUUID.String())
	}
	return strings.Join(UUIDs, ",")
}
