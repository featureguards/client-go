package featureguards

import (
	"context"
	"os"
	"testing"
	"time"

	pb_ft "github.com/featureguards/featureguards-go/v2/proto/feature_toggle"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_hash(t *testing.T) {
	type args struct {
		name  string
		keys  []*pb_ft.Key
		attrs Attributes
	}
	time, err := time.Parse(time.RFC3339, "2019-10-12T07:20:50.52Z")
	if err != nil {
		t.Fatal(err)
	}
	attrs := Attributes(map[string]interface{}{}).
		Int("user_id", 123).
		String("company_slug", "FeatureGuards").
		Bool("is_admin", true).
		Time("created_at", time)
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{name: "basic int", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "user_id", KeyType: pb_ft.Key_FLOAT}}, attrs: attrs}, want: 4353148100880623749, wantErr: false},
		{name: "wrong int type", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "company_slug", KeyType: pb_ft.Key_FLOAT}}, attrs: attrs}, want: 0, wantErr: true},
		{name: "basic string", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "company_slug", KeyType: pb_ft.Key_STRING}}, attrs: attrs}, want: 15324770540884756055, wantErr: false},
		{name: "wrong string type", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "user_id", KeyType: pb_ft.Key_STRING}}, attrs: attrs}, want: 0, wantErr: true},
		{name: "basic bool", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "is_admin", KeyType: pb_ft.Key_BOOLEAN}}, attrs: attrs}, want: 15549163119024811594, wantErr: false},
		{name: "wrong bool type", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "user_id", KeyType: pb_ft.Key_BOOLEAN}}, attrs: attrs}, want: 0, wantErr: true},
		{name: "basic time", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "created_at", KeyType: pb_ft.Key_DATE_TIME}}, attrs: attrs}, want: 148043139920556009, wantErr: false},
		{name: "wrong time type", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "user_id", KeyType: pb_ft.Key_DATE_TIME}}, attrs: attrs}, want: 0, wantErr: true},
		{name: "missing key", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "user_iD", KeyType: pb_ft.Key_FLOAT}}, attrs: attrs}, want: 0, wantErr: true},
		{name: "empty key", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "", KeyType: pb_ft.Key_FLOAT}}, attrs: attrs.String("", "foo")}, want: 0, wantErr: true},
		{name: "unknown key type", args: args{name: "FOO", keys: []*pb_ft.Key{{Key: "user_id", KeyType: 20}}, attrs: attrs}, want: 0, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hash(tt.args.name, tt.args.keys, tt.args.attrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_match(t *testing.T) {
	type args struct {
		name    string
		matches []*pb_ft.Match
		attrs   Attributes
	}
	timestamp, err := time.Parse(time.RFC3339, "2019-10-12T07:20:50.52Z")
	if err != nil {
		t.Fatal(err)
	}
	attrs := Attributes(map[string]interface{}{}).
		Int("user_id", 123).
		String("company_slug", "FeatureGuards").
		Bool("is_admin", true).
		Time("created_at", timestamp)

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "string eq matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "company_slug", KeyType: pb_ft.Key_STRING}, Operation: &pb_ft.Match_StringOp{StringOp: &pb_ft.StringOp{Op: pb_ft.StringOp_EQ, Values: []string{"FeatureGuards"}}}}}}, want: true, wantErr: false},
		{name: "string eq mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "company_slug", KeyType: pb_ft.Key_STRING}, Operation: &pb_ft.Match_StringOp{StringOp: &pb_ft.StringOp{Op: pb_ft.StringOp_EQ, Values: []string{"featureguards"}}}}}}, want: false, wantErr: false},
		{name: "string eq empty", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "company_slug", KeyType: pb_ft.Key_STRING}, Operation: &pb_ft.Match_StringOp{StringOp: &pb_ft.StringOp{Op: pb_ft.StringOp_EQ}}}}}, want: false, wantErr: true},
		{name: "string contains matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "company_slug", KeyType: pb_ft.Key_STRING}, Operation: &pb_ft.Match_StringOp{StringOp: &pb_ft.StringOp{Op: pb_ft.StringOp_CONTAINS, Values: []string{"Guards"}}}}}}, want: true, wantErr: false},
		{name: "string in matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "company_slug", KeyType: pb_ft.Key_STRING}, Operation: &pb_ft.Match_StringOp{StringOp: &pb_ft.StringOp{Op: pb_ft.StringOp_IN, Values: []string{"foo", "FeatureGuards"}}}}}}, want: true, wantErr: false},
		{name: "string in mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "company_slug", KeyType: pb_ft.Key_STRING}, Operation: &pb_ft.Match_StringOp{StringOp: &pb_ft.StringOp{Op: pb_ft.StringOp_IN, Values: []string{"foo"}}}}}}, want: false, wantErr: false},
		{name: "boolean matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "is_admin", KeyType: pb_ft.Key_BOOLEAN}, Operation: &pb_ft.Match_BoolOp{BoolOp: &pb_ft.BoolOp{Value: true}}}}}, want: true, wantErr: false},
		{name: "boolean mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "is_admin", KeyType: pb_ft.Key_BOOLEAN}, Operation: &pb_ft.Match_BoolOp{BoolOp: &pb_ft.BoolOp{Value: false}}}}}, want: false, wantErr: false},
		{name: "float eq matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{123}, Op: pb_ft.FloatOp_EQ}}}}}, want: true, wantErr: false},
		{name: "float eq mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{1234}, Op: pb_ft.FloatOp_EQ}}}}}, want: false, wantErr: false},
		{name: "float neq matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{123}, Op: pb_ft.FloatOp_NEQ}}}}}, want: false, wantErr: false},
		{name: "float neq mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{1234}, Op: pb_ft.FloatOp_NEQ}}}}}, want: true, wantErr: false},
		{name: "float gt matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{122}, Op: pb_ft.FloatOp_GT}}}}}, want: true, wantErr: false},
		{name: "float gt mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{123}, Op: pb_ft.FloatOp_GT}}}}}, want: false, wantErr: false},
		{name: "float gte matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{123}, Op: pb_ft.FloatOp_GTE}}}}}, want: true, wantErr: false},
		{name: "float gte mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{124}, Op: pb_ft.FloatOp_GTE}}}}}, want: false, wantErr: false},
		{name: "float lt matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{124}, Op: pb_ft.FloatOp_LT}}}}}, want: true, wantErr: false},
		{name: "float lt mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{123}, Op: pb_ft.FloatOp_LT}}}}}, want: false, wantErr: false},
		{name: "float lte matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{123}, Op: pb_ft.FloatOp_LTE}}}}}, want: true, wantErr: false},
		{name: "float lte mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{122}, Op: pb_ft.FloatOp_LTE}}}}}, want: false, wantErr: false},
		{name: "float in matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{1, 123}, Op: pb_ft.FloatOp_IN}}}}}, want: true, wantErr: false},
		{name: "float in mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Values: []float32{2, 122}, Op: pb_ft.FloatOp_IN}}}}}, want: false, wantErr: false},
		{name: "time after matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "created_at", KeyType: pb_ft.Key_DATE_TIME}, Operation: &pb_ft.Match_DateTimeOp{DateTimeOp: &pb_ft.DateTimeOp{Timestamp: timestamppb.New(timestamp.Add(-1 * time.Second)), Op: pb_ft.DateTimeOp_AFTER}}}}}, want: true, wantErr: false},
		{name: "time after mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "created_at", KeyType: pb_ft.Key_DATE_TIME}, Operation: &pb_ft.Match_DateTimeOp{DateTimeOp: &pb_ft.DateTimeOp{Timestamp: timestamppb.New(timestamp), Op: pb_ft.DateTimeOp_AFTER}}}}}, want: false, wantErr: false},
		{name: "time before matches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "created_at", KeyType: pb_ft.Key_DATE_TIME}, Operation: &pb_ft.Match_DateTimeOp{DateTimeOp: &pb_ft.DateTimeOp{Timestamp: timestamppb.New(timestamp.Add(1 * time.Second)), Op: pb_ft.DateTimeOp_BEFORE}}}}}, want: true, wantErr: false},
		{name: "time before mismatches", args: args{name: "FOO", attrs: attrs, matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "created_at", KeyType: pb_ft.Key_DATE_TIME}, Operation: &pb_ft.Match_DateTimeOp{DateTimeOp: &pb_ft.DateTimeOp{Timestamp: timestamppb.New(timestamp), Op: pb_ft.DateTimeOp_BEFORE}}}}}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := match(tt.args.name, tt.args.matches, tt.args.attrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isOn(t *testing.T) {
	type args struct {
		ft      *pb_ft.FeatureToggle
		options []FeatureToggleOptions
	}
	timestamp, err := time.Parse(time.RFC3339, "2019-10-12T07:20:50.52Z")
	if err != nil {
		t.Fatal(err)
	}
	attrs := Attributes(map[string]interface{}{}).
		Int("user_id", 123). // hashes to about 62%
		String("company_slug", "FeatureGuards").
		Bool("is_admin", true).
		Time("created_at", timestamp)

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "on/off deleted", args: args{ft: &pb_ft.FeatureToggle{Name: "FOO", DeletedAt: timestamppb.New(timestamp), ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{On: &pb_ft.Variant{Weight: 100}, Off: &pb_ft.Variant{Weight: 0}}}}}, want: false, wantErr: true},
		{name: "on/off disabled", args: args{ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: false, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{On: &pb_ft.Variant{Weight: 100}, Off: &pb_ft.Variant{Weight: 0}}}}}, want: false, wantErr: false},
		{name: "on/off on", args: args{ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{On: &pb_ft.Variant{Weight: 100}, Off: &pb_ft.Variant{Weight: 0}}}}}, want: true, wantErr: false},
		{name: "on/off errors on partial weight", args: args{ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{On: &pb_ft.Variant{Weight: 1}, Off: &pb_ft.Variant{Weight: 99}}}}}, want: false, wantErr: true},
		{name: "on/off errors on equal weights", args: args{ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{On: &pb_ft.Variant{Weight: 0}, Off: &pb_ft.Variant{Weight: 0}}}}}, want: false, wantErr: true},
		{name: "on/off errors on missing off", args: args{ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{On: &pb_ft.Variant{Weight: 0}}}}}, want: false, wantErr: true},
		{name: "on/off off", args: args{ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{On: &pb_ft.Variant{Weight: 0}, Off: &pb_ft.Variant{Weight: 100}}}}}, want: false, wantErr: false},
		{name: "on/off off 100 weight", args: args{ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{On: &pb_ft.Variant{Weight: 0}, Off: &pb_ft.Variant{Weight: 100}}}}}, want: false, wantErr: false},
		{name: "on/off on allowlist", args: args{options: []FeatureToggleOptions{WithAttributes(attrs)}, ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{Off: &pb_ft.Variant{Weight: 100}, On: &pb_ft.Variant{Weight: 0, Matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Op: pb_ft.FloatOp_EQ, Values: []float32{123}}}}}}}}}}, want: true, wantErr: false},
		{name: "on/off off allowlist", args: args{options: []FeatureToggleOptions{WithAttributes(attrs)}, ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{Off: &pb_ft.Variant{Weight: 100}, On: &pb_ft.Variant{Weight: 0, Matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Op: pb_ft.FloatOp_EQ, Values: []float32{1234}}}}}}}}}}, want: false, wantErr: false},
		{name: "on/off off disallowlist", args: args{options: []FeatureToggleOptions{WithAttributes(attrs)}, ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_ON_OFF, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_OnOff{OnOff: &pb_ft.OnOffFeature{On: &pb_ft.Variant{Weight: 100}, Off: &pb_ft.Variant{Matches: []*pb_ft.Match{{Key: &pb_ft.Key{Key: "user_id", KeyType: pb_ft.Key_FLOAT}, Operation: &pb_ft.Match_FloatOp{FloatOp: &pb_ft.FloatOp{Op: pb_ft.FloatOp_EQ, Values: []float32{123}}}}}}}}}}, want: false, wantErr: false},

		// Percentage
		{name: "percentage random on 100", args: args{ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_PERCENTAGE, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_Percentage{Percentage: &pb_ft.PercentageFeature{On: &pb_ft.Variant{Weight: 100}, Off: &pb_ft.Variant{Weight: 0}, Stickiness: &pb_ft.Stickiness{StickinessType: pb_ft.Stickiness_RANDOM}}}}}, want: true, wantErr: false},
		{name: "percentage random off", args: args{options: []FeatureToggleOptions{WithAttributes(attrs)}, ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_PERCENTAGE, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_Percentage{Percentage: &pb_ft.PercentageFeature{On: &pb_ft.Variant{Weight: 0}, Off: &pb_ft.Variant{Weight: 100}, Stickiness: &pb_ft.Stickiness{StickinessType: pb_ft.Stickiness_KEYS, Keys: []*pb_ft.Key{{Key: "user_id", KeyType: pb_ft.Key_FLOAT}}}}}}}, want: false, wantErr: false},
		{name: "percentage sticky 70%", args: args{options: []FeatureToggleOptions{WithAttributes(attrs)}, ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_PERCENTAGE, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_Percentage{Percentage: &pb_ft.PercentageFeature{On: &pb_ft.Variant{Weight: 70}, Off: &pb_ft.Variant{Weight: 30}, Stickiness: &pb_ft.Stickiness{StickinessType: pb_ft.Stickiness_KEYS, Keys: []*pb_ft.Key{{Key: "user_id", KeyType: pb_ft.Key_FLOAT}}}}}}}, want: true, wantErr: false},
		{name: "percentage sticky 60%", args: args{options: []FeatureToggleOptions{WithAttributes(attrs)}, ft: &pb_ft.FeatureToggle{Name: "FOO", ToggleType: pb_ft.FeatureToggle_PERCENTAGE, Enabled: true, FeatureDefinition: &pb_ft.FeatureToggle_Percentage{Percentage: &pb_ft.PercentageFeature{On: &pb_ft.Variant{Weight: 60}, Off: &pb_ft.Variant{Weight: 40}, Stickiness: &pb_ft.Stickiness{StickinessType: pb_ft.Stickiness_KEYS, Keys: []*pb_ft.Key{{Key: "user_id", KeyType: pb_ft.Key_FLOAT}}}}}}}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isOn(tt.args.ft, tt.args.options...)
			if (err != nil) != tt.wantErr {
				t.Errorf("isOn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isOn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_refreshTokens(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr := os.Getenv("GRPC_ADDR")
	apiKey := os.Getenv("API_KEY")
	ft := New(ctx, withDomain(addr), withoutListen(), withTestCerts(),
		WithApiKey(apiKey), WithDefaults(map[string]bool{"BAR": true}))
	ft.ft.mu.Lock()
	accessToken := ft.ft.accessToken
	refreshToken := ft.ft.refreshToken
	ft.ft.mu.Unlock()
	err := ft.ft.refreshTokens(ctx)
	if err != nil {
		t.Errorf("refreshTokens() error = %v, wantErr nil", err)
	}
	ft.ft.mu.Lock()
	defer ft.ft.mu.Unlock()
	if accessToken == ft.ft.accessToken {
		t.Errorf("accessToken = %v, want != %v", ft.ft.accessToken, accessToken)
	}
	if refreshToken == ft.ft.refreshToken {
		t.Errorf("refreshToken = %v, want != %v", ft.ft.refreshToken, refreshToken)
	}
}

func Test_IsOn(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addr := os.Getenv("GRPC_ADDR")
	apiKey := os.Getenv("API_KEY")
	ft := New(ctx, withDomain(addr), withoutListen(), withTestCerts(),
		WithApiKey(apiKey), WithDefaults(map[string]bool{"BAR": true}))
	if ft.ft == nil {
		t.Error("feature toggles should not be nil")
	}
	on, err := ft.IsOn("TEST")
	if err != nil {
		t.Errorf("IsOn() error = %v, wantErr nil", err)
	}
	if !on {
		t.Errorf("on() on = %v, wantOn %v", on, true)
	}

	on, err = ft.IsOn("BAR")
	if err == nil {
		t.Errorf("IsOn() error = %v, wantErr not found", err)
	}
	if !on {
		t.Errorf("on() on = %v, wantOn %v", on, true)
	}

}
