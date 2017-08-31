# go-intelli

A NATS and REST gateway for the Intelli range of devices.  This allows hackers and tinkerers to do cool stuff
with IntelliDose or IntelliClimate devices connected via USB.  Update events can be subscribed to via NATS or
read via the HTTP API.

Current support is limited to Linux and [exprimentally] MAC.

This can be coupled with the [Jelly SDK](https://github.com/AutogrowSystems/go-jelly) to get programmatical access
to an IntelliDose.

## Installation

Install using go:

    go get github.com/AutogrowSystems/go-intelli

Build using go:

    go build github.com/AutogrowSystems/go-intelli/cmd/intellid

Optionally install and run a [NATS](https://github.com/nats-io/gnatsd/releases) server.

## Usage

To use it, simply run the binary (as sudo to get access to the USB files):

    sudo ./intellid

If you have a NATS server running you will see JSON being published to the subject `intelli.*` or `intelli.ASLID06030112` every 15 seconds.  The JSON output is the same as available in the endpoint as shown below.

You also have two endpoints available, `/devices/count` and `/devices`.  By calling the latter you will see output like the following:

```json
[{
	"serial": "ASLID06030112",
	"name": "",
	"type": "idoze",
	"hid": {
		"path": "/dev/hidraw3",
		"vendor_id": 4292,
		"product_id": 33298,
		"version_number": 0,
		"manufacturer": "ASL",
		"product": "ASL IntelliDose",
		"serial_number": "ASLID06030112",
		"usage_page": 65280,
		"usage": 1,
		"input_report_length": 64,
		"output_report_length": 0
	},
	"shadow": {
		"state": {
			"reported": {
				"config": {
					"units": {
						"date_format": "MM/DD/YY",
						"temperature": "F",
						"ec": "EC",
						"tds_conversation_standart": 500
					},
					"times": {
						"day_start": 360,
						"day_end": 1080
					},
					"functions": {
						"nutrients_parts": 4,
						"ph_dosing": "lower",
						"irrigation_mode": "single",
						"irrigation_stations": 0,
						"separate_pump_output": false,
						"use_water": false,
						"external_alarm": false,
						"day_night_ec": false,
						"irrigation_station_1": "day_night",
						"irrigation_station_2": "day_night",
						"irrigation_station_3": "day_night",
						"irrigation_station_4": "day_night",
						"scheduling": false,
						"mute_buzzer": true
					},
					"advanced": {
						"proportinal_dosing": false,
						"sequential_dosing": true,
						"disable_ec": false,
						"disable_ph": false,
						"mntn_reminder_freq": "weekly"
					},
					"general": {
						"device_name": "IDose\u0000\u0000\u0000\u0000\u0000",
						"firmware": 2.13
					}
				},
				"metrics": {
					"ec": 32768,
					"nut_temp": 32768,
					"pH": 32768
				},
				"status": {
					"general": {
						"dose_interval": 1,
						"nutrient_dose_time": 4,
						"water_on_time": 20,
						"irrigation_interval_1": {
							"day": 20,
							"night": 40,
							"every": 840
						},
						"irrigation_interval_2": {
							"day": 0,
							"night": 0,
							"every": 0
						},
						"irrigation_interval_3": {
							"day": 0,
							"night": 0,
							"every": 0
						},
						"irrigation_interval_4": {
							"day": 0,
							"night": 0,
							"every": 0
						},
						"irrigation_duration_1": 180,
						"irrigation_duration_2": 0,
						"irrigation_duration_3": 0,
						"irrigation_duration_4": 0,
						"max_nutrient_dose_time": 10,
						"max_ph_dose_time": 6,
						"mix_1": 50,
						"mix_2": 100,
						"mix_3": 30,
						"mix_4": 100,
						"mix_5": 100,
						"mix_6": 100,
						"mix_7": 100,
						"mix_8": 100,
						"ph_dose_time": 3
					},
					"nutrient": {
						"detent": 0,
						"ec": {
							"enabled": false,
							"max": 300,
							"min": 50
						},
						"nut_temp": {
							"enabled": false,
							"max": 32,
							"min": 8
						},
						"ph": {
							"enabled": false,
							"max": 7,
							"min": 4.5
						}
					},
					"set_points": {
						"nutrient": 100,
						"nutrient_night": 150,
						"ph_dosing": "lower",
						"ph": 5.8
					},
					"status": [{
						"active": false,
						"enabled": true,
						"force_on": false,
						"function": "Nutrient Dosing"
					}, {
						"active": false,
						"enabled": true,
						"force_on": false,
						"function": "ph"
					}, {
						"active": false,
						"enabled": false,
						"force_on": false,
						"function": "irrigation"
					}, {
						"active": false,
						"enabled": false,
						"force_on": false,
						"function": "Irrigation Station 1"
					}, {
						"active": false,
						"enabled": false,
						"force_on": false,
						"function": "Irrigation Station 2"
					}, {
						"active": false,
						"enabled": false,
						"force_on": false,
						"function": "Irrigation Station 3"
					}, {
						"active": false,
						"enabled": false,
						"force_on": false,
						"function": "Irrigation Station 4"
					}, {
						"active": false,
						"enabled": false,
						"force_on": false,
						"function": "Water"
					}],
					"units": {
						"date_format": "",
						"temperature": "",
						"ec": "",
						"tds_conversation_standart": 0
					}
				},
				"source": "Gateway",
				"device": "ASLID06030112",
				"timestamp": 1504147408,
				"connected": true
			}
		}
	},
	"is_open": true
}]
```
