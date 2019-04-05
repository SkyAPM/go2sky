package propagation

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const splitToken string = "-"
const idToken string = "."

var (
	errEmptyHeader = errors.New("empty header")
	errInsufficientHeaderEntities = errors.New("insufficient header entities")
)

// DownstreamContext define the trace context from downstream
type DownstreamContext interface {
	Header() string
}

// UpstreamContext define the trace context to upstream
type UpstreamContext interface {
	SetHeader(header string)
}

// Extractor is a tool specification which define how to
// extract trace parent context from propagation context
type Extractor func() (DownstreamContext, error)

// Injector is a tool specification which define how to
// inject trace context into propagation context
type Injector func(carrier UpstreamContext) error

// Refs defines propagation specification of SkyWalking
type SpanContext struct {
	Sample                  int8
	TraceID                 []int64
	ParentSegmentID         []int64
	ParentSpanID            int32
	ParentServiceInstanceID int32
	EntryServiceInstanceID  int32
	NetworkAddress          string
	NetworkAddressID        int32
	EntryEndpoint           string
	EntryEndpointID         int32
	ParentEndpoint          string
	ParentEndpointID        int32
}

// DecodeSW6 converts string header to Refs
func (tc *SpanContext) DecodeSW6(header string) error {
	if header == "" {
		return errEmptyHeader
	}
	hh := strings.Split(header, splitToken)
	if len(hh) < 7 {
		return errors.WithMessagef(errInsufficientHeaderEntities, "header string: %s", header)
	}
	sample, err := strconv.ParseInt(hh[0], 10, 8)
	if err != nil {
		return errors.Errorf("str to int8 error %s", hh[0])
	}
	tc.Sample = int8(sample)
	tc.TraceID, err = stringConvertGlobalID(hh[1])
	if err != nil {
		return errors.Wrap(err, "trace id parse error")
	}
	tc.ParentSegmentID, err = stringConvertGlobalID(hh[2])
	if err != nil {
		return errors.Wrap(err, "parent segment id parse error")
	}
	tc.ParentSpanID, err = stringConvertInt32(hh[3])
	if err != nil {
		return errors.Wrap(err, "parent span id parse error")
	}
	tc.ParentServiceInstanceID, err = stringConvertInt32(hh[4])
	if err != nil {
		return errors.Wrap(err, "parent service instance id parse error")
	}
	tc.EntryServiceInstanceID, err = stringConvertInt32(hh[5])
	if err != nil {
		return errors.Wrap(err, "entry service instance id parse error")
	}
	tc.NetworkAddress, tc.NetworkAddressID, err = decodeBase64(hh[6])
	if err != nil {
		return errors.Wrap(err, "network address parse error")
	}
	if len(hh) < 9 {
		return nil
	}
	tc.EntryEndpoint, tc.EntryEndpointID, err = decodeBase64(hh[7])
	if err != nil {
		return errors.Wrap(err, "entry endpoint parse error")
	}
	tc.ParentEndpoint, tc.ParentEndpointID, err = decodeBase64(hh[8])
	if err != nil {
		return errors.Wrap(err, "parent endpoint parse error")
	}
	return nil
}

func stringConvertGlobalID(str string) ([]int64, error) {
	idStr, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, errors.Wrapf(err, "decode id error %s", str)
	}
	ss := strings.Split(string(idStr), idToken)
	if len(ss) < 3 {
		return nil, errors.Errorf("decode id entities error %s", string(idStr))
	}
	ii := make([]int64, 0, len(ss))
	for i, s := range ss {
		ii[i], err = strconv.ParseInt(s, 0, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "convert id error convert id entities to int32 error %s", s)
		}
	}
	return ii, nil
}

func stringConvertInt32(str string) (int32, error) {
	i, err := strconv.ParseInt(str, 0, 32)
	return int32(i), err
}

func decodeBase64(str string) (string, int32, error) {
	ret, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", 0, err
	}
	retStr := string(ret)
	if strings.HasPrefix(retStr, "#") {
		return retStr[1:], 0, nil
	}
	i, err := strconv.ParseInt(retStr, 0, 32)
	if err != nil {
		return "", 0, err
	}
	return "", int32(i), nil
}
