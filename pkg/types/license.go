package types

import "fmt"

type LicenseReport struct {
	Image    string     `json:"image"`
	Licenses []*License `json:"licenses"`
}

type License struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

func (l *License) String() string {
	return fmt.Sprintf("%s: %s", l.Source, l.Name)
}

func (l *LicenseReport) Hash() string {
	return string(Hash(l))
}
