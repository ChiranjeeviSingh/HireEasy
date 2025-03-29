package services

import (
	"math"
	"strings"
)

// MockATSScore calculates ATS score based on candidate and job-required skills
func MockATSScore(candidateSkills []string, requiredSkills []string) int {
	if len(requiredSkills) == 0 {
		return 50 // Default score if no required skills are provided
	}

	// Normalize skills (convert to lowercase)
	normalize := func(skills []string) map[string]bool {
		skillMap := make(map[string]bool)
		for _, skill := range skills {
			skillMap[strings.ToLower(skill)] = true
		}
		return skillMap
	}

	candidateSkillsMap := normalize(candidateSkills)
	requiredSkillsMap := normalize(requiredSkills)

	// Count matching skills
	matchCount := 0
	for skill := range candidateSkillsMap {
		if requiredSkillsMap[skill] {
			matchCount++
		}
	}

	// Calculate match percentage
	matchPercentage := (float64(matchCount) / float64(len(requiredSkills))) * 100

	// Scale score between 50 and 100
	score := 50 + int(math.Round(matchPercentage/2))

	if score > 100 {
		score = 100
	}

	return score
}
