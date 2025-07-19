package response

type ProblemDetails struct {
	Type      string `json:"type,omitempty"`
	Title     string `json:"title"`
	ErrorCode string `json:"error_code"`
	Detail    string `json:"detail,omitempty"`
	Instance  string `json:"instance,omitempty"`
}

func New(errorCode string, title string, detail ...string) *ProblemDetails {
	if len(detail) > 0 {
		return &ProblemDetails{ErrorCode: errorCode, Title: title, Detail: detail[0]}
	}
	return &ProblemDetails{ErrorCode: errorCode, Title: title, Detail: ""}
}

func (p *ProblemDetails) WithType(t string) *ProblemDetails {
	p.Type = t
	return p
}
func (p *ProblemDetails) WithInstance(i string) *ProblemDetails {
	p.Instance = i
	return p
}
