package sapi

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/happyhippyhippo/slate"
)

// ----------------------------------------------------------------------------
// defs
// ----------------------------------------------------------------------------

const (
	// ValidationContainerID defines the id to be used
	// as the container registration id of a validation.
	ValidationContainerID = slate.ContainerID + ".validation"

	// ValidationTranslatorContainerID defines the id to be used
	// as the container registration id of a translator.
	ValidationTranslatorContainerID = ValidationContainerID + ".translator"

	// ValidationUniversalTranslatorContainerID defines the id to be used
	// as the container registration id of a universal translator.
	ValidationUniversalTranslatorContainerID = ValidationContainerID + ".universal"

	// ValidationParserContainerID defines the id to be used
	// as the container registration id of an error parser instance.
	ValidationParserContainerID = ValidationContainerID + ".parser"

	// ValidationEnvID defines the validation module base environment variable name.
	ValidationEnvID = slate.EnvID + "_VALIDATION"
)

var (
	// ValidationLocale defines the default locale string to be used when
	// instantiating the translator.
	ValidationLocale = slate.EnvString(ValidationEnvID+"_LOCALE", "en")
)

// ----------------------------------------------------------------------------
// errors
// ----------------------------------------------------------------------------

var (
	// ErrValidationTranslatorNotFound defines an error that denotes
	// that a required error translator was not found.
	ErrValidationTranslatorNotFound = fmt.Errorf("validation translator not found")
)

func errValidationTranslatorNotFound(
	translator string,
	ctx ...map[string]interface{},
) error {
	return slate.NewErrorFrom(ErrValidationTranslatorNotFound, translator, ctx...)
}

// ----------------------------------------------------------------------------
// validation universal translator
// ----------------------------------------------------------------------------

// NewValidationUniversalTranslator @todo doc
func NewValidationUniversalTranslator() *ut.UniversalTranslator {
	lang := en.New()
	return ut.New(lang, lang)
}

// ----------------------------------------------------------------------------
// validation translator
// ----------------------------------------------------------------------------

// NewValidationTranslator @todo doc
func NewValidationTranslator(
	universalTranslator *ut.UniversalTranslator,
) (ut.Translator, error) {
	translator, found := universalTranslator.GetTranslator(ValidationLocale)
	if found == false {
		return nil, errValidationTranslatorNotFound(ValidationLocale)
	}
	return translator, nil
}

// ----------------------------------------------------------------------------
// validation parser
// ----------------------------------------------------------------------------

// ValidationParser @todo doc
type ValidationParser struct {
	mapper     map[string]int
	translator ut.Translator
}

// NewValidationParser instantiate a new validation parser instance
func NewValidationParser(
	translator ut.Translator,
) (*ValidationParser, error) {
	if translator == nil {
		return nil, errNilPointer("translator")
	}

	return &ValidationParser{
		mapper: map[string]int{
			"eqcsfield":     1,
			"eqfield":       2,
			"fieldcontains": 3,
			"fieldexcludes": 4,
			"gtcsfield":     5,
			"gtecsfield":    6,
			"gtefield":      7,
			"gtfield":       8,
			"ltcsfield":     9,
			"ltecsfield":    10,
			"ltefield":      11,
			"ltfield":       12,
			"necsfield":     13,
			"nefield":       14,

			"cidr":             15,
			"cidrv4":           16,
			"cidrv6":           17,
			"datauri":          18,
			"fqdn":             19,
			"hostname":         20,
			"hostname_port":    21,
			"hostname_rfc1123": 22,
			"ip":               23,
			"ip4_addr":         24,
			"ip6_addr":         25,
			"ip_addr":          26,
			"ipv4":             27,
			"ipv6":             28,
			"mac":              29,
			"tcp4_addr":        30,
			"tcp6_addr":        31,
			"tcp_addr":         32,
			"udp4_addr":        33,
			"udp6_addr":        34,
			"udp_addr":         35,
			"unix_addr":        36,
			"uri":              37,
			"url":              38,
			"url_encoded":      39,
			"urn_rfc2141":      40,

			"alpha":           41,
			"alphanum":        42,
			"alphanumunicode": 43,
			"alphaunicode":    44,
			"ascii":           45,
			"contains":        46,
			"containsany":     47,
			"containsrune":    48,
			"endswith":        49,
			"lowercase":       50,
			"multibyte":       51,
			"number":          52,
			"numeric":         53,
			"printascii":      54,
			"startswith":      55,
			"uppercase":       56,

			"base64":          57,
			"base64url":       58,
			"btc_addr":        59,
			"btc_addr_bech32": 60,
			"datetime":        61,
			"e164":            62,
			"email":           63,
			"eth_addr":        64,
			"hexadecimal":     65,
			"hexcolor":        66,
			"hsl":             67,
			"hsla":            68,
			"html":            69,
			"html_encoded":    70,
			"isbn":            71,
			"isbn10":          72,
			"isbn13":          73,
			"json":            74,
			"latitude":        75,
			"longitude":       76,
			"rgb":             77,
			"rgba":            78,
			"ssn":             79,
			"uuid":            80,
			"uuid3":           81,
			"uuid3_rfc4122":   82,
			"uuid4":           83,
			"uuid4_rfc4122":   84,
			"uuid5":           85,
			"uuid5_rfc4122":   86,
			"uuid_rfc4122":    87,

			"eq":  88,
			"gt":  89,
			"gte": 90,
			"lt":  91,
			"lte": 92,
			"ne":  93,

			"dir":                  94,
			"excludes":             95,
			"excludesall":          96,
			"excludesrune":         97,
			"file":                 98,
			"isdefault":            99,
			"len":                  100,
			"max":                  101,
			"min":                  102,
			"oneof":                103,
			"required":             104,
			"required_if":          105,
			"required_unless":      106,
			"required_with":        107,
			"required_with_all":    108,
			"required_without":     109,
			"required_without_all": 110,
			"excluded_with":        111,
			"excluded_with_all":    112,
			"excluded_without":     113,
			"excluded_without_all": 114,
			"unique":               115,
		},
		translator: translator,
	}, nil
}

// Parse method that will convert the list of validation error into
// an envelope struct to be used as the endpoint response.
func (p *ValidationParser) Parse(
	val interface{},
	errs validator.ValidationErrors,
) (*Envelope, error) {
	if val == nil {
		return nil, errNilPointer("value")
	}

	if errs == nil || len(errs) == 0 {
		return nil, nil
	}

	resp := NewEnvelope(http.StatusBadRequest, nil, nil)
	for _, err := range errs {
		parsed, e := p.convert(val, err)
		if e != nil {
			return nil, e
		}
		resp = resp.AddError(parsed)
	}

	return resp, nil
}

// AddError will add a validation mapped error to code value.
func (p *ValidationParser) AddError(
	e string,
	code int,
) {
	p.mapper[e] = code
}

func (p *ValidationParser) convert(
	value interface{},
	e validator.FieldError,
) (*EnvelopeStatusError, error) {
	if e == nil {
		return nil, errNilPointer("error")
	}

	typeof := reflect.TypeOf(value)
	field, _ := typeof.FieldByName(e.StructField())
	iparam := 0
	if param, ok := field.Tag.Lookup("vparam"); ok {
		var err error
		if iparam, err = strconv.Atoi(param); err != nil {
			return nil, err
		}
	}

	return NewEnvelopeStatusError(p.mapper[e.Tag()], e.Translate(p.translator)).SetParam(iparam), nil
}

// ----------------------------------------------------------------------------
// validator
// ----------------------------------------------------------------------------

// Validator is a function type used to define a calling interface of
// function responsible to validate an instance of a structure and return
// an initialized response envelope with the founded error
type Validator func(val interface{}) (*Envelope, error)

// NewValidator instantiates a new validation function
func NewValidator(
	translator ut.Translator,
	parser *ValidationParser,
) (Validator, error) {
	// check validate argument reference
	if translator == nil {
		return nil, errNilPointer("translator")
	}
	// check parser argument reference
	if parser == nil {
		return nil, errNilPointer("parser")
	}
	// register the translator in the used validator
	validate := validator.New()
	_ = translations.RegisterDefaultTranslations(validate, translator)
	// return the validation method instance
	return func(value interface{}) (*Envelope, error) {
		// check the value argument reference
		if value == nil {
			return nil, errNilPointer("value")
		}
		// validate the given structure
		if errs := validate.Struct(value); errs != nil {
			// compose the response envelope with the parsed validation error
			return parser.Parse(value, errs.(validator.ValidationErrors))
		}
		return nil, nil
	}, nil
}

// ----------------------------------------------------------------------------
// validation service register
// ----------------------------------------------------------------------------

// ValidationServiceRegister @todo doc
type ValidationServiceRegister struct {
	slate.ServiceRegister
}

var _ slate.ServiceProvider = &ValidationServiceRegister{}

// NewValidationServiceRegister will generate a new registry instance
func NewValidationServiceRegister(
	app ...*slate.App,
) *ValidationServiceRegister {
	return &ValidationServiceRegister{
		ServiceRegister: *slate.NewServiceRegister(app...),
	}
}

// Provide will register the validation package instances in the
// application container
func (sr ValidationServiceRegister) Provide(
	container *slate.ServiceContainer,
) error {
	// check container argument reference
	if container == nil {
		return errNilPointer("container")
	}
	_ = container.Add(ValidationUniversalTranslatorContainerID, NewValidationUniversalTranslator)
	_ = container.Add(ValidationTranslatorContainerID, NewValidationTranslator)
	_ = container.Add(ValidationParserContainerID, NewValidationParser)
	_ = container.Add(ValidationContainerID, NewValidator)
	return nil
}
