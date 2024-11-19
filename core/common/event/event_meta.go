package event

type EventMetadata struct {
	Uid       string `json:"uid,omitempty"` // Identity
	UserPodId string `json:"user_pod_id,omitempty"`
	Username  string `json:"username,omitempty"`
	RoleType  int    `json:"role_type,omitempty"`
	NetworkId string `json:"network_id,omitempty"`

	//
	PodId string `json:"pod_id,omitempty"`

	DeviceId    string `json:"device_id,omitempty"`
	DeviceModel string `json:"device_model,omitempty"`
	DeviceName  string `json:"device_name,omitempty"`
	DeviceType  string `json:"device_type,omitempty"`

	GroupType string `json:"group_type,omitempty"` // 群类型
	GroupId   string `json:"group_id,omitempty"`   // 群ID
	GroupName string `json:"group_name,omitempty"` // 群名称

	NameSpace string `json:"namespace,omitempty"`

	ViewType string `json:"view_type,omitempty"`
	ViewId   string `json:"view_id,omitempty"`
	ViewName string `json:"view_name,omitempty"`
}
