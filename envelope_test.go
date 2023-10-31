package sapi

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_EnvelopeStatusError(t *testing.T) {
	t.Run("NewEnvelopeStatusError", func(t *testing.T) {
		t.Run("construct", func(t *testing.T) {
			code := 123
			msg := "message"
			e := NewEnvelopeStatusError(code, msg)
			if check := e.Service; check != 0 {
				t.Errorf("(%v) service value instead of zero", check)
			} else if check := e.Endpoint; check != 0 {
				t.Errorf("(%v) endpoint value instead of zero", check)
			} else if check := e.Param; check != 0 {
				t.Errorf("(%v) param value instead of zero", check)
			} else if check := e.Error; check != fmt.Sprintf("%d", code) {
				t.Errorf("(%v) when expecting (%v)", check, code)
			} else if check := e.Code; check != fmt.Sprintf("c:%d", code) {
				t.Errorf("(%v) when expecting (%v)", check, fmt.Sprintf("c:%d", code))
			} else if check := e.Message; check != msg {
				t.Errorf("(%v) when expecting (%v)", check, msg)
			}
		})

		t.Run("construct with string error code", func(t *testing.T) {
			code := "error code"
			msg := "message"
			e := NewEnvelopeStatusError(code, msg)
			if check := e.Service; check != 0 {
				t.Errorf("(%v) service value instead of zero", check)
			} else if check := e.Endpoint; check != 0 {
				t.Errorf("(%v) endpoint value instead of zero", check)
			} else if check := e.Param; check != 0 {
				t.Errorf("(%v) param value instead of zero", check)
			} else if check := e.Error; check != code {
				t.Errorf("(%v) when expecting (%v)", check, code)
			} else if check := e.Code; check != code {
				t.Errorf("(%v) when expecting (%v)", check, code)
			} else if check := e.Message; check != msg {
				t.Errorf("(%v) when expecting (%v)", check, msg)
			}
		})
	})

	t.Run("SetService", func(t *testing.T) {
		t.Run("assign", func(t *testing.T) {
			service := 123
			code := 456
			msg := "message"
			expected := fmt.Sprintf("s:%d.c:%d", service, code)
			err1 := NewEnvelopeStatusError(code, msg)
			err2 := err1.SetService(service)

			if check := err2.Service; check != service {
				t.Errorf("(%v) when expecting (%v)", check, service)
			} else if check := err2.Code; check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})
	})

	t.Run("SetEndpoint", func(t *testing.T) {
		t.Run("assign", func(t *testing.T) {
			endpoint := 123
			code := 456
			msg := "message"
			expected := fmt.Sprintf("e:%d.c:%d", endpoint, code)
			err1 := NewEnvelopeStatusError(code, msg)
			err2 := err1.SetEndpoint(endpoint)

			if check := err2.Endpoint; check != endpoint {
				t.Errorf("(%v) when expecting (%v)", check, endpoint)
			} else if check := err2.Code; check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})
	})

	t.Run("SetParam", func(t *testing.T) {
		t.Run("assign", func(t *testing.T) {
			param := 123
			code := 456
			msg := "message"
			expected := fmt.Sprintf("p:%d.c:%d", param, code)
			err1 := NewEnvelopeStatusError(code, msg)
			err2 := err1.SetParam(param)

			if check := err2.Param; check != param {
				t.Errorf("(%v) when expecting (%v)", check, param)
			} else if check := err2.Code; check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})
	})

	t.Run("SetError", func(t *testing.T) {
		t.Run("assign", func(t *testing.T) {
			newCode := 123
			code := 456
			expected := fmt.Sprintf("c:%d", newCode)
			msg := "message"
			err1 := NewEnvelopeStatusError(code, msg)
			err2 := err1.SetError(newCode)

			if check := err2.Error; check != fmt.Sprintf("%d", newCode) {
				t.Errorf("(%v) when expecting (%v)", check, newCode)
			} else if check := err2.Code; check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})

		t.Run("assign with string error code", func(t *testing.T) {
			code := 456
			newCode := "new error code"
			msg := "message"
			err1 := NewEnvelopeStatusError(code, msg)
			err2 := err1.SetError(newCode)

			if check := err2.Error; check != newCode {
				t.Errorf("(%v) when expecting (%v)", check, newCode)
			} else if check := err2.Code; check != newCode {
				t.Errorf("(%v) when expecting (%v)", check, newCode)
			}
		})
	})

	t.Run("SetMessage", func(t *testing.T) {
		t.Run("assign", func(t *testing.T) {
			code := 456
			msg := "message"
			newMsg := "new message"
			err1 := NewEnvelopeStatusError(code, msg)
			err2 := err1.SetMessage(newMsg)

			if check := err2.Message; check != newMsg {
				t.Errorf("(%v) when expecting (%v)", check, newMsg)
			}
		})
	})

	t.Run("GetCode", func(t *testing.T) {
		t.Run("retrieval", func(t *testing.T) {
			service := 12
			endpoint := 34
			param := 56
			code := 78
			expected := fmt.Sprintf("s:%d.e:%d.p:%d.c:%d", service, endpoint, param, code)
			msg := "message"
			e := NewEnvelopeStatusError(1, msg).
				SetService(service).
				SetEndpoint(endpoint).
				SetParam(param).
				SetError(code)

			if check := e.GetCode(); check != expected {
				t.Errorf("(%v) code instead of (%v)", check, expected)
			}
		})
	})

	t.Run("GetMessage", func(t *testing.T) {
		t.Run("retrieval", func(t *testing.T) {
			msg := "message"
			e := NewEnvelopeStatusError(123, msg)

			if check := e.GetMessage(); check != msg {
				t.Errorf("(%v) when expecting (%v)", check, msg)
			}
		})
	})

	t.Run("code generation", func(t *testing.T) {
		t.Run("assign param error", func(t *testing.T) {
			service := 12
			endpoint := 34
			param := 56
			code := 78
			expected := fmt.Sprintf("s:%d.e:%d.p:%d.c:%d", service, endpoint, param, code)
			originalCode := 1
			msg := "message"
			err1 := NewEnvelopeStatusError(originalCode, msg)
			err2 := err1.SetService(service).SetEndpoint(endpoint).SetParam(param).SetError(code)

			if check := err2.Service; check != service {
				t.Errorf("(%v) when expecting (%v)", check, service)
			} else if check := err2.Endpoint; check != endpoint {
				t.Errorf("(%v) when expecting (%v)", check, endpoint)
			} else if check := err2.Param; check != param {
				t.Errorf("(%v) when expecting (%v)", check, param)
			} else if check := err2.Error; check != fmt.Sprintf("%d", code) {
				t.Errorf("(%v) when expecting (%v)", check, code)
			} else if check := err2.Code; check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})

		t.Run("assign with string error code without param", func(t *testing.T) {
			service := 12
			endpoint := 34
			code := "error code"
			expected := fmt.Sprintf("s:%d.e:%d.%s", service, endpoint, code)
			originalCode := 1
			msg := "message"
			err1 := NewEnvelopeStatusError(originalCode, msg)
			err2 := err1.SetService(service).SetEndpoint(endpoint).SetError(code)

			if check := err2.Service; check != service {
				t.Errorf("(%v) when expecting (%v)", check, service)
			} else if check := err2.Endpoint; check != endpoint {
				t.Errorf("(%v) when expecting (%v)", check, endpoint)
			} else if check := err2.Param; check != 0 {
				t.Errorf("(%v) when expecting (%v)", check, 0)
			} else if check := err2.Error; check != code {
				t.Errorf("(%v) when expecting (%v)", check, code)
			} else if check := err2.Code; check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})

		t.Run("assign with string error code with param", func(t *testing.T) {
			service := 12
			endpoint := 34
			param := 56
			code := "error code"
			expected := fmt.Sprintf("s:%d.e:%d.p:%d.%s", service, endpoint, param, code)
			originalCode := 1
			msg := "message"
			err1 := NewEnvelopeStatusError(originalCode, msg)
			err2 := err1.SetService(service).SetEndpoint(endpoint).SetParam(param).SetError(code)

			if check := err2.Service; check != service {
				t.Errorf("(%v) when expecting (%v)", check, service)
			} else if check := err2.Endpoint; check != endpoint {
				t.Errorf("(%v) when expecting (%v)", check, endpoint)
			} else if check := err2.Param; check != param {
				t.Errorf("(%v) when expecting (%v)", check, param)
			} else if check := err2.Error; check != code {
				t.Errorf("(%v) when expecting (%v)", check, code)
			} else if check := err2.Code; check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})
	})
}

func Test_EnvelopeStatusErrorList(t *testing.T) {
	t.Run("MarshalXML", func(t *testing.T) {
		t.Run("empty list", func(t *testing.T) {
			name := "start"
			buffer := strings.Builder{}
			start := xml.StartElement{Name: xml.Name{Local: name}}
			list := EnvelopeStatusErrorList{}
			expected := "<start></start>"

			if err := list.MarshalXML(xml.NewEncoder(&buffer), start); err != nil {
				t.Errorf("returned the ujnexpected error (%v)", err)
			} else if check := buffer.String(); check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})

		t.Run("single element list", func(t *testing.T) {
			name := "start"
			buffer := strings.Builder{}
			start := xml.StartElement{Name: xml.Name{Local: name}}
			list := EnvelopeStatusErrorList{
				NewEnvelopeStatusError(1, "error message").
					SetService(2).
					SetEndpoint(3),
			}
			expected := `<start><error code="s:2.e:3.c:1" message="error message"></error></start>`

			if err := list.MarshalXML(xml.NewEncoder(&buffer), start); err != nil {
				t.Errorf("returned the ujnexpected error (%v)", err)
			} else if check := buffer.String(); check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})

		t.Run("multiple element list", func(t *testing.T) {
			name := "start"
			buffer := strings.Builder{}
			start := xml.StartElement{Name: xml.Name{Local: name}}
			list := EnvelopeStatusErrorList{
				NewEnvelopeStatusError(1, "error message 1").SetService(2).SetEndpoint(3),
				NewEnvelopeStatusError(2, "error message 2").SetService(2).SetEndpoint(3),
				NewEnvelopeStatusError(3, "error message 3").SetService(2).SetEndpoint(3),
			}
			expected := `<start>`
			expected += `<error code="s:2.e:3.c:1" message="error message 1"></error>`
			expected += `<error code="s:2.e:3.c:2" message="error message 2"></error>`
			expected += `<error code="s:2.e:3.c:3" message="error message 3"></error>`
			expected += `</start>`

			if err := list.MarshalXML(xml.NewEncoder(&buffer), start); err != nil {
				t.Errorf("returned the ujnexpected error (%v)", err)
			} else if check := buffer.String(); check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})
	})
}

func Test_EnvelopeStatus(t *testing.T) {
	t.Run("NewEnvelopeStatus", func(t *testing.T) {
		t.Run("construct", func(t *testing.T) {
			s := NewEnvelopeStatus()

			if check := s.Success; check != true {
				t.Error("initialized the status field as false")
			} else if len(s.Errors) != 0 {
				t.Error("initialized the error list with some elements")
			}
		})
	})

	t.Run("AddError", func(t *testing.T) {
		t.Run("add single error", func(t *testing.T) {
			err := NewEnvelopeStatusError(123, "error message")
			s := NewEnvelopeStatus().AddError(err)

			if s.Success != false {
				t.Error("didn't assign the status as false")
			} else if len(s.Errors) != 1 {
				t.Error("didn't stored the inserted error")
			} else if check := s.Errors[0]; !reflect.DeepEqual(err, check) {
				t.Errorf("(%v) when expecting (%v)", check, err)
			}
		})

		t.Run("add multiple error", func(t *testing.T) {
			err1 := NewEnvelopeStatusError(123, "error message 1")
			err2 := NewEnvelopeStatusError(456, "error message 2")
			err3 := NewEnvelopeStatusError(789, "error message 3")
			s := NewEnvelopeStatus().AddError(err1).AddError(err2).AddError(err3)

			if s.Success != false {
				t.Error("didn't assign the status as false")
			} else if len(s.Errors) != 3 {
				t.Error("didn't stored the inserted error")
			} else if check := s.Errors[0]; !reflect.DeepEqual(err1, check) {
				t.Errorf("(%v) when expecting (%v)", check, err1)
			} else if check := s.Errors[1]; !reflect.DeepEqual(err2, check) {
				t.Errorf("(%v) when expecting (%v)", check, err2)
			} else if check := s.Errors[2]; !reflect.DeepEqual(err3, check) {
				t.Errorf("(%v) when expecting (%v)", check, err3)
			}
		})
	})

	t.Run("SetService", func(t *testing.T) {
		t.Run("assign to all stored error", func(t *testing.T) {
			service := 147
			err1 := NewEnvelopeStatusError(123, "error message 1")
			err2 := NewEnvelopeStatusError(456, "error message 2")
			err3 := NewEnvelopeStatusError(789, "error message 3")
			s := NewEnvelopeStatus().AddError(err1).AddError(err2).AddError(err3)
			s = s.SetService(service)

			if check := s.Errors[0]; check.Service != service {
				t.Errorf("(%v) when expecting (%v)", check.Service, service)
			} else if check := s.Errors[1]; check.Service != service {
				t.Errorf("(%v) when expecting (%v)", check.Service, service)
			} else if check := s.Errors[2]; check.Service != service {
				t.Errorf("(%v) when expecting (%v)", check.Service, service)
			}
		})
	})

	t.Run("SetEndpoint", func(t *testing.T) {
		t.Run("assign to all stored error", func(t *testing.T) {
			endpoint := 147
			err1 := NewEnvelopeStatusError(123, "error message 1")
			err2 := NewEnvelopeStatusError(456, "error message 2")
			err3 := NewEnvelopeStatusError(789, "error message 3")
			s := NewEnvelopeStatus().AddError(err1).AddError(err2).AddError(err3)
			s = s.SetEndpoint(endpoint)

			if check := s.Errors[0]; check.Endpoint != endpoint {
				t.Errorf("(%v) when expecting (%v)", check.Endpoint, endpoint)
			} else if check := s.Errors[1]; check.Endpoint != endpoint {
				t.Errorf("(%v) when expecting (%v)", check.Endpoint, endpoint)
			} else if check := s.Errors[2]; check.Endpoint != endpoint {
				t.Errorf("(%v) when expecting (%v)", check.Endpoint, endpoint)
			}
		})
	})
}

func Test_EnvelopeListReport(t *testing.T) {
	t.Run("NewEnvelopeListReport", func(t *testing.T) {
		t.Run("store the search parameters", func(t *testing.T) {
			scenarios := []struct {
				search string
				start  uint
				count  uint
				total  uint
				prev   string
				next   string
			}{
				{ // report on start position
					search: "search string",
					start:  uint(0),
					count:  uint(2),
					total:  uint(10),
					prev:   "",
					next:   "?search=search string&start=2&count=2",
				},
				{ // report with truncated prev link
					search: "search string",
					start:  uint(1),
					count:  uint(2),
					total:  uint(10),
					prev:   "?search=search string&start=0&count=2",
					next:   "?search=search string&start=3&count=2",
				},
				{ // report with prev link
					search: "search string",
					start:  uint(2),
					count:  uint(2),
					total:  uint(10),
					prev:   "?search=search string&start=0&count=2",
					next:   "?search=search string&start=4&count=2",
				},
				{ // report with prev link (2)
					search: "search string",
					start:  uint(3),
					count:  uint(2),
					total:  uint(10),
					prev:   "?search=search string&start=1&count=2",
					next:   "?search=search string&start=5&count=2",
				},
				{ // report without next page
					search: "search string",
					start:  uint(8),
					count:  uint(2),
					total:  uint(10),
					prev:   "?search=search string&start=6&count=2",
					next:   "",
				},
				{ // report without next page (2)
					search: "search string",
					start:  uint(9),
					count:  uint(2),
					total:  uint(10),
					prev:   "?search=search string&start=7&count=2",
					next:   "",
				},
				{ // report without next page (3)
					search: "search string",
					start:  uint(10),
					count:  uint(2),
					total:  uint(10),
					prev:   "?search=search string&start=8&count=2",
					next:   "",
				},
			}

			for _, s := range scenarios {
				report := NewEnvelopeListReport(s.search, s.start, s.count, s.total)

				if check := report.Search; check != s.search {
					t.Errorf("(%v) when expecting (%v)", check, s.search)
				} else if check := report.Start; check != s.start {
					t.Errorf("(%v) when expecting (%v)", check, s.start)
				} else if check := report.Count; check != s.count {
					t.Errorf("(%v) when expecting (%v)", check, s.count)
				} else if check := report.Total; check != s.total {
					t.Errorf("(%v) when expecting (%v)", check, s.total)
				} else if check := report.Prev; check != s.prev {
					t.Errorf("(%v) when expecting (%v)", check, s.prev)
				} else if check := report.Next; check != s.next {
					t.Errorf("(%v) when expecting (%v)", check, s.next)
				}
			}
		})
	})
}

func Test_Envelope(t *testing.T) {
	t.Run("NewEnvelope", func(t *testing.T) {
		t.Run("construct without list report", func(t *testing.T) {
			statusCode := 123
			data := "message"
			env := NewEnvelope(statusCode, data)

			if check := env.Status.Success; check != true {
				t.Error("initialized the status field as false")
			} else if len(env.Status.Errors) != 0 {
				t.Error("initialized the error list with some elements")
			} else if check := env.Data; !reflect.DeepEqual(check, data) {
				t.Errorf("(%v) when expecting (%v)", check, data)
			} else if env.ListReport != nil {
				t.Errorf("unexpected (%v) list report", env.ListReport)
			}
		})

		t.Run("construct with nil list report", func(t *testing.T) {
			statusCode := 123
			data := "message"
			env := NewEnvelope(statusCode, data, nil)

			if check := env.Status.Success; check != true {
				t.Error("initialized the status field as false")
			} else if len(env.Status.Errors) != 0 {
				t.Error("initialized the error list with some elements")
			} else if check := env.Data; !reflect.DeepEqual(check, data) {
				t.Errorf("(%v) when expecting (%v)", check, data)
			} else if env.ListReport != nil {
				t.Errorf("unexpected (%v) list report", env.ListReport)
			}
		})

		t.Run("construct with list report", func(t *testing.T) {
			statusCode := 123
			data := "message"
			report := NewEnvelopeListReport("search", 1, 2, 10)
			env := NewEnvelope(statusCode, data, report)

			if check := env.Status.Success; check != true {
				t.Error("initialized the status field as false")
			} else if len(env.Status.Errors) != 0 {
				t.Error("initialized the error list with some elements")
			} else if check := env.Data; !reflect.DeepEqual(check, data) {
				t.Errorf("(%v) when expecting (%v)", check, data)
			} else if check := env.ListReport; !reflect.DeepEqual(check, report) {
				t.Errorf("(%v) when expecting (%v)", check, report)
			}
		})
	})

	t.Run("GetStatusCode", func(t *testing.T) {
		t.Run("return stored value", func(t *testing.T) {
			expected := 123
			if check := NewEnvelope(expected, nil).GetStatusCode(); check != expected {
				t.Errorf("(%v) when expecting (%v)", check, expected)
			}
		})
	})

	t.Run("AddError", func(t *testing.T) {
		statusCode := 123
		data := "message"

		t.Run("add single error", func(t *testing.T) {
			err := NewEnvelopeStatusError(123, "error message")
			env := NewEnvelope(statusCode, data).AddError(err)

			if env.Status.Success != false {
				t.Error("didn't assign the status as false")
			} else if len(env.Status.Errors) != 1 {
				t.Error("didn't stored the inserted error")
			} else if check := env.Status.Errors[0]; !reflect.DeepEqual(err, check) {
				t.Errorf("(%v) when expecting (%v)", check, err)
			}
		})

		t.Run("add multiple error", func(t *testing.T) {
			err1 := NewEnvelopeStatusError(123, "error message 1")
			err2 := NewEnvelopeStatusError(456, "error message 2")
			err3 := NewEnvelopeStatusError(789, "error message 3")
			env := NewEnvelope(statusCode, data).AddError(err1).AddError(err2).AddError(err3)

			if env.Status.Success != false {
				t.Error("didn't assign the status as false")
			} else if len(env.Status.Errors) != 3 {
				t.Error("didn't stored the inserted error")
			} else if check := env.Status.Errors[0]; !reflect.DeepEqual(err1, check) {
				t.Errorf("(%v) when expecting (%v)", check, err1)
			} else if check := env.Status.Errors[1]; !reflect.DeepEqual(err2, check) {
				t.Errorf("(%v) when expecting (%v)", check, err2)
			} else if check := env.Status.Errors[2]; !reflect.DeepEqual(err3, check) {
				t.Errorf("(%v) when expecting (%v)", check, err3)
			}
		})
	})

	t.Run("SetService", func(t *testing.T) {
		t.Run("assign to all stored error", func(t *testing.T) {
			service := 147
			statusCode := 123
			data := "message"
			err1 := NewEnvelopeStatusError(123, "error message 1")
			err2 := NewEnvelopeStatusError(456, "error message 2")
			err3 := NewEnvelopeStatusError(789, "error message 3")
			env := NewEnvelope(statusCode, data).AddError(err1).AddError(err2).AddError(err3)
			env = env.SetService(service)

			if check := env.Status.Errors[0]; check.Service != service {
				t.Errorf("(%v) when expecting (%v)", check.Service, service)
			} else if check := env.Status.Errors[1]; check.Service != service {
				t.Errorf("(%v) when expecting (%v)", check.Service, service)
			} else if check := env.Status.Errors[2]; check.Service != service {
				t.Errorf("(%v) when expecting (%v)", check.Service, service)
			}
		})
	})

	t.Run("SetEndpoint", func(t *testing.T) {
		t.Run("assign to all stored error", func(t *testing.T) {
			endpoint := 147
			statusCode := 123
			data := "message"
			err1 := NewEnvelopeStatusError(123, "error message 1")
			err2 := NewEnvelopeStatusError(456, "error message 2")
			err3 := NewEnvelopeStatusError(789, "error message 3")
			env := NewEnvelope(statusCode, data).AddError(err1).AddError(err2).AddError(err3)
			env = env.SetEndpoint(endpoint)

			if check := env.Status.Errors[0]; check.Endpoint != endpoint {
				t.Errorf("(%v) when expecting (%v)", check.Endpoint, endpoint)
			} else if check := env.Status.Errors[1]; check.Endpoint != endpoint {
				t.Errorf("(%v) when expecting (%v)", check.Endpoint, endpoint)
			} else if check := env.Status.Errors[2]; check.Endpoint != endpoint {
				t.Errorf("(%v) when expecting (%v)", check.Endpoint, endpoint)
			}
		})
	})

	t.Run("SetListReport", func(t *testing.T) {
		t.Run("assign the list report", func(t *testing.T) {
			endpoint := 147
			statusCode := 123
			data := "message"
			report := NewEnvelopeListReport("search", 1, 2, 10)
			env := NewEnvelope(statusCode, data)
			env = env.SetEndpoint(endpoint)
			env.SetListReport(report)

			if check := env.Status.Success; check != true {
				t.Error("initialized the status field as false")
			} else if len(env.Status.Errors) != 0 {
				t.Error("initialized the error list with some elements")
			} else if check := env.Data; !reflect.DeepEqual(check, data) {
				t.Errorf("(%v) when expecting (%v)", check, data)
			} else if check := env.ListReport; !reflect.DeepEqual(check, report) {
				t.Errorf("(%v) when expecting (%v)", check, report)
			}
		})

		t.Run("assign/remove the list report if nil given", func(t *testing.T) {
			endpoint := 147
			statusCode := 123
			data := "message"
			report := NewEnvelopeListReport("search", 1, 2, 10)
			env := NewEnvelope(statusCode, data, report)
			env = env.SetEndpoint(endpoint)
			env.SetListReport(nil)

			if check := env.Status.Success; check != true {
				t.Error("initialized the status field as false")
			} else if len(env.Status.Errors) != 0 {
				t.Error("initialized the error list with some elements")
			} else if check := env.Data; !reflect.DeepEqual(check, data) {
				t.Errorf("(%v) when expecting (%v)", check, data)
			} else if env.ListReport != nil {
				t.Errorf("unexpected (%v) list report", env.ListReport)
			}
		})
	})
}
