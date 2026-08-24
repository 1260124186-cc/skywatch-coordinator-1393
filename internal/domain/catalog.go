package domain

import "strings"

type TargetCatalog struct{ aliases map[string]string }

func NewTargetCatalog() *TargetCatalog {
	return &TargetCatalog{aliases: map[string]string{"aurora": "aurora arc", "meteor": "meteor trail", "nebula": "nebula field", "moon": "lunar limb"}}
}
func (c *TargetCatalog) Canonical(label string) string {
	normalized := strings.ToLower(strings.TrimSpace(label))
	if resolved, ok := c.aliases[normalized]; ok {
		return resolved
	}
	return normalized
}
func (c *TargetCatalog) Match(label string) bool {
	_, ok := c.aliases[strings.ToLower(strings.TrimSpace(label))]
	return ok
}
func (c *TargetCatalog) AddAlias(alias, canonical string) bool {
	alias = strings.ToLower(strings.TrimSpace(alias))
	canonical = strings.ToLower(strings.TrimSpace(canonical))
	if alias == "" || canonical == "" {
		return false
	}
	if _, exists := c.aliases[alias]; exists {
		return false
	}
	c.aliases[alias] = canonical
	return true
}
func (c *TargetCatalog) Aliases() map[string]string {
	result := make(map[string]string, len(c.aliases))
	for alias, canonical := range c.aliases {
		result[alias] = canonical
	}
	return result
}
func (c *TargetCatalog) Suggest(prefix string) []string {
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	out := make([]string, 0)
	seen := map[string]bool{}
	for alias, canonical := range c.aliases {
		if strings.HasPrefix(alias, prefix) && !seen[canonical] {
			out = append(out, canonical)
			seen[canonical] = true
		}
	}
	return out
}

type TargetProfile struct {
	Canonical  string `json:"canonical"`
	AliasCount int    `json:"aliasCount"`
	Known      bool   `json:"known"`
}

func (c *TargetCatalog) Profile(label string) TargetProfile {
	canonical := c.Canonical(label)
	count := 0
	for _, value := range c.aliases {
		if value == canonical {
			count++
		}
	}
	return TargetProfile{Canonical: canonical, AliasCount: count, Known: count > 0}
}
func (c *TargetCatalog) ValidateForCampaign(target string) bool {
	profile := c.Profile(target)
	return strings.TrimSpace(profile.Canonical) != ""
}
func NormalizeTarget(label string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(label))), " ")
}
func TargetWords(label string) []string {
	normalized := NormalizeTarget(label)
	if normalized == "" {
		return nil
	}
	return strings.Split(normalized, " ")
}
