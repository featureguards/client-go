package dynamic_settings

import (
	"testing"

	pb_ds "github.com/featureguards/featureguards-go/v2/proto/dynamic_setting"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_basic(t *testing.T) {
	ds := New()
	var boolV bool
	var intV int64
	var stringV string
	ds.Bool("FOO", &boolV, WithDefaultBool(true))
	ds.Int64("BAR", &intV, WithDefaultInt64(10))
	ds.String("DELTA", &stringV, WithDefaultString("123"))

	if boolV != true {
		t.Errorf("Bool(FOO) = %v, wantOn %v", boolV, true)
	}

	if intV != 10 {
		t.Errorf("Int64(BAR) = %v, wantOn %v", intV, 10)
	}

	if stringV != "123" {
		t.Errorf("String(DELTA) = %v, wantOn %v", stringV, "123")
	}

	// Process
	ds.Process([]*pb_ds.DynamicSetting{{Name: "FOO", SettingType: pb_ds.DynamicSetting_BOOL, SettingDefinition: &pb_ds.DynamicSetting_BoolValue{BoolValue: &pb_ds.BoolValue{Value: false}}}, {
		Name: "BAR", DeletedAt: timestamppb.Now(),
	}})
	if boolV != false {
		t.Errorf("Bool(FOO) = %v, wantOn %v", boolV, false)
	}

	// Deleted BAR. Make sure it doesn't update.
	ds.Process([]*pb_ds.DynamicSetting{{Name: "BAR", SettingType: pb_ds.DynamicSetting_INTEGER, SettingDefinition: &pb_ds.DynamicSetting_IntegerValue{IntegerValue: &pb_ds.IntegerValue{Value: 100}}}})
	if intV != 10 {
		t.Errorf("Int64(BAR) = %v, wantOn %v", intV, 10)
	}

}
