package modin

type StaffMessageListInp struct {
	Page   int  `json:"page" d:"1" v:"min:1"`
	Size   int  `json:"size" d:"20" v:"min:1|max:100"`
	Status *int `json:"status" in:"query"`
}

type AdminStaffMessageListInp struct {
	Page     int    `json:"page" d:"1" v:"min:1"`
	Size     int    `json:"size" d:"20" v:"min:1|max:100"`
	Status   *int   `json:"status" in:"query"`
	SenderId uint64 `json:"senderId" in:"query"`
}

type StaffMessageCreateInp struct {
	Subject string `json:"subject" v:"required|length:1,200"`
	Content string `json:"content" v:"required"`
}

type StaffMessageProcessInp struct {
	Id     uint64 `json:"id" in:"path" v:"required"`
	Answer string `json:"answer"`
}
