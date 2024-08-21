package jsonio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type jsonReader struct {
	checkSyntaxError           bool
	checkUnexpectedEOF         bool
	checkUnmarshalTypeError    bool
	checkEOF                   bool
	checkUnknownField          bool
	checkMaxBytesError         bool
	checkInvalidUnmarshalError bool
	disallowUnknownFields      bool
	maxBodySize                int64
}

// option represents a functional option for configuring the JSONReader.
type option func(*jsonReader)

// NewJSONReader creates a new JSONReader with the provided options.
func NewJSONReader(options ...option) *jsonReader {
	jr := &jsonReader{
		checkSyntaxError:           true,
		checkUnexpectedEOF:         true,
		checkUnmarshalTypeError:    true,
		checkEOF:                   true,
		checkUnknownField:          true,
		checkMaxBytesError:         true,
		checkInvalidUnmarshalError: true,
		disallowUnknownFields:      true,
		maxBodySize:                1_048_576, // Default max size: 1MB
	}
	for _, opt := range options {
		opt(jr)
	}
	return jr
}

// ReadJSON reads and decodes JSON from the request body into dst.
func (jr *jsonReader) ReadJSON(r *http.Request, dst any) error {

	if jr.maxBodySize > 0 {
		r.Body = http.MaxBytesReader(nil, r.Body, jr.maxBodySize)
	}
	dec := json.NewDecoder(r.Body)
	if jr.disallowUnknownFields {
		dec.DisallowUnknownFields()
	}
	err := dec.Decode(dst)
	// check body has one json value
	err2 := dec.Decode(&struct{}{})
	if !errors.Is(err2, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return jr.addContextToError(err)
}

func (jr *jsonReader) addContextToError(err error) error {
	if err == nil {
		return nil
	}
	checks := []struct {
		enabled bool
		fn      func(error) (bool, error)
	}{
		{jr.checkSyntaxError, jr.handleSyntaxError},
		{jr.checkUnexpectedEOF, jr.handleUnexpectedEOF},
		{jr.checkUnmarshalTypeError, jr.handleUnmarshalTypeError},
		{jr.checkEOF, jr.handleEOF},
		{jr.checkUnknownField, jr.handleUnknownField},
		{jr.checkMaxBytesError, jr.handleMaxBytesError},
		{jr.checkInvalidUnmarshalError, jr.handleInvalidUnmarshalError},
	}

	for _, check := range checks {
		if check.enabled {
			if hasTransformed, err := check.fn(err); hasTransformed {
				return err
			}
		}
	}
	return err
}

// SetMaxBodySize sets the maximum allowed size for the request body.
func SetMaxBodySize(size int64) option {
	return func(jr *jsonReader) {
		jr.maxBodySize = size
	}
}

// DisableSyntaxErrorCheck disables the syntax error check.
func DisableSyntaxErrorCheck() option {
	return func(jr *jsonReader) {
		jr.checkSyntaxError = false
	}
}

// DisableUnexpectedEOFCheck disables the unexpected EOF check.
func DisableUnexpectedEOFCheck() option {
	return func(jr *jsonReader) {
		jr.checkUnexpectedEOF = false
	}
}

// DisableUnmarshalTypeErrorCheck disables the unmarshal type error check.
func DisableUnmarshalTypeErrorCheck() option {
	return func(jr *jsonReader) {
		jr.checkUnmarshalTypeError = false
	}
}

// DisableEOFCheck disables the EOF check.
func DisableEOFCheck() option {
	return func(jr *jsonReader) {
		jr.checkEOF = false
	}
}

// DisableUnknownFieldCheck disables the unknown field check.
func DisableUnknownFieldCheck() option {
	return func(jr *jsonReader) {
		jr.checkUnknownField = false
	}
}

// DisableMaxBytesErrorCheck disables the max bytes error check.
func DisableMaxBytesErrorCheck() option {
	return func(jr *jsonReader) {
		jr.checkMaxBytesError = false
	}
}

// DisableInvalidUnmarshalErrorCheck disables the invalid unmarshal error check.
func DisableInvalidUnmarshalErrorCheck() option {
	return func(jr *jsonReader) {
		jr.checkInvalidUnmarshalError = false
	}
}

// DisableDisallowUnknownFields disables the DisallowUnknownFields check.
func DisableDisallowUnknownFields() option {
	return func(jr *jsonReader) {
		jr.disallowUnknownFields = false
	}
}

// the JSON can't be parsed.
func (jr *jsonReader) handleSyntaxError(err error) (bool, error) {
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return true, fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
	}
	return false, nil
}

func (jr *jsonReader) handleUnexpectedEOF(err error) (bool, error) {
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return true, errors.New("body contains badly-formed JSON")
	}
	return false, err
}

func (jr *jsonReader) handleUnmarshalTypeError(err error) (bool, error) {
	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		if unmarshalTypeError.Field != "" {
			return true, fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return true, fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
	}
	return false, err
}

func (jr *jsonReader) handleEOF(err error) (bool, error) {
	if errors.Is(err, io.EOF) {
		return true, errors.New("body must not be empty")
	}
	return false, err
}

func (jr *jsonReader) handleUnknownField(err error) (bool, error) {
	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return true, fmt.Errorf("body contains unknown key %s", fieldName)
	}
	return false, err
}

func (jr *jsonReader) handleMaxBytesError(err error) (bool, error) {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		return true, fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
	}
	return false, err
}

func (jr *jsonReader) handleInvalidUnmarshalError(err error) (bool, error) {
	var invalidUnmarshalError *json.InvalidUnmarshalError
	if errors.As(err, &invalidUnmarshalError) {
		panic(err)
	}
	return false, err
}
