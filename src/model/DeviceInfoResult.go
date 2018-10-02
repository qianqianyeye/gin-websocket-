package model

type DeviceInfoResult struct {
	AndroidScreen AndroidScreen
	Device        Device
	DeviceModular DeviceModular
	CommLgz       []CommLgz
	QrCode        string
	WeiMaQi  []CommWeimaqi
}
