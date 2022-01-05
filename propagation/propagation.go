//
// Copyright 2022 SkyAPM org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

/*
Package propagation holds the required function signatures for Injection and
Extraction. It also contains decoder and encoder of SkyWalking propagation protocol.
*/
package propagation

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	Header                        string = "sw8"
	HeaderCorrelation             string = "sw8-correlation"
	headerLen                     int    = 8
	splitToken                    string = "-"
	correlationSplitToken         string = ","
	correlationKeyValueSplitToken string = ":"
)

var (
	errEmptyHeader                = errors.New("empty header")
	errInsufficientHeaderEntities = errors.New("insufficient header entities")
)

// Extractor is a tool specification which define how to
// extract trace parent context from propagation context
type Extractor func(headerKey string) (string, error)

// Injector is a tool specification which define how to
// inject trace context into propagation context
type Injector func(headerKey, headerValue string) error

// SpanContext defines propagation specification of SkyWalking
type SpanContext struct {
	TraceID               string            `json:"trace_id"`
	ParentSegmentID       string            `json:"parent_segment_id"`
	ParentService         string            `json:"parent_service"`
	ParentServiceInstance string            `json:"parent_service_instance"`
	ParentEndpoint        string            `json:"parent_endpoint"`
	AddressUsedAtClient   string            `json:"address_used_at_client"`
	ParentSpanID          int32             `json:"parent_span_id"`
	Sample                int8              `json:"sample"`
	Valid                 bool              `json:"valid"`
	CorrelationContext    map[string]string `json:"correlation_context"`
}

// Decode all SpanContext data from Extractor
func (tc *SpanContext) Decode(extractor Extractor) error {
	tc.Valid = false
	// sw8
	err := tc.decode(extractor, Header, tc.DecodeSW8)
	if err != nil {
		return err
	}

	// correlation
	err = tc.decode(extractor, HeaderCorrelation, tc.DecodeSW8Correlation)
	if err != nil {
		return err
	}
	return nil
}

// Encode all SpanContext data to Injector
func (tc *SpanContext) Encode(injector Injector) error {
	// sw8
	err := injector(Header, tc.EncodeSW8())
	if err != nil {
		return err
	}
	// correlation
	err = injector(HeaderCorrelation, tc.EncodeSW8Correlation())
	if err != nil {
		return err
	}
	return nil
}

// DecodeSW6 converts string header to SpanContext
func (tc *SpanContext) DecodeSW8(header string) error {
	if header == "" {
		return errEmptyHeader
	}
	hh := strings.Split(header, splitToken)
	if len(hh) < headerLen {
		return errors.WithMessagef(errInsufficientHeaderEntities, "header string: %s", header)
	}
	sample, err := strconv.ParseInt(hh[0], 10, 8)
	if err != nil {
		return errors.Errorf("str to int8 error %s", hh[0])
	}
	tc.Sample = int8(sample)
	tc.TraceID, err = decodeBase64(hh[1])
	if err != nil {
		return errors.Wrap(err, "trace id parse error")
	}
	tc.ParentSegmentID, err = decodeBase64(hh[2])
	if err != nil {
		return errors.Wrap(err, "parent segment id parse error")
	}
	tc.ParentSpanID, err = stringConvertInt32(hh[3])
	if err != nil {
		return errors.Wrap(err, "parent span id parse error")
	}
	tc.ParentService, err = decodeBase64(hh[4])
	if err != nil {
		return errors.Wrap(err, "parent service parse error")
	}
	tc.ParentServiceInstance, err = decodeBase64(hh[5])
	if err != nil {
		return errors.Wrap(err, "parent service instance parse error")
	}
	tc.ParentEndpoint, err = decodeBase64(hh[6])
	if err != nil {
		return errors.Wrap(err, "parent endpoint parse error")
	}
	tc.AddressUsedAtClient, err = decodeBase64(hh[7])
	if err != nil {
		return errors.Wrap(err, "network address parse error")
	}
	tc.Valid = true
	return nil
}

// EncodeSW6 converts SpanContext to string header
func (tc *SpanContext) EncodeSW8() string {
	return strings.Join([]string{
		fmt.Sprint(tc.Sample),
		encodeBase64(tc.TraceID),
		encodeBase64(tc.ParentSegmentID),
		fmt.Sprint(tc.ParentSpanID),
		encodeBase64(tc.ParentService),
		encodeBase64(tc.ParentServiceInstance),
		encodeBase64(tc.ParentEndpoint),
		encodeBase64(tc.AddressUsedAtClient),
	}, "-")
}

// DecodeSW8Correlation converts correlation string header to SpanContext
func (tc *SpanContext) DecodeSW8Correlation(header string) error {
	tc.CorrelationContext = make(map[string]string)
	if header == "" {
		return nil
	}

	hh := strings.Split(header, correlationSplitToken)
	for inx := range hh {
		keyValues := strings.Split(hh[inx], correlationKeyValueSplitToken)
		if len(keyValues) != 2 {
			continue
		}
		decodedKey, err := decodeBase64(keyValues[0])
		if err != nil {
			continue
		}
		decodedValue, err := decodeBase64(keyValues[1])
		if err != nil {
			continue
		}

		tc.CorrelationContext[decodedKey] = decodedValue
	}
	return nil
}

// EncodeSW8Correlation converts correlation to string header
func (tc *SpanContext) EncodeSW8Correlation() string {
	if len(tc.CorrelationContext) == 0 {
		return ""
	}

	content := make([]string, 0, len(tc.CorrelationContext))
	for k, v := range tc.CorrelationContext {
		content = append(content, fmt.Sprintf("%s%s%s", encodeBase64(k), correlationKeyValueSplitToken, encodeBase64(v)))
	}
	return strings.Join(content, correlationSplitToken)
}

func stringConvertInt32(str string) (int32, error) {
	i, err := strconv.ParseInt(str, 0, 32)
	return int32(i), err
}

func decodeBase64(str string) (string, error) {
	ret, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func encodeBase64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func (tc *SpanContext) decode(extractor Extractor, headerKey string, decoder func(header string) error) error {
	val, err := extractor(headerKey)
	if err != nil {
		return err
	}
	if val == "" {
		return nil
	}
	err = decoder(val)
	if err != nil {
		return err
	}
	return nil
}
