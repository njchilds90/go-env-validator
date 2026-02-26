package envvalidator

// Schema returns a slice of FieldSchema values that describe every declared
// field in the Validator. The output is deterministic: fields appear in the
// same order they were passed to New.
//
// This method is intended for tooling, documentation generators, and AI agents
// that need a machine-readable description of an application's configuration
// contract.
//
// Example:
//
//	v := envvalidator.New(
//	    envvalidator.Field{Key: "PORT", Kind: envvalidator.KindInteger, Default: "8080", Description: "HTTP server port"},
//	)
//	schema := v.Schema()
//	// schema[0].Key == "PORT", schema[0].Kind == "integer"
func (v *Validator) Schema() []FieldSchema {
	out := make([]FieldSchema, len(v.fields))
	for i, f := range v.fields {
		kind := f.Kind
		if kind == "" {
			kind = KindString
		}
		allowed := f.AllowedValues
		if allowed == nil {
			allowed = []string{}
		}
		out[i] = FieldSchema{
			Key:           f.Key,
			Kind:          string(kind),
			Required:      f.Required,
			Default:       f.Default,
			Description:   f.Description,
			AllowedValues: allowed,
		}
	}
	return out
}
