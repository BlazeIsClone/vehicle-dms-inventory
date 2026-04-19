package vehicle

import "strings"

type ValidationErrors map[string]string

func (e ValidationErrors) HasErrors() bool { return len(e) > 0 }

type vehicleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (req vehicleRequest) validate() ValidationErrors {
	errs := make(ValidationErrors)

	name := strings.TrimSpace(req.Name)
	switch {
	case name == "":
		errs["name"] = "required"
	case len(name) > 255:
		errs["name"] = "must not exceed 255 characters"
	}

	if desc := strings.TrimSpace(req.Description); len(desc) > 1000 {
		errs["description"] = "must not exceed 1000 characters"
	}

	return errs
}
