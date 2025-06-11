package gazelle_pyproject

import (
	"regexp"
)

var re = regexp.MustCompile(`(?mi)^[a-zA-Z0-9](?:[a-zA-Z0-9_\-\.]*[a-zA-Z0-9]|[a-zA-Z0-9]*)`)

func extractPackageNameFromRequirement(requirements string) string {
	match := re.FindString(requirements)
	return match
}

func ExtractPackageNamesFromRequirements(requirements []string) []string {
	packageNames := make([]string, 0)
	for _, requirement := range requirements {
		packageNames = append(packageNames, extractPackageNameFromRequirement(requirement))
	}
	return packageNames
}
