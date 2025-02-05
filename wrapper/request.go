package wrapper

import (
	"bytes"
	"context"
	"io"
	"maps"
	"net/url"
	"strings"
)

type Request interface {
	Context() context.Context
	Params() map[string]string
	Headers() map[string][]string
	Body() io.ReadCloser
	Method() string
	URL() *url.URL
	Query() url.Values
	Path() string
}

func Modifier(request Request) *RequestWrapper {
	return &RequestWrapper{
		ctx:     request.Context(),
		params:  request.Params(),
		headers: request.Headers(),
		body:    request.Body(),
		method:  request.Method(),
		url:     request.URL(),
		query:   request.Query(),
		path:    request.Path(),
	}
}

type RequestWrapper struct {
	ctx     context.Context
	method  string
	url     *url.URL
	query   url.Values
	path    string
	body    io.ReadCloser
	params  map[string]string
	headers url.Values
}

func (r *RequestWrapper) Context() context.Context { return r.ctx }
func (r *RequestWrapper) Method() string           { return r.method }
func (r *RequestWrapper) URL() *url.URL            { return r.url }
func (r *RequestWrapper) Query() url.Values        { return r.query }
func (r *RequestWrapper) Path() string             { return r.path }
func (r *RequestWrapper) Headers() url.Values      { return r.headers }

func (r *RequestWrapper) Body() io.ReadCloser {
	restore, ret := &bytes.Buffer{}, &bytes.Buffer{}

	wr := io.MultiWriter(restore, ret)
	tee := io.TeeReader(r.body, wr)

	_, err := io.ReadAll(tee)
	if err != nil {
		return io.NopCloser(bytes.NewReader([]byte{}))
	}

	r.body = io.NopCloser(restore)

	return io.NopCloser(ret)
}

func (r *RequestWrapper) SetBody(data []byte) {
	r.body = io.NopCloser(bytes.NewBuffer(data))
}

func (r *RequestWrapper) SetQueryParam(key, value string) {
	q, _ := url.ParseQuery(r.url.RawQuery)
	q.Set(key, value)
	r.url.RawQuery = q.Encode()
	r.query = r.url.Query()
}

func (r *RequestWrapper) Form() url.Values {
	var form url.Values

	ct, ok := r.headers["Content-Type"]
	if ok && strings.Join(ct, ";") == "application/x-www-form-urlencoded" {
		data, err := io.ReadAll(r.Body())
		if err == nil {
			form, _ = url.ParseQuery(string(data))
		}
	}

	result := make(url.Values, len(r.query)+len(form))

	maps.Copy(result, r.query)
	maps.Copy(result, form)

	return result
}

func (r *RequestWrapper) SetHeader(key, value string) {
	r.headers[key] = append(r.headers[key], value)
}
