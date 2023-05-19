package models

type User struct {
	BaseModel
	DeviceUUID  string
	Username    string
	DeviceModel string
}
