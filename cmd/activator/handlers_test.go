/*
Copyright 2018 The Knative Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/knative/serving/cmd/util"
	"github.com/knative/serving/pkg/activator"
	"github.com/knative/serving/pkg/controller"
	"go.uber.org/zap"
)

type fakeActivator struct {
	endpoint  activator.Endpoint
	namespace string
	name      string
}

func newFakeActivator(namespace string, name string, server *httptest.Server) activator.Activator {
	url, _ := url.Parse(server.URL)
	host := url.Hostname()
	port, _ := strconv.Atoi(url.Port())

	return &fakeActivator{
		endpoint:  activator.Endpoint{FQDN: host, Port: int32(port)},
		namespace: namespace,
		name:      name,
	}
}

func (fa *fakeActivator) ActiveEndpoint(namespace, name string) (activator.Endpoint, activator.Status, error) {
	if namespace == fa.namespace && name == fa.name {
		return fa.endpoint, http.StatusOK, nil
	}

	return activator.Endpoint{}, http.StatusNotFound, errors.New("not found!")
}

func (fa *fakeActivator) Shutdown() {
}

func TestActivationHandler(t *testing.T) {
	logger := zap.NewExample().Sugar()

	errMsg := func(msg string) string {
		return fmt.Sprintf("Error getting active endpoint: %v\n", msg)
	}

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "everything good!")
		}),
	)
	defer server.Close()

	act := newFakeActivator("real-namespace", "real-name", server)

	examples := []struct {
		label     string
		namespace string
		name      string
		wantBody  string
		wantCode  int
		wantErr   error
	}{
		{
			label:     "active endpoint",
			namespace: "real-namespace",
			name:      "real-name",
			wantBody:  "everything good!",
			wantCode:  http.StatusOK,
			wantErr:   nil,
		},
		{
			label:     "no active endpoint",
			namespace: "fake-namespace",
			name:      "fake-name",
			wantBody:  errMsg("not found!"),
			wantCode:  http.StatusNotFound,
			wantErr:   nil,
		},
		{
			label:     "request error",
			namespace: "real-namespace",
			name:      "real-name",
			wantBody:  "",
			wantCode:  http.StatusBadGateway,
			wantErr:   errors.New("request error!"),
		},
	}

	for _, e := range examples {
		t.Run(e.label, func(t *testing.T) {
<<<<<<< HEAD
			rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				if r.Host != "" {
					t.Errorf("Unexpected request host. Want %q, got %q", "", r.Host)
				}

||||||| merged common ancestors
			rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
=======
			rt := util.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
>>>>>>> Move non-activator specific components from cmd/activator to cmd/util
				if e.wantErr != nil {
					return nil, e.wantErr
				}

				return http.DefaultTransport.RoundTrip(r)
			})

			handler := NewActivationHandler(act, rt, logger)

			resp := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "http://example.com", nil)
			req.Header.Set(controller.GetRevisionHeaderNamespace(), e.namespace)
			req.Header.Set(controller.GetRevisionHeaderName(), e.name)

			handler.ServeHTTP(resp, req)

			if resp.Code != e.wantCode {
				t.Errorf("Unexpected response status. Want %d, got %d", e.wantCode, resp.Code)
			}

			gotBody, _ := ioutil.ReadAll(resp.Body)
			if string(gotBody) != e.wantBody {
				t.Errorf("Unexpected response body. Want %q, got %q", e.wantBody, gotBody)
			}
		})
	}
}
