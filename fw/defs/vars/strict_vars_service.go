// Code generated by clubbygen.
// GENERATED FILE DO NOT EDIT
// +build clubby_strict

package vars

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cesanta.com/common/go/mgrpc"
	"cesanta.com/common/go/mgrpc/frame"
	"cesanta.com/common/go/ourjson"
	"cesanta.com/common/go/ourtrace"
	"github.com/cesanta/errors"
	"golang.org/x/net/trace"

	"github.com/cesanta/ucl"
	"github.com/cesanta/validate-json/schema"
	"github.com/golang/glog"
)

var _ = bytes.MinRead
var _ = fmt.Errorf
var emptyMessage = ourjson.RawMessage{}
var _ = ourtrace.New
var _ = trace.New

const ServiceID = "http://mongoose-iot.com/fwVars"

type GetResult struct {
	Arch         *string `json:"arch,omitempty"`
	Fw_id        *string `json:"fw_id,omitempty"`
	Fw_timestamp *string `json:"fw_timestamp,omitempty"`
	Fw_version   *string `json:"fw_version,omitempty"`
	Mac_address  *string `json:"mac_address,omitempty"`
}

type Service interface {
	Get(ctx context.Context) (*GetResult, error)
}

type Instance interface {
	Call(context.Context, string, *frame.Command) (*frame.Response, error)
	//TraceCall(context.Context, string, *frame.Command) (context.Context, trace.Trace, func(*error))
}

type _validators struct {
	// This comment prevents gofmt from aligning types in the struct.
	GetResult *schema.Validator
}

var (
	validators     *_validators
	validatorsOnce sync.Once
)

func initValidators() {
	validators = &_validators{}

	loader := schema.NewLoader()

	service, err := ucl.Parse(bytes.NewBuffer(_ServiceDefinition))
	if err != nil {
		panic(err)
	}
	// Patch up shortcuts to be proper schemas.
	for _, v := range service.(*ucl.Object).Find("methods").(*ucl.Object).Value {
		if s, ok := v.(*ucl.Object).Find("result").(*ucl.String); ok {
			for kk := range v.(*ucl.Object).Value {
				if kk.Value == "result" {
					v.(*ucl.Object).Value[kk] = &ucl.Object{
						Value: map[ucl.Key]ucl.Value{
							ucl.Key{Value: "type"}: s,
						},
					}
				}
			}
		}
		if v.(*ucl.Object).Find("args") == nil {
			continue
		}
		args := v.(*ucl.Object).Find("args").(*ucl.Object)
		for kk, vv := range args.Value {
			if s, ok := vv.(*ucl.String); ok {
				args.Value[kk] = &ucl.Object{
					Value: map[ucl.Key]ucl.Value{
						ucl.Key{Value: "type"}: s,
					},
				}
			}
		}
	}
	var s *ucl.Object
	_ = s // avoid unused var error
	validators.GetResult, err = schema.NewValidator(service.(*ucl.Object).Find("methods").(*ucl.Object).Find("Get").(*ucl.Object).Find("result"), loader)
	if err != nil {
		panic(err)
	}
}

func NewClient(i Instance, addr string) Service {
	validatorsOnce.Do(initValidators)
	return &_Client{i: i, addr: addr}
}

type _Client struct {
	i    Instance
	addr string
}

func (c *_Client) Get(ctx context.Context) (res *GetResult, err error) {
	cmd := &frame.Command{
		Cmd: "Vars.Get",
	}
	//ctx, tr, finish := c.i.TraceCall(pctx, c.addr, cmd)
	//defer finish(&err)
	//_ = tr
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.Status != 0 {
		return nil, errors.Trace(&mgrpc.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}

	//tr.LazyPrintf("res: %s", ourjson.LazyJSON(&resp))

	bb, err := resp.Response.MarshalJSON()
	if err != nil {
		glog.Errorf("Failed to marshal result as JSON: %+v", err)
	} else {
		rv, err := ucl.Parse(bytes.NewReader(bb))
		if err == nil {
			if err := validators.GetResult.Validate(rv); err != nil {
				glog.Warningf("Got invalid result for Get: %+v", err)
				return nil, errors.Annotatef(err, "invalid response for Get")
			}
		}
	}
	var r *GetResult
	err = resp.Response.UnmarshalInto(&r)
	if err != nil {
		return nil, errors.Annotatef(err, "unmarshaling response")
	}
	return r, nil
}

//func RegisterService(i *clubby.Instance, impl Service) error {
//validatorsOnce.Do(initValidators)
//s := &_Server{impl}
//i.RegisterCommandHandler("Vars.Get", s.Get)
//i.RegisterService(ServiceID, _ServiceDefinition)
//return nil
//}

type _Server struct {
	impl Service
}

func (s *_Server) Get(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	r, err := s.impl.Get(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	bb, err := json.Marshal(r)
	if err == nil {
		v, err := ucl.Parse(bytes.NewBuffer(bb))
		if err != nil {
			glog.Errorf("Failed to parse just serialized JSON value %q: %+v", string(bb), err)
		} else {
			if err := validators.GetResult.Validate(v); err != nil {
				glog.Warningf("Returned invalid response for Get: %+v", err)
				return nil, errors.Annotatef(err, "server generated invalid responce for Get")
			}
		}
	}
	return r, nil
}

var _ServiceDefinition = json.RawMessage([]byte(`{
  "methods": {
    "Get": {
      "doc": "Get device read-only vars",
      "result": {
        "properties": {
          "arch": {
            "type": "string"
          },
          "fw_id": {
            "type": "string"
          },
          "fw_timestamp": {
            "type": "string"
          },
          "fw_version": {
            "type": "string"
          },
          "mac_address": {
            "type": "string"
          }
        },
        "type": "object"
      }
    }
  },
  "name": "Vars",
  "namespace": "http://mongoose-iot.com/fw"
}`))
