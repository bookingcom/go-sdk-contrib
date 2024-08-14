package model_test

import (
	"github.com/open-feature/go-sdk-contrib/providers/go-feature-flag/pkg/model"
	of "github.com/open-feature/go-sdk/openfeature"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewFeatureEvent(t *testing.T) {
	type args struct {
		user      of.EvaluationContext
		flagKey   string
		value     interface{}
		variation string
		failed    bool
		version   string
		source    string
	}
	tests := []struct {
		name string
		args args
		want model.FeatureEvent
	}{
		{
			name: "anonymous user",
			args: args{
				user:      of.NewEvaluationContext("ABCD", map[string]interface{}{"anonymous": true}),
				flagKey:   "random-key",
				value:     "YO",
				variation: "Default",
				failed:    false,
				version:   "",
				source:    "SERVER",
			},
			want: model.FeatureEvent{
				Kind: "feature", ContextKind: "anonymousUser", UserKey: "ABCD", CreationDate: time.Now().Unix(), Key: "random-key",
				Variation: "Default", Value: "YO", Default: false, Source: "SERVER",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, model.NewFeatureEvent(tt.args.user, tt.args.flagKey, tt.args.value, tt.args.variation, tt.args.failed, tt.args.version, tt.args.source), "NewFeatureEvent(%v, %v, %v, %v, %v, %v, %V)", tt.args.user, tt.args.flagKey, tt.args.value, tt.args.variation, tt.args.failed, tt.args.version, tt.args.source)
		})
	}
}

func TestFeatureEvent_MarshalInterface(t *testing.T) {
	tests := []struct {
		name         string
		featureEvent *model.FeatureEvent
		want         *model.FeatureEvent
		wantErr      bool
	}{
		{
			name: "happy path",
			featureEvent: &model.FeatureEvent{
				Kind:         "feature",
				ContextKind:  "anonymousUser",
				UserKey:      "ABCD",
				CreationDate: 1617970547,
				Key:          "random-key",
				Variation:    "Default",
				Value: map[string]interface{}{
					"string": "string",
					"bool":   true,
					"float":  1.23,
					"int":    1,
				},
				Default: false,
			},
			want: &model.FeatureEvent{
				Kind:         "feature",
				ContextKind:  "anonymousUser",
				UserKey:      "ABCD",
				CreationDate: 1617970547,
				Key:          "random-key",
				Variation:    "Default",
				Value:        `{"bool":true,"float":1.23,"int":1,"string":"string"}`,
				Default:      false,
			},
		},
		{
			name: "marshal failed",
			featureEvent: &model.FeatureEvent{
				Kind:         "feature",
				ContextKind:  "anonymousUser",
				UserKey:      "ABCD",
				CreationDate: 1617970547,
				Key:          "random-key",
				Variation:    "Default",
				Value:        make(chan int),
				Default:      false,
			},
			wantErr: true,
		},
		{
			name:         "nil featureEvent",
			featureEvent: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.featureEvent.MarshalInterface(); (err != nil) != tt.wantErr {
				t.Errorf("FeatureEvent.MarshalInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				assert.Equal(t, tt.want, tt.featureEvent)
			}
		})
	}
}