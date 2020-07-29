package endpoint_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/influxdata/influxdb/v2"
	"github.com/influxdata/influxdb/v2/mock"
	"github.com/influxdata/influxdb/v2/notification/endpoint"
	influxTesting "github.com/influxdata/influxdb/v2/testing"
)

const (
	id1 = "020f755c3c082000"
	id3 = "020f755c3c082002"
)

var goodBase = endpoint.Base{
	ID:          influxTesting.MustIDBase16Ptr(id1),
	Name:        "name1",
	OrgID:       influxTesting.MustIDBase16Ptr(id3),
	Status:      influxdb.Active,
	Description: "desc1",
}

func TestValidEndpoint(t *testing.T) {
	cases := []struct {
		name string
		src  influxdb.NotificationEndpoint
		err  error
	}{
		{
			name: "invalid endpoint id",
			src:  &endpoint.Slack{},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "Notification Endpoint ID is invalid",
			},
		},
		{
			name: "invalid status",
			src: &endpoint.PagerDuty{
				Base: endpoint.Base{
					ID:    influxTesting.MustIDBase16Ptr(id1),
					Name:  "name1",
					OrgID: influxTesting.MustIDBase16Ptr(id3),
				},
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "invalid status",
			},
		},
		{
			name: "empty name PagerDuty",
			src: &endpoint.PagerDuty{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
				},
				ClientURL:  "https://events.pagerduty.com/v2/enqueue",
				RoutingKey: influxdb.SecretField{Key: id1 + "-routing-key"},
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "Notification Endpoint Name can't be empty",
			},
		},
		{
			name: "empty name Telegram",
			src: &endpoint.Telegram{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
				},
				Token:   influxdb.SecretField{Key: id1 + "-token"},
				Channel: "-1001406363649",
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "Notification Endpoint Name can't be empty",
			},
		},
		{
			name: "empty slack url",
			src: &endpoint.Slack{
				Base: goodBase,
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "slack endpoint URL must be provided",
			},
		},
		{
			name: "invalid slack url",
			src: &endpoint.Slack{
				Base: goodBase,
				URL:  "posts://er:{DEf1=ghi@:5432/db?ssl",
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "slack endpoint URL is invalid: parse posts://er:{DEf1=ghi@:5432/db?ssl: net/url: invalid userinfo",
			},
		},
		{
			name: "empty slack token",
			src: &endpoint.Slack{
				Base: goodBase,
				URL:  "localhost",
			},
			err: nil,
		},
		{
			name: "empty http http method",
			src: &endpoint.HTTP{
				Base: goodBase,
				URL:  "localhost",
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "invalid http http method",
			},
		},
		{
			name: "empty http token",
			src: &endpoint.HTTP{
				Base:       goodBase,
				URL:        "localhost",
				Method:     "GET",
				AuthMethod: "bearer",
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "invalid http token for bearer auth",
			},
		},
		{
			name: "empty http username",
			src: &endpoint.HTTP{
				Base:       goodBase,
				URL:        "localhost",
				Method:     http.MethodGet,
				AuthMethod: "basic",
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "invalid http username/password for basic auth",
			},
		},
		{
			name: "empty telegram token",
			src: &endpoint.Telegram{
				Base: goodBase,
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "empty telegram bot token",
			},
		},
		{
			name: "empty telegram channel",
			src: &endpoint.Telegram{
				Base:  goodBase,
				Token: influxdb.SecretField{Key: id1 + "-token"},
			},
			err: &influxdb.Error{
				Code: influxdb.EInvalid,
				Msg:  "empty telegram channel",
			},
		},
		{
			name: "valid telegram token",
			src: &endpoint.Telegram{
				Base:    goodBase,
				Token:   influxdb.SecretField{Key: id1 + "-token"},
				Channel: "-1001406363649",
			},
			err: nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.src.Valid()
			influxTesting.ErrorsEqual(t, got, c.err)
		})
	}
}

var timeGen1 = mock.TimeGenerator{FakeValue: time.Date(2006, time.July, 13, 4, 19, 10, 0, time.UTC)}
var timeGen2 = mock.TimeGenerator{FakeValue: time.Date(2006, time.July, 14, 5, 23, 53, 10, time.UTC)}

func TestJSON(t *testing.T) {
	cases := []struct {
		name string
		src  influxdb.NotificationEndpoint
	}{
		{
			name: "simple Slack",
			src: &endpoint.Slack{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				URL:   "https://slack.com/api/chat.postMessage",
				Token: influxdb.SecretField{Key: "token-key-1"},
			},
		},
		{
			name: "Slack without token",
			src: &endpoint.Slack{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				URL: "https://hooks.slack.com/services/x/y/z",
			},
		},
		{
			name: "simple pagerduty",
			src: &endpoint.PagerDuty{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				ClientURL:  "https://events.pagerduty.com/v2/enqueue",
				RoutingKey: influxdb.SecretField{Key: "pagerduty-routing-key"},
			},
		},
		{
			name: "simple http",
			src: &endpoint.HTTP{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				Headers: map[string]string{
					"x-header-1": "header 1",
					"x-header-2": "header 2",
				},
				AuthMethod: "basic",
				URL:        "http://example.com",
				Username:   influxdb.SecretField{Key: "username-key"},
				Password:   influxdb.SecretField{Key: "password-key"},
			},
		},
		{
			name: "simple Telegram",
			src: &endpoint.Telegram{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "nameTelegram",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				Token: influxdb.SecretField{Key: "token-key-1"},
			},
		},
	}
	for _, c := range cases {
		b, err := json.Marshal(c.src)
		if err != nil {
			t.Fatalf("%s marshal failed, err: %s", c.name, err.Error())
		}
		got, err := endpoint.UnmarshalJSON(b)
		if err != nil {
			t.Fatalf("%s unmarshal failed, err: %s", c.name, err.Error())
		}
		if diff := cmp.Diff(got, c.src); diff != "" {
			t.Errorf("failed %s, NotificationEndpoint are different -got/+want\ndiff %s", c.name, diff)
		}
	}
}

func TestBackFill(t *testing.T) {
	cases := []struct {
		name   string
		src    influxdb.NotificationEndpoint
		target influxdb.NotificationEndpoint
	}{
		{
			name: "simple Slack",
			src: &endpoint.Slack{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				URL: "https://slack.com/api/chat.postMessage",
				Token: influxdb.SecretField{
					Value: strPtr("token-value"),
				},
			},
			target: &endpoint.Slack{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				URL: "https://slack.com/api/chat.postMessage",
				Token: influxdb.SecretField{
					Key:   id1 + "-token",
					Value: strPtr("token-value"),
				},
			},
		},
		{
			name: "simple pagerduty",
			src: &endpoint.PagerDuty{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				ClientURL: "https://events.pagerduty.com/v2/enqueue",
				RoutingKey: influxdb.SecretField{
					Value: strPtr("routing-key-value"),
				},
			},
			target: &endpoint.PagerDuty{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				ClientURL: "https://events.pagerduty.com/v2/enqueue",
				RoutingKey: influxdb.SecretField{
					Key:   id1 + "-routing-key",
					Value: strPtr("routing-key-value"),
				},
			},
		},
		{
			name: "http with token",
			src: &endpoint.HTTP{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				AuthMethod: "basic",
				URL:        "http://example.com",
				Username: influxdb.SecretField{
					Value: strPtr("username1"),
				},
				Password: influxdb.SecretField{
					Value: strPtr("password1"),
				},
			},
			target: &endpoint.HTTP{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				AuthMethod: "basic",
				URL:        "http://example.com",
				Username: influxdb.SecretField{
					Key:   id1 + "-username",
					Value: strPtr("username1"),
				},
				Password: influxdb.SecretField{
					Key:   id1 + "-password",
					Value: strPtr("password1"),
				},
			},
		},
		{
			name: "simple Telegram",
			src: &endpoint.Telegram{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				Token: influxdb.SecretField{
					Value: strPtr("token-value"),
				},
			},
			target: &endpoint.Telegram{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				Token: influxdb.SecretField{
					Key:   id1 + "-token",
					Value: strPtr("token-value"),
				},
			},
		},
	}
	for _, c := range cases {
		c.src.BackfillSecretKeys()
		if diff := cmp.Diff(c.target, c.src); diff != "" {
			t.Errorf("failed %s, NotificationEndpoint are different -got/+want\ndiff %s", c.name, diff)
		}
	}
}

func TestSecretFields(t *testing.T) {
	cases := []struct {
		name    string
		src     influxdb.NotificationEndpoint
		secrets []influxdb.SecretField
	}{
		{
			name: "simple Slack",
			src: &endpoint.Slack{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				URL: "https://slack.com/api/chat.postMessage",
				Token: influxdb.SecretField{
					Key:   id1 + "-token",
					Value: strPtr("token-value"),
				},
			},
			secrets: []influxdb.SecretField{
				{
					Key:   id1 + "-token",
					Value: strPtr("token-value"),
				},
			},
		},
		{
			name: "simple pagerduty",
			src: &endpoint.PagerDuty{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				ClientURL: "https://events.pagerduty.com/v2/enqueue",
				RoutingKey: influxdb.SecretField{
					Key:   id1 + "-routing-key",
					Value: strPtr("routing-key-value"),
				},
			},
			secrets: []influxdb.SecretField{
				{
					Key:   id1 + "-routing-key",
					Value: strPtr("routing-key-value"),
				},
			},
		},
		{
			name: "http with user and password",
			src: &endpoint.HTTP{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				AuthMethod: "basic",
				URL:        "http://example.com",
				Username: influxdb.SecretField{
					Key:   id1 + "-username",
					Value: strPtr("user1"),
				},
				Password: influxdb.SecretField{
					Key:   id1 + "-password",
					Value: strPtr("password1"),
				},
			},
			secrets: []influxdb.SecretField{
				{
					Key:   id1 + "-username",
					Value: strPtr("user1"),
				},
				{
					Key:   id1 + "-password",
					Value: strPtr("password1"),
				},
			},
		},
		{
			name: "simple Telegram",
			src: &endpoint.Telegram{
				Base: endpoint.Base{
					ID:     influxTesting.MustIDBase16Ptr(id1),
					Name:   "name1",
					OrgID:  influxTesting.MustIDBase16Ptr(id3),
					Status: influxdb.Active,
					CRUDLog: influxdb.CRUDLog{
						CreatedAt: timeGen1.Now(),
						UpdatedAt: timeGen2.Now(),
					},
				},
				Token: influxdb.SecretField{
					Key:   id1 + "-token",
					Value: strPtr("token-value"),
				},
			},
			secrets: []influxdb.SecretField{
				{
					Key:   id1 + "-token",
					Value: strPtr("token-value"),
				},
			},
		},
	}
	for _, c := range cases {
		secretFields := c.src.SecretFields()
		if diff := cmp.Diff(c.secrets, secretFields); diff != "" {
			t.Errorf("failed %s, NotificationEndpoint are different -got/+want\ndiff %s", c.name, diff)
		}
	}
}

func strPtr(s string) *string {
	ss := new(string)
	*ss = s
	return ss
}
