package address

type ReqSetDefault struct {
	Id int `json:"id"`
}

type ReqInfo struct {
	Id int `form:"id"`
}

type ReqCreate struct {
	Region    string `json:"region"`
	Detail    string `json:"detail"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	TypeId    int    `json:"typeId"`
	IsDefault int    `json:"isDefault"`
}

type ReqUpdate struct {
	Id        int    `json:"id"`
	Region    string `json:"region"`
	Detail    string `json:"detail"`
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	TypeId    int    `json:"typeId"`
	IsDefault int    `json:"isDefault"`
}

type ReqDel struct {
	Id int `json:"id"`
}
