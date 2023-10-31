package sapi

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

// ----------------------------------------------------------------------------
// envelope status error
// ----------------------------------------------------------------------------

// EnvelopeStatusError defines the structure to manipulate an error structure
// that hold the information of an execution error and be assigned to the
// response status error list.
type EnvelopeStatusError struct {
	Service  int    `json:"-" xml:"-"`
	Endpoint int    `json:"-" xml:"-"`
	Param    int    `json:"-" xml:"-"`
	Error    string `json:"-" xml:"-"`
	Code     string `json:"code" xml:"code"`
	Message  string `json:"message" xml:"message"`
}

// NewEnvelopeStatusError instantiates a new error instance.
func NewEnvelopeStatusError(
	e any,
	msg string,
) *EnvelopeStatusError {
	return (&EnvelopeStatusError{
		Error:   fmt.Sprintf("%v", e),
		Message: msg,
	}).compose()
}

// SetService assigns a service code value to the error.
func (e *EnvelopeStatusError) SetService(
	val int,
) *EnvelopeStatusError {
	e.Service = val
	return e.compose()
}

// SetEndpoint assigns an endpoint code value to the error.
func (e *EnvelopeStatusError) SetEndpoint(
	val int,
) *EnvelopeStatusError {
	e.Endpoint = val
	return e.compose()
}

// SetParam assigns a parameter code value to the error.
func (e *EnvelopeStatusError) SetParam(
	param int,
) *EnvelopeStatusError {
	e.Param = param
	return e.compose()
}

// SetError assigns a error code value to the error.
func (e *EnvelopeStatusError) SetError(
	err any,
) *EnvelopeStatusError {
	e.Error = fmt.Sprintf("%v", err)
	return e.compose()
}

// SetMessage assigns a message to the error.
func (e *EnvelopeStatusError) SetMessage(
	msg string,
) *EnvelopeStatusError {
	e.Message = msg
	return e
}

// GetCode retrieves the composed code of the error
func (e *EnvelopeStatusError) GetCode() string {
	return e.Code
}

// GetMessage retrieves the message associated to the error
func (e *EnvelopeStatusError) GetMessage() string {
	return e.Message
}

func (e *EnvelopeStatusError) compose() *EnvelopeStatusError {
	cb := strings.Builder{}
	// compose the service section of the code
	if e.Service != 0 {
		cb.WriteString(fmt.Sprintf("s:%d", e.Service))
	}
	// compose the endpoint section of the code
	if e.Endpoint != 0 {
		if cb.Len() != 0 {
			cb.WriteString(".")
		}
		cb.WriteString(fmt.Sprintf("e:%d", e.Endpoint))
	}
	// compose the param section of the code
	if e.Param != 0 {
		if cb.Len() != 0 {
			cb.WriteString(".")
		}
		cb.WriteString(fmt.Sprintf("p:%d", e.Param))
	}
	// compose the error section of the code
	if e.Error != "" {
		if cb.Len() != 0 {
			cb.WriteString(".")
		}

		if i, err := strconv.Atoi(e.Error); err != nil {
			cb.WriteString(e.Error)
		} else {
			cb.WriteString(fmt.Sprintf("c:%d", i))
		}
	}
	// assign the written code string to the error code structure parameter
	e.Code = cb.String()
	return e
}

// ----------------------------------------------------------------------------
// envelope status error list
// ----------------------------------------------------------------------------

// EnvelopeStatusErrorList defines a type of data  that holds a list
// of error structures.
type EnvelopeStatusErrorList []*EnvelopeStatusError

// MarshalXML serialize the error list into a xml string
func (s EnvelopeStatusErrorList) MarshalXML(
	e *xml.Encoder,
	start xml.StartElement,
) error {
	// encode the list starting tag
	_ = e.EncodeToken(start)
	// iterate through all the stored error
	for _, v := range s {
		// create the iterated error starting tag name
		name := xml.Name{Space: "", Local: "error"}
		// encode the error instance tag with the code and message attributes
		_ = e.EncodeToken(xml.StartElement{
			Name: name,
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "code"}, Value: v.Code},
				{Name: xml.Name{Local: "message"}, Value: v.Message},
			},
		})
		// encode the terminating error tag
		_ = e.EncodeToken(xml.EndElement{Name: name})
	}
	// encode the terminating list tag
	_ = e.EncodeToken(xml.EndElement{Name: start.Name})
	_ = e.Flush()
	return nil
}

// ----------------------------------------------------------------------------
// envelope status
// ----------------------------------------------------------------------------

// EnvelopeStatus defines the structure to manipulate a
// response status information structure.
type EnvelopeStatus struct {
	Success bool                    `json:"success" xml:"success"`
	Errors  EnvelopeStatusErrorList `json:"error" xml:"error"`
}

// NewEnvelopeStatus instantiates a new request result status structure.
func NewEnvelopeStatus() *EnvelopeStatus {
	return &EnvelopeStatus{
		Success: true,
		Errors:  EnvelopeStatusErrorList{},
	}
}

// AddError append a new error to the status error list
func (s *EnvelopeStatus) AddError(
	e *EnvelopeStatusError,
) *EnvelopeStatus {
	s.Success = false
	s.Errors = append(s.Errors, e)
	return s
}

// SetService assign a service code to all stored error.
func (s *EnvelopeStatus) SetService(
	val int,
) *EnvelopeStatus {
	for i := range s.Errors {
		s.Errors[i] = s.Errors[i].SetService(val)
	}
	return s
}

// SetEndpoint assign an endpoint code to all stored error.
func (s *EnvelopeStatus) SetEndpoint(
	val int,
) *EnvelopeStatus {
	for i := range s.Errors {
		s.Errors[i] = s.Errors[i].SetEndpoint(val)
	}
	return s
}

// ----------------------------------------------------------------------------
// envelope list report
// ----------------------------------------------------------------------------

// EnvelopeListReport defines the structure of a response list report
// containing all the request information, but also the total amount of
// filtering records and links for the previous and next pages
type EnvelopeListReport struct {
	Search string `json:"search" xml:"search"`
	Start  uint   `json:"start" xml:"start"`
	Count  uint   `json:"count" xml:"count"`
	Total  uint   `json:"total" xml:"total"`
	Prev   string `json:"prev" xml:"prev"`
	Next   string `json:"next" xml:"next"`
}

// NewEnvelopeListReport instantiates a new response list report by
// populating the prev and next link information regarding the given
// filtering information
func NewEnvelopeListReport(
	search string,
	start,
	count,
	total uint,
) *EnvelopeListReport {
	// store the prev URL query parameters if the start value
	// is greater than zero
	prev := ""
	if start > 0 {
		// discover the previous page starting value
		nstart := uint(0)
		if count < start {
			nstart = start - count
		}
		// compose the URL prev page query parameters
		prev = fmt.Sprintf(
			"?search=%s&start=%d&count=%d",
			search,
			nstart,
			count,
		)
	}
	// store the next URL query parameters if the total number of
	// record are greater than the current start plus the number of
	// presented records
	next := ""
	if start+count < total {
		// compose the URL next page query parameters
		next = fmt.Sprintf(
			"?search=%s&start=%d&count=%d",
			search,
			start+count,
			count,
		)
	}
	// return the list report instance reference
	return &EnvelopeListReport{
		Search: search,
		Start:  start,
		Count:  count,
		Total:  total,
		Prev:   prev,
		Next:   next,
	}
}

// ----------------------------------------------------------------------------
// envelope
// ----------------------------------------------------------------------------

// Envelope identifies the structure of a response structured format.
type Envelope struct {
	XMLName    xml.Name            `json:"-" xml:"envelope"`
	StatusCode int                 `json:"-" xml:"-"`
	Status     *EnvelopeStatus     `json:"status" xml:"status"`
	ListReport *EnvelopeListReport `json:"report,omitempty" xml:"report,omitempty"`
	Data       interface{}         `json:"data,omitempty" xml:"data,omitempty"`
}

// NewEnvelope instantiates a new response data envelope structure
func NewEnvelope(
	statusCode int,
	data interface{},
	listReport ...*EnvelopeListReport,
) *Envelope {
	// initialize the envelope structure
	env := &Envelope{
		StatusCode: statusCode,
		Status:     NewEnvelopeStatus(),
		ListReport: nil,
		Data:       data,
	}
	// assign the list report if given as argument
	if len(listReport) > 0 && listReport[0] != nil {
		env.ListReport = listReport[0]
	}
	return env
}

// GetStatusCode returned the stored enveloped response status code
func (s *Envelope) GetStatusCode() int {
	return s.StatusCode
}

// SetService assign the service identifier to all stored error codes
func (s *Envelope) SetService(
	val int,
) *Envelope {
	s.Status = s.Status.SetService(val)
	return s
}

// SetEndpoint assign the endpoint identifier to all stored error codes
func (s *Envelope) SetEndpoint(
	val int,
) *Envelope {
	s.Status = s.Status.SetEndpoint(val)
	return s
}

// SetListReport assign the list report to the envelope
func (s *Envelope) SetListReport(
	listReport *EnvelopeListReport,
) *Envelope {
	s.ListReport = listReport
	return s
}

// AddError add a new error to the response envelope instance
func (s *Envelope) AddError(
	e *EnvelopeStatusError,
) *Envelope {
	s.Status = s.Status.AddError(e)
	return s
}
