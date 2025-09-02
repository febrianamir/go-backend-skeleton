package request

type SendEmailPayload struct {
	To           []string       `json:"to"`
	TemplateName string         `json:"template_name"`
	TemplateData map[string]any `json:"template_data"`
	Subject      string         `json:"subject"`
}
