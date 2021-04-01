package decode_test

import (
	"fmt"
	"testing"

	"go.einride.tech/can"
	"go.einride.tech/can/internal/generate"
	"go.einride.tech/can/pkg/descriptor"
)

var (
	db  = getDatabase()
	dbc = []byte(`
VERSION ""
NS_ :
BS_:
BU_: DBG DRIVER IO MOTOR SENSOR

BO_ 1 EmptyMessage: 0 DBG
BO_ 100 DriverHeartbeat: 1 DRIVER
 SG_ Command : 0|8@1+ (1,0) [0|0] "" SENSOR,MOTOR
BO_ 101 MotorCommand: 1 DRIVER
 SG_ Steer : 0|4@1- (1,-5) [-5|5] "" MOTOR
 SG_ Drive : 4|4@1+ (1,0) [0|9] "" MOTOR
BO_ 400 MotorStatus: 3 MOTOR
 SG_ WheelError : 0|1@1+ (1,0) [0|0] "" DRIVER,IO
 SG_ SpeedKph : 8|16@1+ (0.001,0) [0|0] "km/h" DRIVER,IO

BO_ 200 SensorSonars: 8 SENSOR
 SG_ Mux M : 0|4@1+ (1,0) [0|0] "" DRIVER,IO
 SG_ ErrCount : 4|12@1+ (1,0) [0|0] "" DRIVER,IO
 SG_ Left m0 : 16|12@1+ (0.1,0) [0|0] "" DRIVER,IO
 SG_ Middle m0 : 28|12@1+ (0.1,0) [0|0] "" DRIVER,IO
 SG_ Right m0 : 40|12@1+ (0.1,0) [0|0] "" DRIVER,IO
 SG_ Rear m0 : 52|12@1+ (0.1,0) [0|0] "" DRIVER,IO
 SG_ NoFiltLeft m1 : 16|12@1+ (0.1,0) [0|0] "" DBG
 SG_ NoFiltMiddle m1 : 28|12@1+ (0.1,0) [0|0] "" DBG
 SG_ NoFiltRight m1 : 40|12@1+ (0.1,0) [0|0] "" DBG
 SG_ NoFiltRear m1 : 52|12@1+ (0.1,0) [0|0] "" DBG

BO_ 500 IODebug: 6 IO
 SG_ TestUnsigned : 0|8@1+ (1,0) [0|0] "" DBG
 SG_ TestEnum : 8|6@1+ (1,0) [0|0] "" DBG
 SG_ TestSigned : 16|8@1- (1,0) [0|0] "" DBG
 SG_ TestFloat : 24|8@1+ (0.5,0) [0|0] "" DBG
 SG_ TestBoolEnum : 32|1@1+ (1,0) [0|0] "" DBG
 SG_ TestScaledEnum : 40|2@1+ (2,0) [0|6] "" DBG

BO_ 1530 DisconnectState: 14 MOTOR
 SG_ LockCountRearRight : 91|20@0+ (1,0) [0|1048575] ""  IO
 SG_ DisconnectStateRearRight : 95|4@0+ (1,0) [0|5] ""  IO
 SG_ CurrentRearRight : 79|16@0+ (1,0) [0|65535] ""  IO
 SG_ DisconnectStateRearRightTarget : 64|1@0+ (1,0) [0|1] ""  IO
 SG_ TargetSpeedRearRight : 63|15@0+ (0.125,-2048) [-2048|2047.875] "rad/s"  IO
 SG_ LockCountRearLeft : 35|20@0+ (1,0) [0|1048575] ""  IO
 SG_ DisconnectStateRearLeft : 39|4@0+ (1,0) [0|5] ""  IO
 SG_ CurrentRearLeft : 23|16@0+ (1,0) [0|65535] ""  IO
 SG_ DisconnectStateRearLeftTarget : 8|1@0+ (1,0) [0|1] ""  IO
 SG_ TargetSpeedRearLeft : 7|15@0+ (0.125,-2048) [-2048|2047.875] "rad/s"  IO

BO_ 1927 EvaporatorVariables: 14 SENSOR
 SG_ EvaporatorAirTempIn : 96|8@0+ (0.4,-20) [-19.6|82.4] "degC"  IO
 SG_ CompressorSpeedActual : 80|16@0+ (1,-1) [0|8600] "RPM"  IO
 SG_ FBTargetCompressorSpeed : 64|16@0+ (1,-8600) [-8599|8600] "RPM"  IO
 SG_ FFTargetCompressorSpeed : 48|16@0+ (1,-1) [0|8600] "RPM"  IO
 SG_ FBTargetSuctionPressure : 32|16@0+ (0.05,0) [0.05|3276.8] "kPa"  IO
 SG_ FFTargetSuctionPressure : 16|16@0+ (0.05,0) [0.05|3276.8] "kPa"  IO
 SG_ ChillerCoolantOutTempTarget : 8|8@0+ (0.4,-20) [-19.6|82.4] "degC"  IO
 SG_ EvaporatorAirOutTempTarget : 0|8@0+ (0.4,-20) [-19.6|82.4] "degC"  IO

EV_ BrakeEngaged: 0 [0|1] "" 0 10 DUMMY_NODE_VECTOR0 Vector__XXX;
EV_ Torque: 1 [0|30000] "mNm" 500 16 DUMMY_NODE_VECTOR0 Vector__XXX;
CM_ EV_ BrakeEngaged "Brake fully engaged";
CM_ BU_ DRIVER "The driver controller driving the car";
CM_ BU_ MOTOR "The motor controller of the car";
CM_ BU_ SENSOR "The sensor controller of the car";
CM_ BO_ 100 "Sync message used to synchronize the controllers";

BA_DEF_ "BusType" STRING ;
BA_DEF_ BO_  "GenMsgSendType" ENUM  "None","Cyclic","OnEvent";
BA_DEF_ BO_ "GenMsgCycleTime" INT 0 0;
BA_DEF_ SG_  "FieldType" STRING ;
BA_DEF_ SG_  "GenSigStartValue" INT 0 10000;
BA_DEF_DEF_ "BusType" "CAN";
BA_DEF_DEF_ "FieldType" "";
BA_DEF_DEF_ "GenMsgCycleTime" 0;
BA_DEF_DEF_ "GenSigStartValue" 0;

BA_ "GenMsgSendType" BO_ 1 0;
BA_ "GenMsgSendType" BO_ 100 1;
BA_ "GenMsgCycleTime" BO_ 100 1000;
BA_ "GenMsgSendType" BO_ 101 1;
BA_ "GenMsgCycleTime" BO_ 101 100;
BA_ "GenMsgSendType" BO_ 200 1;
BA_ "GenMsgCycleTime" BO_ 200 100;
BA_ "GenMsgSendType" BO_ 400 1;
BA_ "GenMsgCycleTime" BO_ 400 100;
BA_ "GenMsgSendType" BO_ 500 2;
BA_ "FieldType" SG_ 100 Command "Command";
BA_ "FieldType" SG_ 500 TestEnum "TestEnum";
BA_ "GenSigStartValue" SG_ 500 TestEnum 2;

VAL_ 100 Command 2 "Reboot" 1 "Sync" 0 "None" ;
VAL_ 500 TestEnum 2 "Two" 1 "One" ;
VAL_ 500 TestScaledEnum 3 "Six" 2 "Four" 1 "Two" 0 "Zero" ;
VAL_ 500 TestBoolEnum 1 "One" 0 "Zero" ;
VAL_ 1530 DisconnectStateRearRight 0 "Undefined" 1 "Locked" 2 "Unlocked" 3 "Locking" 4 "Unlocking" 5 "Faulted" ;
VAL_ 1530 DisconnectStateRearLeft 0 "Undefined" 1 "Locked" 2 "Unlocked" 3 "Locking" 4 "Unlocking" 5 "Faulted" ;

`)
)

type signal struct {
	name        string
	value       float64
	description string
	unit        string
}

func getDatabase() descriptor.Database {
	c, _ := generate.Compile("test.dbc", dbc)
	db := *c.Database
	return db
}

func TestDecodeEvaporatorVariables(t *testing.T) {

	c, err := generate.Compile("test.dbc", dbc)
	if err != nil {
		t.Errorf("err = %v; want nil", err)
	}

	db := *c.Database
	message, _ := db.Message(uint32(1927))

	canDataHexString := "008232204e027600ca4b0007d296"

	payload, err := can.PayloadFromHex(canDataHexString)
	// ci := packet.Metadata().CaptureInfo
	if err != nil {
		t.Errorf("err = %v; want nil", err)
	}

	fmt.Println(payload.Hex())

	expected := []signal{
		{
			name:        "EvaporatorAirOutTempTarget",
			value:       6.0,
			description: "",
			unit:        "degC",
		},
		{
			name:        "ChillerCoolantOutTempTarget",
			value:       -10.0,
			description: "",
			unit:        "degC",
		},
		{
			name:        "FFTargetSuctionPressure",
			value:       206.75,
			description: "",
			unit:        "kPa",
		},
		{
			name:        "FBTargetSuctionPressure",
			value:       15.75,
			description: "",
			unit:        "kPa",
		},
		{
			name:        "FFTargetCompressorSpeed",
			value:       100,
			description: "",
			unit:        "RPM",
		},
		{
			name:        "FBTargetCompressorSpeed",
			value:       1000,
			description: "",
			unit:        "RPM",
		},
		{
			name:        "CompressorSpeedActual",
			value:       1000,
			description: "",
			unit:        "RPM",
		},
		{
			name:        "EvaporatorAirTempIn",
			value:       10.0,
			description: "",
			unit:        "degC",
		},
	}

	var expectedMap = make(map[string]signal)
	for _, s := range expected {
		expectedMap[s.name] = s
	}

	for _, signal := range message.Signals {
		value := signal.UnmarshalPhysicalPayload(&payload)
		valueDesc, _ := signal.UnmarshalValueDescriptionPayload(&payload)
		name := signal.Name

		if value != expectedMap[name].value {
			t.Errorf("signal[%s] value = %f ; expected: %f", name, value, expectedMap[name].value)
		}

		if valueDesc != expectedMap[name].description {
			t.Errorf("signal[%s] value = %s ; expected: %s", name, valueDesc, expectedMap[name].description)
		}
	}

	for _, signal := range message.Decode(&payload) {
		value := signal.Value
		valueDesc := signal.Description
		name := signal.Signal.Name

		if value != expectedMap[name].value {
			t.Errorf("signal[%s] value = %f ; expected: %f", name, value, expectedMap[name].value)
		}

		if valueDesc != expectedMap[name].description {
			t.Errorf("signal[%s] value = %s ; expected: %s", name, valueDesc, expectedMap[name].description)
		}
	}
}

func TestDecodeDisconnectState(t *testing.T) {
	byteStringHex := "8000000420061880000005200600"

	p, err := can.PayloadFromHex(byteStringHex)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.Data)

	expected := []signal{
		{

			name:        "TargetSpeedRearLeft",
			value:       0.0,
			description: "",
			unit:        "rad/s",
		},
		{

			name:        "DisconnectStateRearLeftTarget",
			value:       0,
			description: "",
			unit:        "",
		},
		{

			name:        "CurrentRearLeft",
			value:       4,
			description: "",
			unit:        "",
		},
		{

			name:        "DisconnectStateRearLeft",
			value:       2,
			description: "Unlocked",
			unit:        "",
		},
		{

			name:        "LockCountRearLeft",
			value:       1560,
			description: "",
			unit:        "",
		},
		{

			name:        "TargetSpeedRearRight",
			value:       0,
			description: "",
			unit:        "rad/s",
		},
		{

			name:        "DisconnectStateRearRightTarget",
			value:       0,
			description: "",
			unit:        "",
		},
		{

			name:        "CurrentRearRight",
			value:       5,
			description: "",
			unit:        "",
		},
		{

			name:        "DisconnectStateRearRight",
			value:       2,
			description: "Unlocked",
			unit:        "",
		},
		{

			name:        "LockCountRearRight",
			value:       1536,
			description: "",
			unit:        "",
		},
	}

	var expectedMap = make(map[string]signal)
	for _, s := range expected {
		expectedMap[s.name] = s
	}

	message, _ := db.Message(uint32(1530))
	for _, signal := range message.Signals {
		value := signal.UnmarshalPhysicalPayload(&p)
		valueDesc, _ := signal.UnmarshalValueDescriptionPayload(&p)
		name := signal.Name

		if value != expectedMap[name].value {
			t.Errorf("signal[%s] value = %f ; expected: %f", name, value, expectedMap[name].value)
		}

		if valueDesc != expectedMap[name].description {
			t.Errorf("signal[%s] value = %s ; expected: %s", name, valueDesc, expectedMap[name].description)
		}
	}
}

func TestDecodeSensorSonarsData(t *testing.T) {

	data := can.Data{0x01, 0x01, 0x01, 0x02, 0x01, 0x00}
	payload := can.Payload{Data: data[:]}

	message, _ := db.Message(uint32(500))
	fmt.Println(message.Length)

	for _, signal := range message.Signals {
		value := signal.UnmarshalPhysicalPayload(&payload)
		valueDesc, _ := signal.UnmarshalValueDescriptionPayload(&payload)

		valueFromData := signal.UnmarshalPhysical(data)
		descFromData, _ := signal.UnmarshalValueDescription(data)
		name := signal.Name

		if value != valueFromData {
			t.Errorf("signal[%s] value = %f ; expected: %f", name, value, valueFromData)
		}

		if valueDesc != descFromData {
			t.Errorf("signal[%s] value = %s ; expected: %s", name, valueDesc, descFromData)
		}
	}
}

func TestMessageDecodeSensorSonarsData(t *testing.T) {

	data := can.Data{0x01, 0x01, 0x01, 0x02, 0x01, 0x00}
	payload := can.Payload{Data: data[:]}

	message, _ := db.Message(uint32(500))
	fmt.Println(message.Length)

	decodedSignals := message.Decode(&payload)
	for _, signal := range decodedSignals {
		value := signal.Value
		valueDesc := signal.Description
		name := signal.Signal.Name

		valueFromData := signal.Signal.UnmarshalPhysical(data)
		descFromData, _ := signal.Signal.UnmarshalValueDescription(data)

		if value != valueFromData {
			t.Errorf("signal[%s] value = %f ; expected: %f", name, value, valueFromData)
		}

		if valueDesc != descFromData {
			t.Errorf("signal[%s] value = %s ; expected: %s", name, valueDesc, descFromData)
		}
	}
}

func BenchmarkDecodeData(b *testing.B) {

	message, _ := db.Message(uint32(500))
	decodeSignal := func() {
		data := can.Data{0x01, 0x01, 0x01, 0x02, 0x01, 0x00}
		for _, signal := range message.Signals {
			_ = signal.UnmarshalPhysical(data)
			_, _ = signal.UnmarshalValueDescription(data)
		}
	}
	for i := 0; i < b.N; i++ {
		decodeSignal()
	}
}

func BenchmarkDecodePayload(b *testing.B) {

	// {0x01, 0x01, 0x01, 0x02, 0x01, 0x00}

	message, _ := db.Message(uint32(500))
	decodeSignal := func() {
		data := can.Payload{Data: []byte{0x01, 0x01, 0x01, 0x02, 0x01, 0x00}}
		for _, signal := range message.Signals {
			_ = signal.UnmarshalPhysicalPayload(&data)
			_, _ = signal.UnmarshalValueDescriptionPayload(&data)
		}
	}
	for i := 0; i < b.N; i++ {
		decodeSignal()
	}
}
