package device

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/AutogrowSystems/go-intelli/util/encoding"
	"github.com/AutogrowSystems/go-intelli/util/tell"
	"github.com/snksoft/crc"
)

const (
	doserRaise    = "raise"
	doserLower    = "lower"
	doserModeNone = "none"
	doserModeBoth = "both"

	fan1Function         = "fan_1"
	fan2Function         = "fan_2"
	airConFunction       = "air_conditioner"
	co2InjectionFunction = "co2_injection"
	co2ExtractFunction   = "co2_extraction"
	heaterFunction       = "heater"
	dehumidifierFunction = "dehumidifier"
	humidifierFunction   = "humidifier"
	lightBank1Function   = "light_bank_1"
	lightBank2Function   = "light_bank_2"
	foggerFunction       = "pulsed_fogger"
	purgingFunction      = "purge"

	phFunction = "ph"

	nutrientDosingFunction = "Nutrient Dosing"
	waterFunction          = "Water"

	irrigationFunction = "irrigation"

	irrigationStation1Function = "Irrigation Station 1"
	irrigationStation2Function = "Irrigation Station 2"
	irrigationStation3Function = "Irrigation Station 3"
	irrigationStation4Function = "Irrigation Station 4"

	irrigationModeSingle      = "single"
	irrigationModeSequential  = "sequential"
	irrigationModeIndependent = "independent"

	irrigationModeDuringDayOnly = "during_day_only"
	irrigationModeSameTime      = "same_time"
	irrigationModeDayNight      = "day_night"

	temperatureC = "C"
	temperatureF = "F"

	nutrientConfigEC  = "EC"
	nutrientConfigCF  = "CF"
	nutrientConfigTDS = "TDS"

	dateFormat    = "DD/MM/YY"
	dateFormatUSA = "MM/DD/YY"

	dehumidifyAirCon = "air_conditioner"
	dehumidifyPurge  = "purge"
	dehumidifyNone   = doserModeNone

	requestLength = 64

	valueUndefined = 32768.0
)

var (
	d0Request = []byte{
		0x00,
		0x44, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7C, 0x54}

	d1Request = []byte{
		0x00,
		0x44, 0x31, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81, 0x94}

	d2Request = []byte{
		0x00,
		0x44, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x85, 0x95}

	d3Request = []byte{
		0x00,
		0x44, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x78, 0x55}
)

// IClimate
type iClimateShadow struct {
	State StateIClimate `json:"state"`
}

// StateIClimate represents the State data structure from an IntelliClimate packet
type StateIClimate struct {
	Reported ReportedIClimate `json:"reported"`
}

// ReportedIClimate represents the Reported data structure from an IntelliClimate packet
type ReportedIClimate struct {
	Config    ConfigIClimate  `json:"config"`
	Metrics   MetricsIClimate `json:"metrics"`
	Status    StatusIClimate  `json:"status"`
	Source    string          `json:"source"`
	Device    string          `json:"device"`
	Timestamp int64           `json:"timestamp"`
	Connected bool            `json:"connected"`
}

// ConfigIClimate represents the Config data structure from an IntelliClimate packet
type ConfigIClimate struct {
	Units     UnitsIClimate     `json:"units"`
	Functions FunctionsIClimate `json:"functions"`
	Advanced  AdvancedIClimate  `json:"advanced"`
	General   GeneralIClimate   `json:"general"`
}

// MetricsIClimate represents the Metrics data structure from an IntelliClimate packet
type MetricsIClimate struct {
	AirTemp        float64 `json:"air_temp"`
	DayNight       string  `json:"day_night"`
	FailSafeAlarms bool    `json:"fail_safe_alarms"`
	Light          float64 `json:"light"`
	PowerFail      bool    `json:"power_fail"`
	Rh             float64 `json:"rh"`
	Vpd            float64 `json:"vpd"`
	Co2            float64 `json:"co2"`
	Intruder       bool    `json:"intruder_alarm"`
	OutsideTemp    float64 `json:"outside_temp_sensor"`
	EnviroAirTemp1 float64 `json:"enviro_air_temp_1"`
	EnviroAirTemp2 float64 `json:"enviro_air_temp_2"`
	EnviroRH1      float64 `json:"enviro_rh_1"`
	EnviroRH2      float64 `json:"enviro_rh_2"`
	EnviroCO21     float64 `json:"enviro_co2_1"`
	EnviroCO22     float64 `json:"enviro_co2_2"`
	EnviroLight1   float64 `json:"enviro_light_1"`
	EnviroLight2   float64 `json:"enviro_light_2"`
}

// StatusIClimate represents the Status data structure from an IntelliClimate packet
type StatusIClimate struct {
	Readings         ReadingsIClimate         `json:"readings"`
	Statistics       StatisticsIClimate       `json:"statistics"`
	ModeAlarmHistory ModeAlarmHistoryIClimate `json:"mode_alarm_history"`
	SetPoints        []SetPointIClimate       `json:"set_points"`
	Status           []StatusStatusIClimate   `json:"status"`
}

// ReadingsIClimate represents the Readings data structure from an IntelliClimate packet
type ReadingsIClimate struct {
	AirTemp        AirTempIClimate        `json:"air_temp"`
	Detent         byte                   `json:"detent"`
	FailSafeAlarms FailSafeAlarmsIClimate `json:"fail_safe_alarms"`
	Light          LightIClimate          `json:"light"`
	PowerFail      PowerFailIClimate      `json:"power_fail"`
	Rh             RhIClimate             `json:"rh"`
	CO2            CO2IClimate            `json:"co2"`
	IntruderAlarm  IntruderAlarmIClimate  `json:"intruder"`
}

// IntruderAlarmIClimate represents the IntruderAlarm data structure from an IntelliClimate packet
type IntruderAlarmIClimate struct {
	Enabled bool `json:"enabled"`
	Page    bool `json:"page"`
}

// CO2IClimate represents the CO2 data structure from an IntelliClimate packet
type CO2IClimate struct {
	Target  float64 `json:"target"`
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	Page    bool    `json:"page"`
}

// AirTempIClimate represents the AirTemp data structure from an IntelliClimate packet
type AirTempIClimate struct {
	Cool    float64 `json:"cool"`
	Enabled bool    `json:"enabled"`
	Heat    float64 `json:"heat"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
	Page    bool    `json:"page"`
}

// FailSafeAlarmsIClimate represents the FailSafeAlarms data structure from an IntelliClimate packet
type FailSafeAlarmsIClimate struct {
	Enabled bool `json:"enabled"`
	Page    bool `json:"page"`
}

// LightIClimate represents the Light data structure from an IntelliClimate packet
type LightIClimate struct {
	Enabled bool    `json:"enabled"`
	Min     float64 `json:"min"`
	Page    bool    `json:"page"`
}

// PowerFailIClimate represents the PowerFail data structure from an IntelliClimate packet
type PowerFailIClimate struct {
	Enabled bool `json:"enabled"`
	Page    bool `json:"page"`
}

// RhIClimate represents the Rh data structure from an IntelliClimate packet
type RhIClimate struct {
	Enabled bool `json:"enabled"`
	Max     byte `json:"max"`
	Min     byte `json:"min"`
	Page    bool `json:"page"`
	Target  byte `json:"target"`
}

// StatisticsIClimate represents the Statistics data structure from an IntelliClimate packet
type StatisticsIClimate struct {
	Lights float64 `json:"lights"`
	CO2    float64 `json:"CO2"`
}

// ModeAlarmHistoryIClimate represents the ModeAlarmHistory data structure from an IntelliClimate packet
type ModeAlarmHistoryIClimate struct {
	Alarms []AlarmsIClimate `json:"alarms"`
	Mode   []ModeIClimate   `json:"mode"`
}

// AlarmsIClimate represents the Alarms data structure from an IntelliClimate packet
type AlarmsIClimate struct {
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// ModeIClimate represents the Mode data structure from an IntelliClimate packet
type ModeIClimate struct {
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// SetPointIClimate represents the SetPoint data structure from an IntelliClimate packet
type SetPointIClimate struct {
	LightBank     string  `json:"light_bank"`
	LightOn       int     `json:"light_on"`
	LightDuration int     `json:"light_duration"`
	DayTemp       float64 `json:"day_temp"`
	NightDropDeg  float64 `json:"night_drop_deg"`
	RhDay         int     `json:"rh_day"`
	RhMax         int     `json:"rh_max"`
	RhNight       int     `json:"rh_night"`
	CO2           int     `json:"co2"`
}

// StatusStatusIClimate represents the StatusStatus data structure from an IntelliClimate packet
type StatusStatusIClimate struct {
	Active    bool   `json:"active"`
	Enabled   bool   `json:"enabled"`
	ForceOn   bool   `json:"force_on"`
	Function  string `json:"function"`
	Installed bool   `json:"installed"`
}

// UnitsIClimate represents the Units data structure from an IntelliClimate packet
type UnitsIClimate struct {
	DateFormat  string `json:"date_format"`
	Temperature string `json:"temperature"`
}

// FunctionsIClimate represents the Functions data structure from an IntelliClimate packet
type FunctionsIClimate struct {
	Fan1                        bool   `json:"fan_1"`
	Fan2                        bool   `json:"fan_2"`
	AirConditioner              bool   `json:"air_conditioner"`
	Heater                      bool   `json:"heater"`
	Co2Sensor                   bool   `json:"co2_sensor"`
	Co2SensorRange              string `json:"co2_sensor_range"`
	Co2Injection                bool   `json:"co2_injection"`
	Co2Extraction               bool   `json:"co2_extraction"`
	Dehumidifier                bool   `json:"dehumidifier"`
	Humidifier                  bool   `json:"humidifier"`
	PulsedFogger                bool   `json:"pulsed_fogger"`
	LightBank1                  bool   `json:"light_bank_1"`
	LightsAirColored            bool   `json:"lights_air_colored"`
	LightBank2                  bool   `json:"light_bank_2"`
	LampOverTempShutdownSensors bool   `json:"lamp_over_temp_shutdown_sensors"`
	OutsideTempSensor           bool   `json:"outside_temp_sensor"`
	SecondEnviroSensor          bool   `json:"second_enviro_sensor"`
	IntruderAlarm               bool   `json:"intruder_alarm"`
	DehumidifyBy                string `json:"dehumidify_by"`
	Setup                       string `json:"setup"`
	MuteBuzzer                  bool   `json:"mute_buzzer"`
}

// AdvancedIClimate represents the Advanced data structure from an IntelliClimate packet
type AdvancedIClimate struct {
	ViewAdvancedSetting bool                     `json:"view_advanced_setting"`
	SwitchingOffsets    SwitchingOffsetsIClimate `json:"switching_offsets"`
	FailSafeSettings    FailSafeSettingsIClimate `json:"fail_safe_settings"`
	Rules               RulesIClimate            `json:"rules"`
}

// RulesIClimate represents the Rules data structure from an IntelliClimate packet
type RulesIClimate struct {
	HumidifyTempRules     HumidifyTempRulesIClimate     `json:"humidify_temp_rules"`
	MinimumAirChangeRules MinimumAirChangeRulesIClimate `json:"minimum_air_change_rules"`
	AllowAirCon           bool                          `json:"allow_air_con"`
	SetpointRamping       SetpointRampingIClimate       `json:"setpoint_ramping"`
	AirCon                AirConIClimate                `json:"air_con"`
	CO2Rules              CO2RulesIClimate              `json:"co2_rules"`
	Humidification        HumidificationIClimate        `json:"humidification"`
	Lighting              LightingIClimate              `json:"lighting"`
	FoggingRules          FoggingRulesIClimate          `json:"fogging_rules"`
	PurgingRules          PurgingRulesIClimate          `json:"purging_rules"`
}

// HumidifyTempRulesIClimate represents the HumidifyTempRules data structure from an IntelliClimate packet
type HumidifyTempRulesIClimate struct {
	LowerCoolingTemp float64 `json:"lower_cooling_temp"`
	RaiseHeatingTemp float64 `json:"raise_heating_temp"`
	RhLowThenRaise   float64 `json:"rh_low_then_raise"`
	PreventHeater    byte    `json:"prevent_heater"`
	HeatingOffset    float64 `json:"heating_offset"`
}

// MinimumAirChangeRulesIClimate represents the MinimumAirChangeRules data structure from an IntelliClimate packet
type MinimumAirChangeRulesIClimate struct {
	DaySecs        int `json:"day_secs"`
	EveryDayMins   int `json:"every_day_mins"`
	NightSecs      int `json:"night_secs"`
	EveryNightMins int `json:"every_night_mins"`
}

// SetpointRampingIClimate represents the SetpointRamping data structure from an IntelliClimate packet
type SetpointRampingIClimate struct {
	RampSetpoints byte `json:"ramp_setpoints"`
}

// AirConIClimate represents the AirCon data structure from an IntelliClimate packet
type AirConIClimate struct {
	ForceAirCon      bool    `json:"force_air_con"`
	AutoChangeAirCon float64 `json:"auto_change_air_con"`
	StartBefore      byte    `json:"start_before"`
	AutoStartAirCon  float64 `json:"auto_start_air_con"`
}

// CO2RulesIClimate represents the CO2Rules data structure from an IntelliClimate packet
type CO2RulesIClimate struct {
	Co2InjectionAllowed  bool    `json:"co2_injection_allowed"`
	InjectIfLightGreater float64 `json:"inject_if_light_greater"`
	Co2InjectionAvoid    bool    `json:"co2_injection_avoid"`
	Co2Cycling           float64 `json:"co2_cycling"`
	RiseVentTemp         float64 `json:"rise_vent_temp"`
	InjectTimeMin        byte    `json:"inject_time_min"`
	InjectTimeMax        byte    `json:"inject_time_max"`
	WaitTimeMin          byte    `json:"wait_time_min"`
	WaitTimeMax          byte    `json:"wait_time_max"`
	VentTimeMin          byte    `json:"vent_time_min"`
	VentTimeMax          byte    `json:"vent_time_max"`
}

// HumidificationIClimate represents the Humidification data structure from an IntelliClimate packet
type HumidificationIClimate struct {
	AllowHumidification  bool `json:"allow_humidification"`
	ChangeHumidification byte `json:"change_humidification"`
}

// LightingIClimate represents the Lighting data structure from an IntelliClimate packet
type LightingIClimate struct {
	LampCoolDownTime  byte `json:"lamp_cool_down_time"`
	SwOnNextLightBank byte `json:"sw_on_next_light_bank"`
}

// FoggingRulesIClimate represents the FoggingRules data structure from an IntelliClimate packet
type FoggingRulesIClimate struct {
	FogToCool      byte    `json:"fog_to_cool"`
	FogToAchieveRh float64 `json:"fog_to_achieve_rh"`
	FogTimes       int     `json:"fog_times"`
	FogTimeMax     byte    `json:"fog_time_max"`
	FogTimeMin     byte    `json:"fog_time_min"`
}

// PurgingRulesIClimate represents the PurgingRules data structure from an IntelliClimate packet
type PurgingRulesIClimate struct {
	PurgeMins byte `json:"purge_mins"`
	PurgeMin  byte `json:"purge_min"`
	PurgeMax  byte `json:"purge_max"`
}

// SwitchingOffsetsIClimate represents the SwitchingOffsets data structure from an IntelliClimate packet
type SwitchingOffsetsIClimate struct {
	AirConditionerOn  float64 `json:"air_conditioner_on"`
	AirConditionerOff float64 `json:"air_conditioner_off"`
	CO2On             float64 `json:"co2_on"`
	CO2Off            float64 `json:"co2_off"`
	DehumidifierOn    float64 `json:"dehumidifier_on"`
	DehumidifierOff   float64 `json:"dehumidifier_off"`
	FansOn            float64 `json:"fans_on"`
	FansOff           float64 `json:"fans_off"`
	HeaterOn          float64 `json:"heater_on"`
	HeaterOff         float64 `json:"heater_off"`
	HumidifierOn      float64 `json:"humidifier_on"`
	HumidifierOff     float64 `json:"humidifier_off"`
	PulsedFoggerOn    float64 `json:"pulsed_fogger_on"`
	PulsedFoggerOff   float64 `json:"pulsed_fogger_off"`
}

// FailSafeSettingsIClimate represents the FailSafeSettings data structure from an IntelliClimate packet
type FailSafeSettingsIClimate struct {
	FanFailOverride      FanFailOverrideIClimate      `json:"fan_fail_override"`
	AirConOverride       AirConOverrideIClimate       `json:"air_con_override"`
	DehumidifierOverride DehumidifierOverrideIClimate `json:"dehumidifier_override"`
	Co2FailSafe          Co2FailSafeIClimate          `json:"co2_fail_safe"`
	Co2InjectionOverride Co2InjectionOverrideIClimate `json:"co2_injection_override"`
	PowerFailure         PowerFailureIClimate         `json:"power_failure"`
	LightingOverride     LightingOverrideIClimate     `json:"light_falls_alarm_minimum"`
}

// FanFailOverrideIClimate represents the FanFailOverride data structure from an IntelliClimate packet
type FanFailOverrideIClimate struct {
	SwOffLightTempExceed  float64 `json:"sw_off_light_temp_exceed"`
	SwOffLightsTempExceed float64 `json:"sw_off_lights_temp_exceed"`
}

// AirConOverrideIClimate represents the AirConOverride data structure from an IntelliClimate packet
type AirConOverrideIClimate struct {
	SwAllExhaustFans float64 `json:"sw_all_exhaust_fans"`
}

// DehumidifierOverrideIClimate represents the DehumidifierOverride data structure from an IntelliClimate packet
type DehumidifierOverrideIClimate struct {
	SwOnFansRhExceed byte `json:"sw_on_fans_rh_exceed"`
	SwAcRhExceed     byte `json:"sw_ac_rh_exceed"`
}

// Co2FailSafeIClimate represents the Co2FailSafe data structure from an IntelliClimate packet
type Co2FailSafeIClimate struct {
	SwOnFansCo2Exceed int `json:"sw_on_fans_co2_exceed"`
}

// Co2InjectionOverrideIClimate represents the Co2InjectionOverride data structure from an IntelliClimate packet
type Co2InjectionOverrideIClimate struct {
	RevertFansCo2Falls int `json:"revert_fans_co2_falls"`
}

// PowerFailureIClimate represents the PowerFailure data structure from an IntelliClimate packet
type PowerFailureIClimate struct {
	SwLightsAfterCoolDown byte `json:"sw_lights_after_cool_down"`
}

// LightingOverrideIClimate represents the LightingOverride data structure from an IntelliClimate packet
type LightingOverrideIClimate struct {
	LightFallsAlarmMinimum bool `json:"light_falls_alarm_minimum"`
}

// GeneralIClimate represents the General data structure from an IntelliClimate packet
type GeneralIClimate struct {
	DeviceName string  `json:"device_name"`
	Firmware   float64 `json:"firmware"`
}

// IDose
type iDoseShadow struct {
	State StateIDose `json:"state"`
}

// StateIDose represents the State data structure from an IntelliDose packet
type StateIDose struct {
	Reported ReportedIDose `json:"reported"`
}

// ReportedIDose represents the Reported data structure from an IntelliDose packet
type ReportedIDose struct {
	Config    ConfigIDose  `json:"config"`
	Metrics   MetricsIDose `json:"metrics"`
	Status    StatusIDose  `json:"status"`
	Source    string       `json:"source"`
	Device    string       `json:"device"`
	Timestamp int64        `json:"timestamp"`
	Connected bool         `json:"connected"`
}

// ConfigIDose represents the Config data structure from an IntelliDose packet
type ConfigIDose struct {
	Units     UnitsIDose     `json:"units"`
	Times     TimesIDose     `json:"times"`
	Functions FunctionsIDose `json:"functions"`
	Advanced  AdvancedIDose  `json:"advanced"`
	General   GeneralIDose   `json:"general"`
}

// MetricsIDose represents the Metrics data structure from an IntelliDose packet
type MetricsIDose struct {
	Ec      float64 `json:"ec"`
	NutTemp float64 `json:"nut_temp"`
	PH      float64 `json:"pH"`
}

// StatusIDose represents the Status data structure from an IntelliDose packet
type StatusIDose struct {
	General   GeneralStatusIDose  `json:"general"`
	Nutrient  NutrientIDose       `json:"nutrient"`
	SetPoints SetPointsIDose      `json:"set_points"`
	Status    []StatusStatusIDose `json:"status"`
	Units     UnitsIDose          `json:"units"`
}

// GeneralStatusIDose represents the GeneralStatus data structure from an IntelliDose packet
type GeneralStatusIDose struct {
	DoseInterval        byte                    `json:"dose_interval"`
	NutrientDoseTime    byte                    `json:"nutrient_dose_time"`
	WaterOnTime         byte                    `json:"water_on_time"`
	IrrigationInterval1 IrrigationIntervalIDose `json:"irrigation_interval_1"`
	IrrigationInterval2 IrrigationIntervalIDose `json:"irrigation_interval_2"`
	IrrigationInterval3 IrrigationIntervalIDose `json:"irrigation_interval_3"`
	IrrigationInterval4 IrrigationIntervalIDose `json:"irrigation_interval_4"`
	IrrigationDuration1 int                     `json:"irrigation_duration_1"`
	IrrigationDuration2 int                     `json:"irrigation_duration_2"`
	IrrigationDuration3 int                     `json:"irrigation_duration_3"`
	IrrigationDuration4 int                     `json:"irrigation_duration_4"`
	MaxNutrientDoseTime byte                    `json:"max_nutrient_dose_time"`
	MaxPhDoseTime       byte                    `json:"max_ph_dose_time"`
	Mix1                byte                    `json:"mix_1"`
	Mix2                byte                    `json:"mix_2"`
	Mix3                byte                    `json:"mix_3"`
	Mix4                byte                    `json:"mix_4"`
	Mix5                byte                    `json:"mix_5"`
	Mix6                byte                    `json:"mix_6"`
	Mix7                byte                    `json:"mix_7"`
	Mix8                byte                    `json:"mix_8"`
	PhDoseTime          byte                    `json:"ph_dose_time"`
}

// IrrigationIntervalIDose represents the IrrigationInterval data structure from an IntelliDose packet
type IrrigationIntervalIDose struct {
	Day   int `json:"day"`
	Night int `json:"night"`
	Every int `json:"every"`
}

// NutrientIDose represents the Nutrient data structure from an IntelliDose packet
type NutrientIDose struct {
	Detent  byte         `json:"detent"`
	Ec      EcIDose      `json:"ec"`
	NutTemp NutTempIDose `json:"nut_temp"`
	Ph      PhIDose      `json:"ph"`
}

// EcIDose represents the Ec data structure from an IntelliDose packet
type EcIDose struct {
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

// NutTempIDose represents the NutTemp data structure from an IntelliDose packet
type NutTempIDose struct {
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

// PhIDose represents the Ph data structure from an IntelliDose packet
type PhIDose struct {
	Enabled bool    `json:"enabled"`
	Max     float64 `json:"max"`
	Min     float64 `json:"min"`
}

// SetPointsIDose represents the SetPoints data structure from an IntelliDose packet
type SetPointsIDose struct {
	Nutrient      float64 `json:"nutrient"`
	NutrientNight float64 `json:"nutrient_night"`
	PhDosing      string  `json:"ph_dosing"`
	Ph            float64 `json:"ph"`
}

// StatusStatusIDose represents the StatusStatus data structure from an IntelliDose packet
type StatusStatusIDose struct {
	Active   bool   `json:"active"`
	Enabled  bool   `json:"enabled"`
	ForceOn  bool   `json:"force_on"`
	Function string `json:"function"`
}

// UnitsIDose represents the Units data structure from an IntelliDose packet
type UnitsIDose struct {
	DateFormat              string `json:"date_format"`
	Temperature             string `json:"temperature"`
	Ec                      string `json:"ec"`
	TdsConversationStandart int    `json:"tds_conversation_standart"`
}

// TimesIDose represents the Times data structure from an IntelliDose packet
type TimesIDose struct {
	DayStart int `json:"day_start"`
	DayEnd   int `json:"day_end"`
}

// FunctionsIDose represents the Functions data structure from an IntelliDose packet
type FunctionsIDose struct {
	NutrientsParts     byte   `json:"nutrients_parts"`
	PhDosing           string `json:"ph_dosing"`
	IrrigationMode     string `json:"irrigation_mode"`
	IrrigationStations byte   `json:"irrigation_stations"`
	SeparatePumpOutput bool   `json:"separate_pump_output"`
	UseWater           bool   `json:"use_water"`
	ExternalAlarm      bool   `json:"external_alarm"`
	DayNightEc         bool   `json:"day_night_ec"`
	IrrigationStation1 string `json:"irrigation_station_1"`
	IrrigationStation2 string `json:"irrigation_station_2"`
	IrrigationStation3 string `json:"irrigation_station_3"`
	IrrigationStation4 string `json:"irrigation_station_4"`
	Scheduling         bool   `json:"scheduling"`
	MuteBuzzer         bool   `json:"mute_buzzer"`
}

// AdvancedIDose represents the Advanced data structure from an IntelliDose packet
type AdvancedIDose struct {
	ProportinalDosing bool   `json:"proportinal_dosing"`
	SequentialDosing  bool   `json:"sequential_dosing"`
	DisableEc         bool   `json:"disable_ec"`
	DisablePh         bool   `json:"disable_ph"`
	MntnReminderFreq  string `json:"mntn_reminder_freq"`
}

// GeneralIDose represents the General data structure from an IntelliDose packet
type GeneralIDose struct {
	DeviceName string  `json:"device_name"`
	Firmware   float64 `json:"firmware"`
}

func (device *Device) updateShadow() {
	device.updating.Lock()
	defer device.updating.Unlock()

	var currentState interface{}
	err := device.updateState()
	if err != nil && !device.checkStates() {
		device.IsOpen = false
		tell.Errorf("failed to update device state: %s", err)
		return
	}

	timestamp := time.Now().Unix()
	switch device.DeviceType {
	case IntelliDoseDeviceType:
		currentState = parseByteResponseForIDose(device.states.d0State, device.states.d1State, device.states.d2State, device.SerialNumber, timestamp)
	case IntelliClimateDeviceType:
		currentState = parseByteResponseForIClimate(device.states.d0State, device.states.d1State, device.states.d2State, device.states.d3State, device.SerialNumber, timestamp)
	}

	device.update(currentState)
}

func (device Device) checkStates() bool {
	if len(device.states.d0State) != requestLength || len(device.states.d1State) != requestLength || len(device.states.d2State) != requestLength {
		return false
	}
	if device.DeviceType == IntelliClimateDeviceType {
		if len(device.states.d3State) != requestLength {
			return false
		}
	}
	return true
}

func (device Device) updateState() error {
	device.readWriteLock.Lock()
	defer device.readWriteLock.Unlock()
	var err error
	deviceType := device.DeviceType
	device.states.d0State, err = device.sentRequest(d0Request)
	if err != nil {
		return err
	}
	tell.Debugf(fmt.Sprintf("%s D0: % x", device.SerialNumber, device.states.d0State))
	device.states.d1State, err = device.sentRequest(d1Request)
	if err != nil {
		return err
	}
	tell.Debugf(fmt.Sprintf("%s D1: % x", device.SerialNumber, device.states.d1State))
	device.states.d2State, err = device.sentRequest(d2Request)
	if err != nil {
		return err
	}
	tell.Debugf(fmt.Sprintf("%s D2: % x", device.SerialNumber, device.states.d2State))
	if deviceType == IntelliClimateDeviceType {
		device.states.d3State, err = device.sentRequest(d3Request)
		if err != nil {
			return err
		}
		tell.Debugf(fmt.Sprintf("%s D3: % x", device.SerialNumber, device.states.d3State))
	}
	return nil
}

func (device Device) writeDoseData(state iDoseShadow) {
	device.readWriteLock.Lock()
	defer device.readWriteLock.Unlock()
	d0State := device.states.d0State

	s0Request := device.states.d1State
	prepareS0RequestForIDose(&d0State, &s0Request, state)
	s0Request = append([]byte{0x00}, s0Request...)
	device.sentRequest(s0Request)

	tell.Debugf(fmt.Sprintf("%s S0 iDose request: % x", device.SerialNumber, s0Request))

	s1Request := device.states.d2State
	prepareS1RequestForIDose(&d0State, &s1Request, state)
	s1Request = append([]byte{0x00}, s1Request...)
	device.sentRequest(s1Request)

	tell.Debugf(fmt.Sprintf("%s S1 iDose request: % x", device.SerialNumber, s1Request))
}

func prepareS0RequestForIDose(d0bytes *[]byte, bytes *[]byte, doseShadow iDoseShadow) {
	nutTempMax := int(doseShadow.State.Reported.Status.Nutrient.NutTemp.Max * 100)
	nutTempMin := int(doseShadow.State.Reported.Status.Nutrient.NutTemp.Min * 100)
	nutTempEnabled := doseShadow.State.Reported.Status.Nutrient.NutTemp.Enabled
	pHMax := int(doseShadow.State.Reported.Status.Nutrient.Ph.Max * 10)
	pHMin := int(doseShadow.State.Reported.Status.Nutrient.Ph.Min * 10)
	pHEnabled := doseShadow.State.Reported.Status.Nutrient.Ph.Enabled
	eCMax := int(doseShadow.State.Reported.Status.Nutrient.Ec.Max / 10)
	eCMin := int(doseShadow.State.Reported.Status.Nutrient.Ec.Min / 10)
	eCEnabled := doseShadow.State.Reported.Status.Nutrient.Ec.Enabled
	nutrient := int(doseShadow.State.Reported.Status.SetPoints.Nutrient)
	nutrientNight := int(doseShadow.State.Reported.Status.SetPoints.NutrientNight)
	phMode := getBoolFromPHMode(doseShadow.State.Reported.Status.SetPoints.PhDosing)
	ph := doseShadow.State.Reported.Status.SetPoints.Ph * 10
	nutrientDosingStatus := getStatusIDoseFunctionByName(doseShadow.State.Reported.Status.Status, nutrientDosingFunction)
	pHStatus := getStatusIDoseFunctionByName(doseShadow.State.Reported.Status.Status, phFunction)
	phDosingHi, phDosingLo := getPhDosingBoolByString(doseShadow.State.Reported.Config.Functions.PhDosing)
	dayNightEc := doseShadow.State.Reported.Config.Functions.DayNightEc
	mutebuzzer := doseShadow.State.Reported.Config.Functions.MuteBuzzer
	separatePumpOutput := doseShadow.State.Reported.Config.Functions.SeparatePumpOutput

	irrigation0Status := getStatusIDoseFunctionByName(doseShadow.State.Reported.Status.Status, irrigationStation1Function)
	irrigation1Status := getStatusIDoseFunctionByName(doseShadow.State.Reported.Status.Status, irrigationStation2Function)
	irrigation2Status := getStatusIDoseFunctionByName(doseShadow.State.Reported.Status.Status, irrigationStation3Function)
	irrigation3Status := getStatusIDoseFunctionByName(doseShadow.State.Reported.Status.Status, irrigationStation4Function)

	waterStatus := getStatusIDoseFunctionByName(doseShadow.State.Reported.Status.Status, waterFunction)

	irrigationMode := doseShadow.State.Reported.Config.Functions.IrrigationMode

	dateFormat := doseShadow.State.Reported.Config.Units.DateFormat
	temperature := doseShadow.State.Reported.Config.Units.Temperature
	nutrientConfigH, nutrientConfigL := isNutrientConfig(doseShadow.State.Reported.Config.Units.Ec)

	var irrigationIndepended bool
	var irrigationSequential bool
	if irrigationMode == irrigationModeIndependent {
		irrigationIndepended = true
		irrigationSequential = false
	} else if irrigationMode == irrigationModeSequential {
		irrigationIndepended = true
		irrigationSequential = true
	} else {
		irrigationIndepended = false
		irrigationSequential = false
	}

	if irrigationIndepended == irrigationSequential {
		irrigation0Status = getStatusIDoseFunctionByName(doseShadow.State.Reported.Status.Status, irrigationFunction)
	}

	irrigateTOD0, irrigateDayOnly0 := getIrrigationStationConfigurationByString(doseShadow.State.Reported.Config.Functions.IrrigationStation1)
	irrigateTOD1, irrigateDayOnly1 := getIrrigationStationConfigurationByString(doseShadow.State.Reported.Config.Functions.IrrigationStation2)
	irrigateTOD2, irrigateDayOnly2 := getIrrigationStationConfigurationByString(doseShadow.State.Reported.Config.Functions.IrrigationStation3)
	irrigateTOD3, irrigateDayOnly3 := getIrrigationStationConfigurationByString(doseShadow.State.Reported.Config.Functions.IrrigationStation4)

	useWater := doseShadow.State.Reported.Config.Functions.UseWater

	irrigationDuration := doseShadow.State.Reported.Status.General.IrrigationDuration1
	irrigationDay := doseShadow.State.Reported.Status.General.IrrigationInterval1.Day
	irrigationNight := doseShadow.State.Reported.Status.General.IrrigationInterval1.Night
	irrigationEvery := doseShadow.State.Reported.Status.General.IrrigationInterval1.Every

	tdsConversationStandart := doseShadow.State.Reported.Config.Units.TdsConversationStandart

	irrigationStations := doseShadow.State.Reported.Config.Functions.IrrigationStations
	irrigationInstalled := false
	if irrigationStations > 0 {
		irrigationInstalled = true
	}

	(*bytes)[0] = 0x53
	(*bytes)[1] = 0x30
	(*bytes)[2] = updateByte((*bytes)[2], doseShadow.State.Reported.Config.Advanced.ProportinalDosing, 0)
	(*bytes)[2] = updateByte((*bytes)[2], doseShadow.State.Reported.Config.Advanced.SequentialDosing, 1)
	(*bytes)[2] = updateByte((*bytes)[2], doseShadow.State.Reported.Config.Functions.ExternalAlarm, 2)
	(*bytes)[2] = updateByte((*bytes)[2], phMode, 3)
	(*bytes)[2] = updateByte((*bytes)[2], nutrientDosingStatus.Enabled, 4)
	(*bytes)[2] = updateByte((*bytes)[2], irrigation0Status.Enabled, 5)
	(*bytes)[2] = updateByte((*bytes)[2], waterStatus.Enabled, 6)
	(*bytes)[2] = updateByte((*bytes)[2], irrigateDayOnly0, 7)
	(*bytes)[3] = updateByte((*bytes)[3], phDosingLo, 0)
	(*bytes)[3] = updateByte((*bytes)[3], phDosingHi, 1)
	(*bytes)[3] = updateByte((*bytes)[3], dayNightEc, 2)
	(*bytes)[3] = updateByte((*bytes)[3], irrigateTOD0, 3)
	(*bytes)[3] = updateByte((*bytes)[3], irrigation0Status.ForceOn, 4)
	(*bytes)[3] = updateByte((*bytes)[3], waterStatus.ForceOn, 5)
	(*bytes)[3] = updateByte((*bytes)[3], pHStatus.ForceOn, 6)
	(*bytes)[3] = updateByte((*bytes)[3], nutrientDosingStatus.ForceOn, 7)
	(*bytes)[4] = updateByte((*bytes)[4], useWater, 2)
	(*bytes)[4] = updateByte((*bytes)[4], irrigationInstalled, 3)
	(*bytes)[4] = updateByte((*bytes)[4], isUSADateFormat(dateFormat), 4)
	(*bytes)[4] = updateByte((*bytes)[4], nutrientConfigH, 5)
	(*bytes)[4] = updateByte((*bytes)[4], nutrientConfigL, 6)
	(*bytes)[4] = updateByte((*bytes)[4], isTemperatureF(temperature), 7)
	(*bytes)[5] = updateByte((*bytes)[5], pHStatus.Enabled, 0)
	(*bytes)[5] = updateByte((*bytes)[5], separatePumpOutput, 1)
	(*bytes)[5] = updateByte((*bytes)[5], doseShadow.State.Reported.Config.Advanced.DisablePh, 2)
	(*bytes)[5] = updateByte((*bytes)[5], doseShadow.State.Reported.Config.Advanced.DisableEc, 3)
	(*bytes)[5] = updateByte((*bytes)[5], mutebuzzer, 4)
	(*bytes)[5] = updateByte((*bytes)[5], nutTempEnabled, 5)
	(*bytes)[5] = updateByte((*bytes)[5], pHEnabled, 6)
	(*bytes)[5] = updateByte((*bytes)[5], eCEnabled, 7)
	(*bytes)[6] = doseShadow.State.Reported.Status.Nutrient.Detent
	(*bytes)[7] = byte(eCMax)
	(*bytes)[8] = byte(eCMin)
	(*bytes)[9] = byte(pHMax)
	(*bytes)[10] = byte(pHMin)
	(*bytes)[11] = byte(nutTempMax)
	(*bytes)[12] = byte(nutTempMax >> 8)
	(*bytes)[13] = byte(nutTempMin)
	(*bytes)[14] = byte(nutTempMin >> 8)
	(*bytes)[15] = byte(nutrient)
	(*bytes)[16] = byte(nutrient >> 8)
	(*bytes)[17] = byte(nutrientNight)
	(*bytes)[18] = byte(nutrientNight >> 8)
	(*bytes)[19] = byte(ph)
	(*bytes)[20] = doseShadow.State.Reported.Status.General.MaxNutrientDoseTime
	(*bytes)[21] = doseShadow.State.Reported.Status.General.NutrientDoseTime
	(*bytes)[22] = doseShadow.State.Reported.Status.General.Mix1
	(*bytes)[23] = doseShadow.State.Reported.Status.General.Mix2
	(*bytes)[24] = doseShadow.State.Reported.Status.General.Mix3
	(*bytes)[25] = doseShadow.State.Reported.Status.General.MaxPhDoseTime
	(*bytes)[26] = doseShadow.State.Reported.Status.General.PhDoseTime
	(*bytes)[28] = doseShadow.State.Reported.Status.General.DoseInterval
	(*bytes)[29] = byte(irrigationDuration)
	(*bytes)[30] = byte(irrigationDuration >> 8)
	(*bytes)[31] = byte(irrigationDay)
	(*bytes)[32] = byte(irrigationDay >> 8)
	(*bytes)[33] = byte(irrigationNight)
	(*bytes)[34] = byte(irrigationNight >> 8)
	(*bytes)[35] = byte(irrigationEvery)
	(*bytes)[36] = byte(irrigationEvery >> 8)
	(*bytes)[37] = doseShadow.State.Reported.Status.General.WaterOnTime
	(*bytes)[38] = byte(doseShadow.State.Reported.Config.Times.DayStart / 6)
	(*bytes)[39] = byte(doseShadow.State.Reported.Config.Times.DayEnd / 6)
	(*bytes)[40] = doseShadow.State.Reported.Config.Functions.NutrientsParts
	(*bytes)[48] = doseShadow.State.Reported.Status.General.Mix4
	(*bytes)[49] = doseShadow.State.Reported.Status.General.Mix5
	(*bytes)[50] = doseShadow.State.Reported.Status.General.Mix6
	(*bytes)[51] = doseShadow.State.Reported.Status.General.Mix7
	(*bytes)[52] = doseShadow.State.Reported.Status.General.Mix8
	(*bytes)[53] = byte(tdsConversationStandart)
	(*bytes)[54] = byte(tdsConversationStandart >> 8)
	(*bytes)[57] = updateByte((*bytes)[57], irrigation1Status.Enabled, 0)
	(*bytes)[57] = updateByte((*bytes)[57], irrigateDayOnly1, 2)
	(*bytes)[57] = updateByte((*bytes)[57], irrigateTOD1, 3)
	(*bytes)[57] = updateByte((*bytes)[57], irrigation1Status.ForceOn, 4)
	(*bytes)[57] = updateByte((*bytes)[57], irrigation1Status.ForceOn, 4)
	(*bytes)[58] = updateByte((*bytes)[58], irrigation2Status.Enabled, 0)
	(*bytes)[58] = updateByte((*bytes)[58], irrigateDayOnly2, 2)
	(*bytes)[58] = updateByte((*bytes)[58], irrigateTOD2, 3)
	(*bytes)[58] = updateByte((*bytes)[58], irrigation2Status.ForceOn, 4)
	(*bytes)[58] = updateByte((*bytes)[58], irrigationSequential, 7)
	(*bytes)[59] = updateByte((*bytes)[59], irrigation3Status.Enabled, 0)
	(*bytes)[59] = updateByte((*bytes)[59], irrigateDayOnly3, 2)
	(*bytes)[59] = updateByte((*bytes)[59], irrigateTOD3, 3)
	(*bytes)[59] = updateByte((*bytes)[59], irrigation3Status.ForceOn, 4)
	(*bytes)[59] = updateByte((*bytes)[59], irrigationIndepended, 7)
	(*bytes)[60] = irrigationStations
	(*bytes)[61] = (*d0bytes)[61]

	createCheckSum(bytes)
}

func prepareS1RequestForIDose(d0bytes *[]byte, bytes *[]byte, doseShadow iDoseShadow) {
	irrigationDuration1 := doseShadow.State.Reported.Status.General.IrrigationDuration2
	irrigationDuration2 := doseShadow.State.Reported.Status.General.IrrigationDuration3
	irrigationDuration3 := doseShadow.State.Reported.Status.General.IrrigationDuration4

	irrigationDay1 := doseShadow.State.Reported.Status.General.IrrigationInterval2.Day
	irrigationNight1 := doseShadow.State.Reported.Status.General.IrrigationInterval2.Night
	irrigationEvery1 := doseShadow.State.Reported.Status.General.IrrigationInterval2.Every
	irrigationDay2 := doseShadow.State.Reported.Status.General.IrrigationInterval3.Day
	irrigationNight2 := doseShadow.State.Reported.Status.General.IrrigationInterval3.Night
	irrigationEvery2 := doseShadow.State.Reported.Status.General.IrrigationInterval3.Every
	irrigationDay3 := doseShadow.State.Reported.Status.General.IrrigationInterval4.Day
	irrigationNight3 := doseShadow.State.Reported.Status.General.IrrigationInterval4.Night
	irrigationEvery3 := doseShadow.State.Reported.Status.General.IrrigationInterval4.Every

	deviceName := doseShadow.State.Reported.Config.General.DeviceName
	var deviceNameByteArray [10]byte
	copy(deviceNameByteArray[:], deviceName)

	(*bytes)[0] = 0x53
	(*bytes)[1] = 0x31
	(*bytes)[2] = deviceNameByteArray[0]
	(*bytes)[3] = deviceNameByteArray[1]
	(*bytes)[4] = deviceNameByteArray[2]
	(*bytes)[5] = deviceNameByteArray[3]
	(*bytes)[6] = deviceNameByteArray[4]
	(*bytes)[7] = deviceNameByteArray[5]
	(*bytes)[8] = deviceNameByteArray[6]
	(*bytes)[9] = deviceNameByteArray[7]
	(*bytes)[10] = deviceNameByteArray[8]
	(*bytes)[11] = deviceNameByteArray[9]
	(*bytes)[21] = byte(irrigationDuration1)
	(*bytes)[22] = byte(irrigationDuration1 >> 8)
	(*bytes)[23] = byte(irrigationDay1)
	(*bytes)[24] = byte(irrigationDay1 >> 8)
	(*bytes)[25] = byte(irrigationNight1)
	(*bytes)[26] = byte(irrigationNight1 >> 8)
	(*bytes)[27] = byte(irrigationEvery1)
	(*bytes)[28] = byte(irrigationEvery1 >> 8)
	(*bytes)[29] = byte(irrigationDuration2)
	(*bytes)[30] = byte(irrigationDuration2 >> 8)
	(*bytes)[31] = byte(irrigationDay2)
	(*bytes)[32] = byte(irrigationDay2 >> 8)
	(*bytes)[33] = byte(irrigationNight2)
	(*bytes)[34] = byte(irrigationNight2 >> 8)
	(*bytes)[35] = byte(irrigationEvery2)
	(*bytes)[36] = byte(irrigationEvery2 >> 8)
	(*bytes)[37] = byte(irrigationDuration3)
	(*bytes)[38] = byte(irrigationDuration3 >> 8)
	(*bytes)[39] = byte(irrigationDay3)
	(*bytes)[40] = byte(irrigationDay3 >> 8)
	(*bytes)[41] = byte(irrigationNight3)
	(*bytes)[42] = byte(irrigationNight3 >> 8)
	(*bytes)[43] = byte(irrigationEvery3)
	(*bytes)[44] = byte(irrigationEvery3 >> 8)
	(*bytes)[61] = (*d0bytes)[61]
	createCheckSum(bytes)
}

func (device Device) writeClimateData(state iClimateShadow) {
	device.readWriteLock.Lock()
	defer device.readWriteLock.Unlock()
	d0State := device.states.d0State

	s0Request := device.states.d1State
	prepareS0RequestForIClimate(&d0State, &s0Request, state)
	s0Request = append([]byte{0x00}, s0Request...)
	device.sentRequest(s0Request)

	tell.Debugf(fmt.Sprintf("%s S0 iClimat request: % x", device.SerialNumber, s0Request))

	s1Request := device.states.d2State
	prepareS1RequestForIClimate(&d0State, &s1Request, state)
	s1Request = append([]byte{0x00}, s1Request...)
	device.sentRequest(s1Request)

	tell.Debugf(fmt.Sprintf("%s S1 iClimat request: % x", device.SerialNumber, s1Request))

	s2Request := device.states.d3State
	prepareS2RequestForIClimate(&d0State, &s2Request, state)
	s2Request = append([]byte{0x00}, s2Request...)
	device.sentRequest(s2Request)

	tell.Debugf(fmt.Sprintf("%s S2 iClimat request: % x", device.SerialNumber, s2Request))
}

func prepareS0RequestForIClimate(d0bytes *[]byte, bytes *[]byte, climateShadow iClimateShadow) {
	tempCool := int(climateShadow.State.Reported.Status.Readings.AirTemp.Cool * 100)
	tempHeat := int(climateShadow.State.Reported.Status.Readings.AirTemp.Heat * 100)
	tempMin := int(climateShadow.State.Reported.Status.Readings.AirTemp.Min * 100)
	tempMax := int(climateShadow.State.Reported.Status.Readings.AirTemp.Max * 100)
	tempAlarmEnabled := climateShadow.State.Reported.Status.Readings.AirTemp.Enabled
	rHAlarmEnabled := climateShadow.State.Reported.Status.Readings.Rh.Enabled
	lightAlarmEnabled := climateShadow.State.Reported.Status.Readings.Light.Enabled
	lightMin := int(climateShadow.State.Reported.Status.Readings.Light.Min)
	powerFailAlarmEnabled := climateShadow.State.Reported.Status.Readings.PowerFail.Enabled
	failSafeAlarmEnabled := climateShadow.State.Reported.Status.Readings.FailSafeAlarms.Enabled
	co2AlarmEnabled := climateShadow.State.Reported.Status.Readings.CO2.Enabled
	co2Max := int(climateShadow.State.Reported.Status.Readings.CO2.Max / 25)
	co2Min := int(climateShadow.State.Reported.Status.Readings.CO2.Min / 25)
	co2Target := int(climateShadow.State.Reported.Status.Readings.CO2.Target)
	intruderAlarmEanbled := climateShadow.State.Reported.Status.Readings.IntruderAlarm.Enabled
	buzzerMuted := climateShadow.State.Reported.Config.Functions.MuteBuzzer

	co2SensorConfigEnabled := climateShadow.State.Reported.Config.Functions.Co2Sensor
	fan1ConfigEnabled := climateShadow.State.Reported.Config.Functions.Fan1
	fan2ConfigEnabled := climateShadow.State.Reported.Config.Functions.Fan2
	airConditionerConfigEnabled := climateShadow.State.Reported.Config.Functions.AirConditioner
	co2InjectionConfigEnabled := climateShadow.State.Reported.Config.Functions.Co2Injection
	heaterConfigEnabled := climateShadow.State.Reported.Config.Functions.Heater
	dehumidifierConfigEnabled := climateShadow.State.Reported.Config.Functions.Dehumidifier
	humidifierConfigEnabled := climateShadow.State.Reported.Config.Functions.Humidifier
	light1ConfigEnabled := climateShadow.State.Reported.Config.Functions.LightBank1
	light2ConfigEnabled := climateShadow.State.Reported.Config.Functions.LightBank2
	lightsAirConfigEnabled := climateShadow.State.Reported.Config.Functions.LightsAirColored
	foggerConfigEnabled := climateShadow.State.Reported.Config.Functions.PulsedFogger
	intruderAlarmConfigEnabled := climateShadow.State.Reported.Config.Functions.IntruderAlarm
	outsideTempConfigEnabled := climateShadow.State.Reported.Config.Functions.OutsideTempSensor
	secondEnviroSensorConfigEnabled := climateShadow.State.Reported.Config.Functions.SecondEnviroSensor
	lampOverTempShutdownSensorsConfigEnabled := climateShadow.State.Reported.Config.Functions.LampOverTempShutdownSensors

	dehumidifyBy := climateShadow.State.Reported.Config.Functions.DehumidifyBy

	purge, ac := isDehumidifyBy(dehumidifyBy)

	fan1Status := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, fan1Function)
	fan2Status := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, fan2Function)
	airConditionerStatus := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, airConFunction)
	co2ExtractStatus := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, co2ExtractFunction)
	co2InjectStatus := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, co2InjectionFunction)
	heaterStatus := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, heaterFunction)
	dehumidifierStatus := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, dehumidifierFunction)
	humidifierStatus := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, humidifierFunction)
	light1Status := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, lightBank1Function)
	light2Status := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, lightBank2Function)
	foggerStatus := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, foggerFunction)
	purgingStatus := getStatusIClimateFunctionByName(climateShadow.State.Reported.Status.Status, purgingFunction)
	dateFormat := climateShadow.State.Reported.Config.Units.DateFormat
	temperature := climateShadow.State.Reported.Config.Units.Temperature

	tell.Infof("D1 SetBool5: %s", encoding.ByteToBitString((*bytes)[6]))

	deviceName := climateShadow.State.Reported.Config.General.DeviceName
	var deviceNameByteArray [10]byte
	copy(deviceNameByteArray[:], deviceName)

	(*bytes)[0], (*bytes)[1] = byte('S'), byte('0')

	// slice to have correct index
	specificationBytes := (*bytes)[2:]

	specificationBytes[0] = encoding.ComposeByte(fan1ConfigEnabled, fan2ConfigEnabled, airConditionerConfigEnabled, heaterConfigEnabled,
		co2SensorConfigEnabled, co2InjectionConfigEnabled, foggerConfigEnabled, humidifierConfigEnabled)

	specificationBytes[1] = encoding.ComposeByte(dehumidifierConfigEnabled, light1ConfigEnabled, light2ConfigEnabled, outsideTempConfigEnabled,
		secondEnviroSensorConfigEnabled, intruderAlarmConfigEnabled, false, lightsAirConfigEnabled)

	specificationBytes[2] = encoding.ComposeByte(isTemperatureF(temperature), isUSADateFormat(dateFormat), false, // control by VPD never used by specs for d1
		failSafeAlarmEnabled, false, buzzerMuted, ac, purge)

	var co2InjectionEnabledComposite bool
	if co2InjectionConfigEnabled {
		co2InjectionEnabledComposite = co2InjectStatus.Enabled
	} else {
		co2InjectionEnabledComposite = co2ExtractStatus.Enabled
	}
	specificationBytes[3] = encoding.ComposeByte(tempAlarmEnabled, rHAlarmEnabled, co2AlarmEnabled, lightAlarmEnabled,
		intruderAlarmEanbled, powerFailAlarmEnabled, co2InjectionEnabledComposite, lampOverTempShutdownSensorsConfigEnabled)

	var fan2On, fun2En bool
	if fan2ConfigEnabled {
		fan2On, fun2En = fan2Status.ForceOn, fan2Status.Enabled
	} else {
		fan2On, fun2En = airConditionerStatus.ForceOn, airConditionerStatus.Enabled
	}
	specificationBytes[4] = encoding.ComposeByte(fan1Status.Enabled, fan1Status.ForceOn, fun2En, fan2On,
		heaterStatus.Enabled, heaterStatus.ForceOn, dehumidifierStatus.Enabled, dehumidifierStatus.ForceOn)

	var humidOn, humidEn bool
	if humidifierConfigEnabled {
		humidOn, humidEn = humidifierStatus.ForceOn, humidifierStatus.Enabled
	} else {
		humidOn, humidEn = foggerStatus.ForceOn, foggerStatus.Enabled
	}
	specificationBytes[5] = encoding.ComposeByte(humidEn, humidOn, light1Status.Enabled, light1Status.ForceOn,
		light2Status.Enabled, light2Status.ForceOn, purgingStatus.Enabled, purgingStatus.ForceOn)

	specificationBytes[7] = climateShadow.State.Reported.Status.Readings.Detent
	specificationBytes[8], specificationBytes[9] = encoding.UnsignedIntToBytes(tempCool)
	specificationBytes[10], specificationBytes[11] = encoding.UnsignedIntToBytes(tempMin)
	specificationBytes[12], specificationBytes[13] = encoding.UnsignedIntToBytes(tempMax)
	specificationBytes[14] = climateShadow.State.Reported.Status.Readings.Rh.Target
	specificationBytes[15] = climateShadow.State.Reported.Status.Readings.Rh.Min
	specificationBytes[16] = climateShadow.State.Reported.Status.Readings.Rh.Max
	specificationBytes[23], specificationBytes[24] = encoding.UnsignedIntToBytes(co2Target)
	specificationBytes[25] = byte(co2Min)
	specificationBytes[26] = byte(co2Max)
	specificationBytes[27], specificationBytes[28] = encoding.UnsignedIntToBytes(lightMin)
	for i := range deviceNameByteArray {
		specificationBytes[31+i] = deviceNameByteArray[i]
	}
	fillSetPointData(&specificationBytes, climateShadow.State.Reported.Status.SetPoints[0])
	specificationBytes[57], specificationBytes[58] = encoding.UnsignedIntToBytes(tempHeat)
	specificationBytes[59] = (*d0bytes)[61]

	createCheckSum(bytes)
}

// extracts and set SetPointIClimate data to array
func fillSetPointData(bytes *[]byte, iClimate SetPointIClimate) {
	(*bytes)[41] = byte(getLightBoxModeToInt(iClimate.LightBank))
	(*bytes)[42], (*bytes)[43] = encoding.UnsignedIntToBytes(iClimate.LightOn)
	(*bytes)[44], (*bytes)[45] = encoding.UnsignedIntToBytes(iClimate.LightDuration)
	(*bytes)[46], (*bytes)[47] = encoding.UnsignedIntToBytes(int(iClimate.DayTemp * 100))
	(*bytes)[48], (*bytes)[49] = encoding.UnsignedIntToBytes(int(iClimate.NightDropDeg * 100))
	(*bytes)[50], (*bytes)[51] = encoding.UnsignedIntToBytes(iClimate.RhDay)
	(*bytes)[52], (*bytes)[53] = encoding.UnsignedIntToBytes(iClimate.RhNight)
	(*bytes)[54] = byte(iClimate.RhMax)
	(*bytes)[55], (*bytes)[56] = encoding.UnsignedIntToBytes(int(iClimate.CO2))
}

func prepareS1RequestForIClimate(d0bytes *[]byte, bytes *[]byte, climateShadow iClimateShadow) {
	co2ExtractionEnabled := climateShadow.State.Reported.Config.Functions.Co2Extraction

	daySecs := climateShadow.State.Reported.Config.Advanced.Rules.MinimumAirChangeRules.DaySecs
	everyDayMins := climateShadow.State.Reported.Config.Advanced.Rules.MinimumAirChangeRules.EveryDayMins
	nightSecs := climateShadow.State.Reported.Config.Advanced.Rules.MinimumAirChangeRules.NightSecs
	everyNightMins := climateShadow.State.Reported.Config.Advanced.Rules.MinimumAirChangeRules.EveryNightMins
	swOffLightTempExceed := int(climateShadow.State.Reported.Config.Advanced.FailSafeSettings.FanFailOverride.SwOffLightTempExceed * 100)
	swOffLightsTempExceed := int(climateShadow.State.Reported.Config.Advanced.FailSafeSettings.FanFailOverride.SwOffLightsTempExceed * 100)
	swOnFansCo2Exceed := climateShadow.State.Reported.Config.Advanced.FailSafeSettings.Co2FailSafe.SwOnFansCo2Exceed
	revertFansCo2Falls := climateShadow.State.Reported.Config.Advanced.FailSafeSettings.Co2InjectionOverride.RevertFansCo2Falls
	swAllExhaustFans := int(climateShadow.State.Reported.Config.Advanced.FailSafeSettings.AirConOverride.SwAllExhaustFans)
	lowerCoolingTemp := int(climateShadow.State.Reported.Config.Advanced.Rules.HumidifyTempRules.LowerCoolingTemp * 100)
	raiseHeatingTemp := int(climateShadow.State.Reported.Config.Advanced.Rules.HumidifyTempRules.RaiseHeatingTemp * 100)
	rhLowThenRaise := int(climateShadow.State.Reported.Config.Advanced.Rules.HumidifyTempRules.RhLowThenRaise * 100)
	autoChangeAirCon := int(climateShadow.State.Reported.Config.Advanced.Rules.AirCon.AutoChangeAirCon * 100)
	injectIfLightGreater := int(climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.InjectIfLightGreater)
	co2Cycling := int(climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.Co2Cycling * 100)
	riseVentTemp := int(climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.RiseVentTemp * 100)
	fogToAchieveRh := int(climateShadow.State.Reported.Config.Advanced.Rules.FoggingRules.FogToAchieveRh * 100)
	fogTimes := int(climateShadow.State.Reported.Config.Advanced.Rules.FoggingRules.FogTimes)
	autoStartAirCon := int(climateShadow.State.Reported.Config.Advanced.Rules.AirCon.AutoStartAirCon * 100)

	(*bytes)[0] = 0x53
	(*bytes)[1] = 0x31

	(*bytes)[2] = updateByte((*bytes)[2], getBoolFromCo2SensorRange(climateShadow.State.Reported.Config.Functions.Co2SensorRange), 0)
	(*bytes)[2] = updateByte((*bytes)[2], co2ExtractionEnabled, 1)
	(*bytes)[2] = updateByte((*bytes)[2], climateShadow.State.Reported.Config.Advanced.Rules.Humidification.AllowHumidification, 2)
	(*bytes)[2] = updateByte((*bytes)[2], climateShadow.State.Reported.Config.Advanced.Rules.AirCon.ForceAirCon, 3)
	(*bytes)[2] = updateByte((*bytes)[2], climateShadow.State.Reported.Config.Advanced.FailSafeSettings.LightingOverride.LightFallsAlarmMinimum, 4)
	(*bytes)[2] = updateByte((*bytes)[2], climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.Co2InjectionAllowed, 5)
	(*bytes)[2] = updateByte((*bytes)[2], climateShadow.State.Reported.Config.Advanced.Rules.AllowAirCon, 6)
	(*bytes)[2] = updateByte((*bytes)[2], climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.Co2InjectionAvoid, 7)
	(*bytes)[3] = byte(swAllExhaustFans)
	(*bytes)[4] = byte(swAllExhaustFans >> 8)
	(*bytes)[5] = byte(swOffLightTempExceed)
	(*bytes)[6] = byte(swOffLightTempExceed >> 8)
	(*bytes)[7] = byte(swOffLightsTempExceed)
	(*bytes)[8] = byte(swOffLightsTempExceed >> 8)
	(*bytes)[9] = climateShadow.State.Reported.Config.Advanced.FailSafeSettings.DehumidifierOverride.SwOnFansRhExceed
	(*bytes)[10] = climateShadow.State.Reported.Config.Advanced.FailSafeSettings.DehumidifierOverride.SwAcRhExceed
	(*bytes)[11] = byte(swOnFansCo2Exceed)
	(*bytes)[12] = byte(swOnFansCo2Exceed >> 8)
	(*bytes)[13] = byte(revertFansCo2Falls)
	(*bytes)[14] = byte(revertFansCo2Falls >> 8)
	(*bytes)[15] = climateShadow.State.Reported.Config.Advanced.Rules.Lighting.LampCoolDownTime
	(*bytes)[16] = climateShadow.State.Reported.Config.Advanced.FailSafeSettings.PowerFailure.SwLightsAfterCoolDown
	(*bytes)[17] = climateShadow.State.Reported.Config.Advanced.Rules.Lighting.SwOnNextLightBank
	(*bytes)[18] = climateShadow.State.Reported.Config.Advanced.Rules.SetpointRamping.RampSetpoints
	(*bytes)[19] = climateShadow.State.Reported.Config.Advanced.Rules.HumidifyTempRules.PreventHeater
	(*bytes)[20] = climateShadow.State.Reported.Config.Advanced.Rules.Humidification.ChangeHumidification
	(*bytes)[21] = byte(autoStartAirCon & 0xff)
	(*bytes)[22] = byte(autoStartAirCon >> 8)
	(*bytes)[23] = climateShadow.State.Reported.Config.Advanced.Rules.AirCon.StartBefore
	(*bytes)[24] = byte(lowerCoolingTemp & 0xff)
	(*bytes)[25] = byte(lowerCoolingTemp >> 8)
	(*bytes)[27] = byte(raiseHeatingTemp & 0xff)
	(*bytes)[26] = climateShadow.State.Reported.Config.Advanced.Rules.FoggingRules.FogToCool
	(*bytes)[28] = byte(raiseHeatingTemp >> 8)
	(*bytes)[29] = byte(rhLowThenRaise & 0xff)
	(*bytes)[30] = byte(rhLowThenRaise >> 8)
	(*bytes)[31] = climateShadow.State.Reported.Config.Advanced.Rules.PurgingRules.PurgeMins
	(*bytes)[32] = climateShadow.State.Reported.Config.Advanced.Rules.PurgingRules.PurgeMin
	(*bytes)[33] = climateShadow.State.Reported.Config.Advanced.Rules.PurgingRules.PurgeMax
	(*bytes)[34] = byte(fogToAchieveRh)
	(*bytes)[35] = byte(fogToAchieveRh >> 8)
	(*bytes)[36] = byte(fogTimes)
	(*bytes)[37] = climateShadow.State.Reported.Config.Advanced.Rules.FoggingRules.FogTimeMin
	(*bytes)[38] = climateShadow.State.Reported.Config.Advanced.Rules.FoggingRules.FogTimeMax
	(*bytes)[39] = byte(autoChangeAirCon)
	(*bytes)[40] = byte(autoChangeAirCon >> 8)
	(*bytes)[41] = byte(everyDayMins)
	(*bytes)[42] = byte(everyDayMins >> 8)
	(*bytes)[43] = byte(everyNightMins)
	(*bytes)[44] = byte(everyNightMins >> 8)
	(*bytes)[45] = byte(daySecs)
	(*bytes)[46] = byte(daySecs >> 8)
	(*bytes)[47] = byte(nightSecs)
	(*bytes)[48] = byte(nightSecs >> 8)
	(*bytes)[49] = byte(injectIfLightGreater)
	(*bytes)[50] = byte(injectIfLightGreater >> 8)
	(*bytes)[51] = byte(co2Cycling & 0xff)
	(*bytes)[52] = byte(co2Cycling >> 8)
	(*bytes)[53] = byte(riseVentTemp & 0xff)
	(*bytes)[54] = byte(riseVentTemp >> 8)
	(*bytes)[55] = climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.InjectTimeMin
	(*bytes)[56] = climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.InjectTimeMax
	(*bytes)[57] = climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.WaitTimeMin
	(*bytes)[58] = climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.WaitTimeMax
	(*bytes)[59] = climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.VentTimeMin
	(*bytes)[60] = climateShadow.State.Reported.Config.Advanced.Rules.CO2Rules.VentTimeMax
	(*bytes)[61] = (*d0bytes)[61]
	createCheckSum(bytes)
}

func prepareS2RequestForIClimate(d0bytes *[]byte, bytes *[]byte, climateShadow iClimateShadow) {
	fansOn := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.FansOn * 100)
	fansOff := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.FansOff * 100)
	heaterOn := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.HeaterOn * 100)
	heaterOff := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.HeaterOff * 100)
	humidifierOn := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.HumidifierOn)
	humidifierOff := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.HumidifierOff)
	dehumidifierOn := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.DehumidifierOn)
	dehumidifierOff := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.DehumidifierOff)
	cO2On := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.CO2On)
	cO2Off := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.CO2Off)
	airConditionerOn := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.AirConditionerOn * 100)
	airConditionerOff := int(climateShadow.State.Reported.Config.Advanced.SwitchingOffsets.AirConditionerOff * 100)
	heatingOffset := int(climateShadow.State.Reported.Config.Advanced.Rules.HumidifyTempRules.HeatingOffset * 100)
	fogTimes := int(climateShadow.State.Reported.Config.Advanced.Rules.FoggingRules.FogTimes)

	(*bytes)[0] = 0x53
	(*bytes)[1] = 0x32
	(*bytes)[3] = byte(heaterOn & 0xff)
	(*bytes)[4] = byte(heaterOn >> 8)
	(*bytes)[5] = byte(heaterOff & 0xff)
	(*bytes)[6] = byte(heaterOff >> 8)
	(*bytes)[7] = byte(fansOn & 0xff)
	(*bytes)[8] = byte(fansOn >> 8)
	(*bytes)[9] = byte(fansOff & 0xff)
	(*bytes)[10] = byte(fansOff >> 8)
	(*bytes)[11] = byte(airConditionerOn & 0xff)
	(*bytes)[12] = byte(airConditionerOn >> 8)
	(*bytes)[13] = byte(airConditionerOff & 0xff)
	(*bytes)[14] = byte(airConditionerOff >> 8)
	(*bytes)[15] = byte(humidifierOn & 0xff)
	(*bytes)[16] = byte(humidifierOn >> 8)
	(*bytes)[17] = byte(humidifierOff & 0xff)
	(*bytes)[18] = byte(humidifierOff >> 8)
	(*bytes)[19] = byte(dehumidifierOn & 0xff)
	(*bytes)[20] = byte(dehumidifierOn >> 8)
	(*bytes)[21] = byte(dehumidifierOff & 0xff)
	(*bytes)[22] = byte(dehumidifierOff >> 8)
	(*bytes)[23] = byte(cO2On & 0xff)
	(*bytes)[24] = byte(cO2On >> 8)
	(*bytes)[25] = byte(cO2Off & 0xff)
	(*bytes)[26] = byte(cO2Off >> 8)
	(*bytes)[27] = byte(heatingOffset & 0xff)
	(*bytes)[28] = byte(heatingOffset >> 8)
	(*bytes)[35] = byte(fogTimes >> 8)
	(*bytes)[61] = (*d0bytes)[61]
	createCheckSum(bytes)
}

func parseByteResponseForIDose(d0Response []byte, d1Response []byte, d2Response []byte, name string, timestamp int64) iDoseShadow {
	irrigationInstalled := getBoolFromByte(d1Response[4], 3)
	irrigationIndepended := getBoolFromByte(d1Response[59], 7)
	irrigationSequential := getBoolFromByte(d1Response[58], 7)
	irrigationStations := d1Response[60]
	irrigationMode := getIrrigationMode(irrigationSequential, irrigationIndepended)
	if !irrigationInstalled {
		irrigationStations = 0
	} else {
		if irrigationStations == 0 {
			irrigationStations = 1
		}
	}

	return iDoseShadow{
		StateIDose{
			Reported: ReportedIDose{
				Config: ConfigIDose{
					Units: UnitsIDose{
						DateFormat:  getDateFormat(getBoolFromByte(d1Response[4], 4)),
						Temperature: getTemperatureF(getBoolFromByte(d1Response[4], 7)),
						Ec:          getNutrientConfig(getBoolFromByte(d1Response[4], 6), getBoolFromByte(d1Response[4], 5)),
						TdsConversationStandart: getFloatFrom2Bytes(d1Response[54], d1Response[53]),
					},
					Times: TimesIDose{
						DayStart: int(d1Response[38]) * 6,
						DayEnd:   int(d1Response[39]) * 6,
					},
					Functions: FunctionsIDose{
						NutrientsParts:     getDefaultValue(d1Response[40], 1),
						PhDosing:           getPhDosingModeByBool(getBoolFromByte(d1Response[3], 1), getBoolFromByte(d1Response[3], 0)),
						IrrigationMode:     irrigationMode,
						IrrigationStations: irrigationStations,
						SeparatePumpOutput: getBoolFromByte(d1Response[5], 1),
						UseWater:           getBoolFromByte(d1Response[4], 2),
						ExternalAlarm:      getBoolFromByte(d1Response[2], 2),
						DayNightEc:         getBoolFromByte(d1Response[3], 2),
						IrrigationStation1: getIrrigationStationConfiguration(getBoolFromByte(d1Response[3], 3), getBoolFromByte(d1Response[2], 7)),
						IrrigationStation2: getIrrigationStationConfiguration(getBoolFromByte(d1Response[57], 3), getBoolFromByte(d1Response[57], 2)),
						IrrigationStation3: getIrrigationStationConfiguration(getBoolFromByte(d1Response[58], 3), getBoolFromByte(d1Response[58], 2)),
						IrrigationStation4: getIrrigationStationConfiguration(getBoolFromByte(d1Response[59], 3), getBoolFromByte(d1Response[59], 2)),
						Scheduling:         false,
						MuteBuzzer:         getBoolFromByte(d1Response[5], 4),
					},
					Advanced: AdvancedIDose{
						ProportinalDosing: getBoolFromByte(d1Response[2], 0),
						SequentialDosing:  getBoolFromByte(d1Response[2], 1),
						DisableEc:         getBoolFromByte(d1Response[5], 3),
						DisablePh:         getBoolFromByte(d1Response[5], 2),
						MntnReminderFreq:  "weekly",
					},
					General: GeneralIDose{
						DeviceName: string(d2Response[2:12]),
						Firmware:   prepareInt(getSignedFloatFrom2Bytes(d0Response[8], d0Response[7]), 2, 100),
					},
				},
				Status: StatusIDose{
					Nutrient: NutrientIDose{
						Detent: d1Response[6],
						Ec: EcIDose{
							Min:     prepareByte(d1Response[8], 1, 1) * 10,
							Max:     prepareByte(d1Response[7], 1, 1) * 10,
							Enabled: getBoolFromByte(d1Response[5], 7),
						},
						Ph: PhIDose{
							Min:     prepareByte(d1Response[10], 1, 10),
							Max:     prepareByte(d1Response[9], 1, 10),
							Enabled: getBoolFromByte(d1Response[5], 6),
						},
						NutTemp: NutTempIDose{
							Min:     prepareInt(getFloatFrom2Bytes(d1Response[14], d1Response[13]), 1, 100),
							Max:     prepareInt(getFloatFrom2Bytes(d1Response[12], d1Response[11]), 1, 100),
							Enabled: getBoolFromByte(d1Response[5], 5),
						},
					},
					Status: []StatusStatusIDose{
						{
							Active:   getBoolFromByte(d0Response[15], 7),
							Enabled:  getBoolFromByte(d1Response[2], 4),
							ForceOn:  getBoolFromByte(d1Response[3], 7),
							Function: nutrientDosingFunction,
						},
						{
							Active:   getBoolFromByte(d0Response[15], 6),
							Enabled:  getBoolFromByte(d1Response[5], 0),
							ForceOn:  getBoolFromByte(d1Response[3], 3),
							Function: phFunction,
						},
						{
							Active:   getBoolFromByte(d0Response[17], 7),
							Enabled:  getBoolFromByte(d1Response[2], 5),
							ForceOn:  getBoolFromByte(d1Response[3], 4),
							Function: irrigationFunction,
						},
						{
							Active:   getBoolFromByte(d0Response[17], 7),
							Enabled:  getBoolFromByte(d1Response[2], 5),
							ForceOn:  getBoolFromByte(d1Response[3], 4),
							Function: irrigationStation1Function,
						},
						{
							Active:   getBoolFromByte(d0Response[17], 6),
							Enabled:  getBoolFromByte(d1Response[57], 0),
							ForceOn:  getBoolFromByte(d1Response[57], 4),
							Function: irrigationStation2Function,
						},
						{
							Active:   getBoolFromByte(d0Response[17], 5),
							Enabled:  getBoolFromByte(d1Response[58], 0),
							ForceOn:  getBoolFromByte(d1Response[58], 4),
							Function: irrigationStation3Function,
						},
						{
							Active:   getBoolFromByte(d0Response[17], 4),
							Enabled:  getBoolFromByte(d1Response[59], 0),
							ForceOn:  getBoolFromByte(d1Response[59], 4),
							Function: irrigationStation4Function,
						},
						{
							Active:   getBoolFromByte(d0Response[15], 4),
							Enabled:  getBoolFromByte(d1Response[2], 6),
							ForceOn:  getBoolFromByte(d1Response[3], 5),
							Function: waterFunction,
						},
					},
					SetPoints: SetPointsIDose{
						Nutrient:      prepareInt(getFloatFrom2Bytes(d1Response[16], d1Response[15]), 0, 1),
						NutrientNight: prepareInt(getFloatFrom2Bytes(d1Response[18], d1Response[17]), 0, 1),
						PhDosing:      getPHMode(getBoolFromByte(d1Response[2], 3)),
						Ph:            prepareByte(d1Response[19], 1, 10),
					},
					General: GeneralStatusIDose{
						NutrientDoseTime:    d1Response[21],
						MaxNutrientDoseTime: d1Response[20],
						DoseInterval:        d1Response[28],
						WaterOnTime:         d1Response[37],
						IrrigationInterval1: IrrigationIntervalIDose{
							Day:   getFloatFrom2Bytes(d1Response[32], d1Response[31]),
							Night: getFloatFrom2Bytes(d1Response[34], d1Response[33]),
							Every: getFloatFrom2Bytes(d1Response[36], d1Response[35]),
						},
						IrrigationInterval2: IrrigationIntervalIDose{
							Day:   getFloatFrom2Bytes(d2Response[24], d2Response[23]),
							Night: getFloatFrom2Bytes(d2Response[26], d2Response[25]),
							Every: getFloatFrom2Bytes(d2Response[28], d2Response[27]),
						},
						IrrigationInterval3: IrrigationIntervalIDose{
							Day:   getFloatFrom2Bytes(d2Response[32], d2Response[31]),
							Night: getFloatFrom2Bytes(d2Response[34], d2Response[33]),
							Every: getFloatFrom2Bytes(d2Response[36], d2Response[35]),
						},
						IrrigationInterval4: IrrigationIntervalIDose{
							Day:   getFloatFrom2Bytes(d2Response[40], d2Response[39]),
							Night: getFloatFrom2Bytes(d2Response[42], d2Response[41]),
							Every: getFloatFrom2Bytes(d2Response[44], d2Response[43]),
						},
						IrrigationDuration1: getFloatFrom2Bytes(d1Response[30], d1Response[29]),
						IrrigationDuration2: getFloatFrom2Bytes(d2Response[22], d2Response[21]),
						IrrigationDuration3: getFloatFrom2Bytes(d2Response[30], d2Response[29]),
						IrrigationDuration4: getFloatFrom2Bytes(d2Response[38], d2Response[37]),
						Mix1:                d1Response[22],
						Mix2:                d1Response[23],
						Mix3:                d1Response[24],
						Mix4:                d1Response[48],
						Mix5:                d1Response[49],
						Mix6:                d1Response[50],
						Mix7:                d1Response[51],
						Mix8:                d1Response[52],
						PhDoseTime:          d1Response[26],
						MaxPhDoseTime:       d1Response[25],
					},
				},
				Metrics: MetricsIDose{
					Ec:      checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[10], d0Response[9]), 1, 1),
					PH:      checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[12], d0Response[11]), 1, 100),
					NutTemp: checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[14], d0Response[13]), 1, 100),
				},
				Device:    name,
				Timestamp: timestamp,
				Source:    "Gateway",
				Connected: true,
			},
		},
	}
}

func parseByteResponseForIClimate(d0Response []byte, d1Response []byte, d2Response []byte, d3Response []byte, name string, timestamp int64) iClimateShadow {

	return iClimateShadow{
		State: StateIClimate{
			Reported: ReportedIClimate{
				Config: ConfigIClimate{
					Units: UnitsIClimate{
						DateFormat:  getDateFormat(getBoolFromByte(d1Response[4], 6)),
						Temperature: getTemperatureF(getBoolFromByte(d1Response[4], 7)),
					},
					Functions: FunctionsIClimate{
						Fan1:                        getBoolFromByte(d1Response[2], 7),
						Fan2:                        getBoolFromByte(d1Response[2], 6),
						AirConditioner:              getBoolFromByte(d1Response[2], 5),
						Heater:                      getBoolFromByte(d1Response[2], 4),
						Co2Sensor:                   getBoolFromByte(d1Response[2], 3),
						Co2SensorRange:              getCo2SensorRangeFromBool(getBoolFromByte(d2Response[2], 0)),
						Co2Injection:                getBoolFromByte(d1Response[2], 2),
						Co2Extraction:               getBoolFromByte(d2Response[2], 1),
						Dehumidifier:                getBoolFromByte(d1Response[3], 7),
						Humidifier:                  getBoolFromByte(d1Response[2], 0),
						PulsedFogger:                getBoolFromByte(d1Response[2], 1),
						LightBank1:                  getBoolFromByte(d1Response[3], 6),
						LightsAirColored:            getBoolFromByte(d1Response[3], 0),
						LightBank2:                  getBoolFromByte(d1Response[3], 5),
						LampOverTempShutdownSensors: getBoolFromByte(d1Response[5], 0),
						OutsideTempSensor:           getBoolFromByte(d1Response[3], 4),
						SecondEnviroSensor:          getBoolFromByte(d1Response[3], 3),
						IntruderAlarm:               getBoolFromByte(d1Response[3], 2),
						DehumidifyBy:                getDehumidifyBy(getBoolFromByte(d1Response[4], 0), getBoolFromByte(d1Response[4], 1)),
						Setup:                       "Manual",
						MuteBuzzer:                  getBoolFromByte(d1Response[4], 2),
					},
					Advanced: AdvancedIClimate{
						ViewAdvancedSetting: false,
						SwitchingOffsets: SwitchingOffsetsIClimate{
							FansOn:            prepareInt(getSignedFloatFrom2Bytes(d3Response[8], d3Response[7]), 1, 100),
							FansOff:           prepareInt(getSignedFloatFrom2Bytes(d3Response[10], d3Response[9]), 1, 100),
							HeaterOn:          prepareInt(getSignedFloatFrom2Bytes(d3Response[4], d3Response[3]), 1, 100),
							HeaterOff:         prepareInt(getSignedFloatFrom2Bytes(d3Response[6], d3Response[5]), 1, 100),
							HumidifierOn:      prepareInt(getSignedFloatFrom2Bytes(d3Response[16], d3Response[15]), 1, 1),
							HumidifierOff:     prepareInt(getSignedFloatFrom2Bytes(d3Response[18], d3Response[17]), 1, 1),
							DehumidifierOn:    prepareInt(getSignedFloatFrom2Bytes(d3Response[20], d3Response[19]), 1, 1),
							DehumidifierOff:   prepareInt(getSignedFloatFrom2Bytes(d3Response[22], d3Response[21]), 1, 1),
							CO2On:             prepareInt(getSignedFloatFrom2Bytes(d3Response[24], d3Response[23]), 1, 1),
							CO2Off:            prepareInt(getSignedFloatFrom2Bytes(d3Response[26], d3Response[25]), 1, 1),
							PulsedFoggerOn:    10,
							PulsedFoggerOff:   12,
							AirConditionerOn:  prepareInt(getSignedFloatFrom2Bytes(d3Response[12], d3Response[11]), 1, 100),
							AirConditionerOff: prepareInt(getSignedFloatFrom2Bytes(d3Response[14], d3Response[13]), 1, 100),
						},
						Rules: RulesIClimate{
							MinimumAirChangeRules: MinimumAirChangeRulesIClimate{
								DaySecs:        getFloatFrom2Bytes(d2Response[46], d2Response[45]),
								EveryDayMins:   getFloatFrom2Bytes(d2Response[42], d2Response[41]),
								NightSecs:      getFloatFrom2Bytes(d2Response[48], d2Response[47]),
								EveryNightMins: getFloatFrom2Bytes(d2Response[44], d2Response[43]),
							},
							HumidifyTempRules: HumidifyTempRulesIClimate{
								LowerCoolingTemp: prepareInt(getSignedFloatFrom2Bytes(d2Response[25], d2Response[24]), 1, 100),
								RaiseHeatingTemp: prepareInt(getSignedFloatFrom2Bytes(d2Response[28], d2Response[27]), 1, 100),
								RhLowThenRaise:   prepareInt(getSignedFloatFrom2Bytes(d2Response[30], d2Response[29]), 1, 100),
								PreventHeater:    d2Response[19],
								HeatingOffset:    prepareInt(getSignedFloatFrom2Bytes(d3Response[28], d3Response[27]), 1, 100),
							},
							AllowAirCon: getBoolFromByte(d2Response[2], 6),
							SetpointRamping: SetpointRampingIClimate{
								RampSetpoints: d2Response[18],
							},
							AirCon: AirConIClimate{
								ForceAirCon:      getBoolFromByte(d2Response[2], 3),
								AutoChangeAirCon: prepareInt(getFloatFrom2Bytes(d2Response[40], d2Response[39]), 1, 100),
								StartBefore:      d2Response[23],
								AutoStartAirCon:  prepareInt(getSignedFloatFrom2Bytes(d2Response[22], d2Response[21]), 1, 100),
							},
							CO2Rules: CO2RulesIClimate{
								Co2InjectionAllowed:  getBoolFromByte(d2Response[2], 5),
								InjectIfLightGreater: prepareInt(getFloatFrom2Bytes(d2Response[50], d2Response[49]), 1, 1),
								Co2InjectionAvoid:    getBoolFromByte(d2Response[2], 7),
								Co2Cycling:           prepareInt(getSignedFloatFrom2Bytes(d2Response[52], d2Response[51]), 1, 100),
								RiseVentTemp:         prepareInt(getSignedFloatFrom2Bytes(d2Response[54], d2Response[53]), 1, 100),
								InjectTimeMin:        d2Response[55],
								InjectTimeMax:        d2Response[56],
								WaitTimeMin:          d2Response[57],
								WaitTimeMax:          d2Response[58],
								VentTimeMin:          d2Response[59],
								VentTimeMax:          d2Response[60],
							},
							Humidification: HumidificationIClimate{
								ChangeHumidification: d2Response[20],
								AllowHumidification:  getBoolFromByte(d2Response[2], 2),
							},
							Lighting: LightingIClimate{
								LampCoolDownTime:  d2Response[15],
								SwOnNextLightBank: d2Response[17],
							},
							FoggingRules: FoggingRulesIClimate{
								FogToCool:      d2Response[26],
								FogToAchieveRh: prepareInt(getFloatFrom2Bytes(d2Response[35], d2Response[34]), 1, 100),
								FogTimes:       getFloatFrom2Bytes(d3Response[29], d2Response[36]),
								FogTimeMax:     d2Response[38],
								FogTimeMin:     d2Response[37],
							},
							PurgingRules: PurgingRulesIClimate{
								PurgeMins: d2Response[31],
								PurgeMin:  d2Response[32],
								PurgeMax:  d2Response[33],
							},
						},
						FailSafeSettings: FailSafeSettingsIClimate{
							FanFailOverride: FanFailOverrideIClimate{
								SwOffLightsTempExceed: prepareInt(getFloatFrom2Bytes(d2Response[8], d2Response[7]), 1, 100),
								SwOffLightTempExceed:  prepareInt(getFloatFrom2Bytes(d2Response[6], d2Response[5]), 1, 100),
							},
							DehumidifierOverride: DehumidifierOverrideIClimate{
								SwOnFansRhExceed: d2Response[9],
								SwAcRhExceed:     d2Response[10],
							},
							Co2FailSafe: Co2FailSafeIClimate{
								SwOnFansCo2Exceed: getFloatFrom2Bytes(d2Response[12], d2Response[11]),
							},
							Co2InjectionOverride: Co2InjectionOverrideIClimate{
								RevertFansCo2Falls: getFloatFrom2Bytes(d2Response[14], d2Response[13]),
							},
							PowerFailure: PowerFailureIClimate{
								SwLightsAfterCoolDown: d2Response[16],
							},
							LightingOverride: LightingOverrideIClimate{
								LightFallsAlarmMinimum: getBoolFromByte(d2Response[2], 4),
							},
							AirConOverride: AirConOverrideIClimate{
								SwAllExhaustFans: prepareInt(getFloatFrom2Bytes(d2Response[4], d2Response[3]), 1, 100),
							},
						},
					},
					General: GeneralIClimate{
						DeviceName: string(d1Response[33:43]),
						Firmware:   prepareInt(getSignedFloatFrom2Bytes(d0Response[8], d0Response[7]), 2, 100),
					},
				},
				Status: StatusIClimate{
					Readings: ReadingsIClimate{
						AirTemp: AirTempIClimate{
							Cool:    prepareInt(getFloatFrom2Bytes(d1Response[11], d1Response[10]), 1, 100),
							Heat:    prepareInt(getFloatFrom2Bytes(d1Response[60], d1Response[59]), 1, 100),
							Min:     prepareInt(getFloatFrom2Bytes(d1Response[13], d1Response[12]), 1, 100),
							Max:     prepareInt(getFloatFrom2Bytes(d1Response[15], d1Response[14]), 1, 100),
							Enabled: getBoolFromByte(d1Response[5], 7),
							Page:    false,
						},
						Rh: RhIClimate{
							Target:  d1Response[16],
							Min:     d1Response[17],
							Max:     d1Response[18],
							Enabled: getBoolFromByte(d1Response[5], 6),
							Page:    false,
						},
						Light: LightIClimate{
							Min:     prepareInt(getFloatFrom2Bytes(d1Response[30], d1Response[29]), 0, 1),
							Enabled: getBoolFromByte(d1Response[5], 4),
							Page:    false,
						},
						PowerFail: PowerFailIClimate{
							Enabled: getBoolFromByte(d1Response[5], 2),
							Page:    false,
						},
						FailSafeAlarms: FailSafeAlarmsIClimate{
							Enabled: getBoolFromByte(d1Response[4], 4),
							Page:    true,
						},
						Detent: d1Response[9],
						CO2: CO2IClimate{
							Min:     prepareByte(d1Response[27], 0, 1) * 25,
							Max:     prepareByte(d1Response[28], 0, 1) * 25,
							Target:  prepareInt(getFloatFrom2Bytes(d1Response[26], d1Response[25]), 0, 1),
							Enabled: getBoolFromByte(d1Response[5], 5),
							Page:    false,
						},
						IntruderAlarm: IntruderAlarmIClimate{
							Enabled: getBoolFromByte(d1Response[5], 3),
							Page:    false,
						},
					},
					Statistics: StatisticsIClimate{
						Lights: prepareInt(getFloatFrom2Bytes(d1Response[35], d1Response[34]), 1, 100),
						CO2:    prepareInt(getFloatFrom2Bytes(d0Response[41], d0Response[40]), 0, 1),
					},
					SetPoints: extractSetPointArray(d1Response),
					Status: []StatusStatusIClimate{
						{
							Active:    getBoolFromByte(d0Response[43], 7),
							Enabled:   getBoolFromByte(d1Response[6], 7),
							ForceOn:   getBoolFromByte(d1Response[6], 6),
							Function:  fan1Function,
							Installed: getBoolFromByte(d1Response[2], 7),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 6),
							Enabled:   getBoolFromByte(d1Response[6], 5),
							ForceOn:   getBoolFromByte(d1Response[6], 4),
							Function:  fan2Function,
							Installed: getBoolFromByte(d1Response[2], 6),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 6),
							Enabled:   getBoolFromByte(d1Response[6], 5),
							ForceOn:   getBoolFromByte(d1Response[6], 4),
							Function:  airConFunction,
							Installed: getBoolFromByte(d1Response[2], 5),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 3),
							Enabled:   getBoolFromByte(d1Response[5], 1),
							ForceOn:   false,
							Function:  co2InjectionFunction,
							Installed: getBoolFromByte(d1Response[2], 2),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 3),
							Enabled:   getBoolFromByte(d1Response[5], 1),
							ForceOn:   true,
							Function:  co2ExtractFunction,
							Installed: getBoolFromByte(d1Response[2], 2),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 5),
							Enabled:   getBoolFromByte(d1Response[6], 3),
							ForceOn:   getBoolFromByte(d1Response[6], 2),
							Function:  heaterFunction,
							Installed: getBoolFromByte(d1Response[2], 4),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 2),
							Enabled:   getBoolFromByte(d1Response[6], 1),
							ForceOn:   getBoolFromByte(d1Response[6], 0),
							Function:  dehumidifierFunction,
							Installed: getBoolFromByte(d1Response[3], 7),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 4),
							Enabled:   getBoolFromByte(d1Response[7], 7),
							ForceOn:   getBoolFromByte(d1Response[7], 6),
							Function:  humidifierFunction,
							Installed: getBoolFromByte(d1Response[2], 0),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 1),
							Enabled:   getBoolFromByte(d1Response[7], 5),
							ForceOn:   getBoolFromByte(d1Response[7], 4),
							Function:  lightBank1Function,
							Installed: getBoolFromByte(d1Response[3], 6),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 0),
							Enabled:   getBoolFromByte(d1Response[7], 3),
							ForceOn:   getBoolFromByte(d1Response[7], 2),
							Function:  lightBank2Function,
							Installed: getBoolFromByte(d1Response[3], 5),
						},
						{
							Active:    getBoolFromByte(d0Response[43], 4),
							Enabled:   getBoolFromByte(d1Response[7], 7),
							ForceOn:   getBoolFromByte(d1Response[7], 6),
							Function:  foggerFunction,
							Installed: getBoolFromByte(d1Response[2], 1),
						},
						{
							Active:    false,
							Enabled:   true,
							ForceOn:   false,
							Function:  purgingFunction,
							Installed: true,
						},
					},
					ModeAlarmHistory: ModeAlarmHistoryIClimate{
						Mode: []ModeIClimate{
							{
								Description: "Unknown",
								Timestamp:   "16:00:00",
							},
						},
						Alarms: []AlarmsIClimate{
							{
								Description: "Unknown",
								Timestamp:   "16:00:00",
							},
						},
					},
				},
				Metrics: MetricsIClimate{
					AirTemp:        checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[14], d0Response[13]), 1, 100),
					Rh:             prepareByte(d0Response[17], 0, 1),
					Vpd:            checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[19], d0Response[18]), 1, 1000),
					Light:          checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[35], d0Response[34]), 0, 1),
					PowerFail:      getBoolFromByte(d1Response[2], 1),
					FailSafeAlarms: getBoolFromByte(d1Response[2], 1),
					DayNight:       formatDayNightValue(getBoolFromByte(d1Response[4], 4)),
					Co2:            prepareInt(getFloatFrom2Bytes(d0Response[29], d0Response[28]), 0, 1),
					Intruder:       getBoolFromByte(d0Response[42], 7),
					OutsideTemp:    checkIntForNotAvailableValue(getSignedFloatFrom2Bytes(d0Response[37], d0Response[36]), 1, 100),
					EnviroAirTemp1: checkIntForNotAvailableValue(getSignedFloatFrom2Bytes(d0Response[10], d0Response[9]), 1, 100),
					EnviroAirTemp2: checkIntForNotAvailableValue(getSignedFloatFrom2Bytes(d0Response[12], d0Response[11]), 1, 100),
					EnviroCO21:     checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[25], d0Response[24]), 0, 1),
					EnviroCO22:     checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[27], d0Response[26]), 0, 1),
					EnviroRH1:      prepareByte(d0Response[15], 0, 1),
					EnviroRH2:      prepareByte(d0Response[16], 0, 1),
					EnviroLight1:   checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[31], d0Response[30]), 0, 1),
					EnviroLight2:   checkIntForNotAvailableValue(getFloatFrom2Bytes(d0Response[33], d0Response[32]), 0, 1),
				},
				Device:    name,
				Timestamp: timestamp,
				Source:    "Gateway",
				Connected: true,
			},
		},
	}
}

func extractSetPointArray(d1Response []byte) []SetPointIClimate {
	d1 := d1Response[2:]
	return []SetPointIClimate{
		{
			LightBank:     getLightBoxTextByInt(int(d1[41])),
			LightOn:       getFloatFrom2Bytes(d1[43], d1[42]),
			LightDuration: getFloatFrom2Bytes(d1[45], d1[44]),
			DayTemp:       prepareInt(getFloatFrom2Bytes(d1[47], d1[46]), 1, 100),
			NightDropDeg:  prepareInt(getFloatFrom2Bytes(d1[49], d1[48]), 1, 100),
			RhDay:         getFloatFrom2Bytes(d1[51], d1[50]),
			RhNight:       getFloatFrom2Bytes(d1[53], d1[52]),
			RhMax:         int(d1[54]),
			CO2:           getFloatFrom2Bytes(d1[56], d1[55]),
		},
	}
}

func getBoolFromByte(b byte, bit int) bool {
	binarystr := strconv.FormatInt(int64(b), 2)
	out := []rune(binarystr)
	if len(out) < 8 {
		for i := len(out); i < 8; i++ {
			out = append([]rune{'0'}, out...)
		}
	}
	binary := string(out[bit])
	value, err := strconv.ParseBool(binary)
	tell.IfErrorf(err, "Error parsing bit value")
	return value
}

func getFloatFrom2Bytes(l byte, h byte) int {
	var byteSlice = []byte{l, h}
	data := binary.BigEndian.Uint16(byteSlice)
	return int(data)
}

func getSignedFloatFrom2Bytes(l byte, h byte) int {
	return int(int16(uint16(h) + uint16(l)<<8))
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func prepareByte(b byte, n int, d int) float64 {
	return toFixed(float64(b)/float64(d), n)
}

func prepareInt(b int, n int, d int) float64 {
	return toFixed(float64(b)/float64(d), n)
}

/*
if value is valueUndefined => value is not available
*/
func checkIntForNotAvailableValue(b int, n int, d int) float64 {
	if math.Abs(float64(b)) == valueUndefined {
		return valueUndefined
	}
	return prepareInt(b, n, d)
}

func formatDayNightValue(b bool) string {
	if b {
		return "Day"
	}
	return "Night"
}

func createCheckSum(bytes *[]byte) {
	ccittCrc := crc.CalculateCRC(crc.CRC16, (*bytes)[:len(*bytes)-2])
	h := byte(ccittCrc)
	l := byte(ccittCrc >> 8)
	(*bytes)[62] = h
	(*bytes)[63] = l
}

func updateByte(oldByte byte, bol bool, pos int) byte {
	binarystr := strconv.FormatInt(int64(oldByte), 2)
	out := []rune(binarystr)
	if len(out) < 8 {
		for i := len(out); i < 8; i++ {
			out = append([]rune{'0'}, out...)
		}
	}
	out[pos] = parseBool(bol)
	binarystr = string(out)
	i, err := strconv.ParseInt(binarystr, 2, 64)
	tell.IfErrorf(err, "Error parsing bit value")
	return byte(i)
}

func getDefaultValue(b byte, defB byte) byte {
	if b == 0 {
		return defB
	}
	return b
}

func parseBool(bol bool) rune {
	if bol {
		return '1'
	}
	return '0'
}

func getCo2SensorRangeFromBool(bol bool) string {
	if bol {
		return "5000"
	}
	return "2000"
}

func getBoolFromCo2SensorRange(str string) bool {
	if str == "5000" {
		return true
	}
	return false
}

func getStatusIClimateFunctionByName(statuses []StatusStatusIClimate, statusName string) StatusStatusIClimate {
	for _, status := range statuses {
		if status.Function == statusName {
			return status
		}
	}
	return StatusStatusIClimate{
		Active:  false,
		ForceOn: false,
		Enabled: false,
	}
}

func getStatusIDoseFunctionByName(statuses []StatusStatusIDose, statusName string) StatusStatusIDose {
	for _, status := range statuses {
		if status.Function == statusName {
			return status
		}
	}
	return StatusStatusIDose{
		Active:  false,
		ForceOn: false,
		Enabled: false,
	}
}

func getPHMode(bol bool) string {
	if bol {
		return doserRaise
	}
	return doserLower
}

func getBoolFromPHMode(str string) bool {
	if str == doserRaise {
		return true
	}
	return false
}

func getPhDosingModeByBool(bolHi bool, bolLo bool) string {
	if bolHi && bolLo {
		return doserModeBoth
	}
	if bolHi {
		return doserRaise
	}
	if bolLo {
		return doserLower
	}
	return doserModeNone
}

func getPhDosingBoolByString(mode string) (bool, bool) {
	switch mode {
	case doserModeBoth:
		return true, true
	case doserRaise:
		return true, false
	case doserLower:
		return false, true
	default:
		return false, false
	}
}

func getLightBoxTextByInt(mode int) string {
	switch mode {
	case 0:
		return doserModeNone
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "alt"
	case 4:
		return doserModeBoth
	default:
		return "error"
	}
}

func getLightBoxModeToInt(mode string) int {
	switch mode {
	case doserModeNone:
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "alt":
		return 3
	case doserModeBoth:
		return 4
	default:
		return 99
	}
}

func getIrrigationMode(sequential bool, independent bool) string {
	if sequential && independent {
		return irrigationModeSequential
	}
	if independent {
		return irrigationModeIndependent
	}
	return irrigationModeSingle
}

func getIrrigationStationConfiguration(irrigateTOD bool, irrigateDayOnly bool) string {
	if irrigateDayOnly {
		return irrigationModeDuringDayOnly
	}
	if irrigateTOD {
		return irrigationModeSameTime
	}
	return irrigationModeDayNight
}

func getIrrigationStationConfigurationByString(mode string) (bool, bool) {
	irrigateTOD := false
	irrigateDayOnly := false
	if mode == irrigationModeDuringDayOnly {
		irrigateDayOnly = true
	}
	if mode == irrigationModeSameTime {
		irrigateTOD = true
	}
	return irrigateTOD, irrigateDayOnly
}

func isUSADateFormat(dateFormat string) bool {
	if dateFormat == dateFormatUSA {
		return true
	}
	return false
}

func getDateFormat(usaDateFormat bool) string {
	if usaDateFormat {
		return dateFormatUSA
	}
	return dateFormat
}

func isTemperatureF(temperature string) bool {
	if temperature == temperatureF {
		return true
	}
	return false
}

func getTemperatureF(temperature bool) string {
	if temperature {
		return temperatureF
	}
	return temperatureC
}

func getNutrientConfig(nutConfL bool, nutConfH bool) string {
	if nutConfL {
		return nutrientConfigCF
	}
	if nutConfH {
		return nutrientConfigTDS
	}
	return nutrientConfigEC
}

func isNutrientConfig(nutConf string) (bool, bool) {
	switch nutConf {
	case nutrientConfigCF:
		return false, true
	case nutrientConfigTDS:
		return true, false
	default:
		return false, false
	}
}

func isDehumidifyBy(nutConf string) (bool, bool) {
	switch nutConf {
	case dehumidifyAirCon:
		return false, true
	case dehumidifyPurge:
		return true, false
	default:
		return false, false
	}
}

func getDehumidifyBy(purge bool, ac bool) string {
	if purge {
		return dehumidifyPurge
	}
	if ac {
		return dehumidifyAirCon
	}
	return dehumidifyNone
}
