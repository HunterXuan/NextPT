package iamin

type InviteListInp struct {
	Page   int   `json:"page" d:"1" v:"min:1#{#iam.invite.page_min}"`
	Size   int   `json:"size" d:"20" v:"max:100#{#iam.invite.size_max}"`
	Status *uint `json:"status" v:"max:4#{#iam.invite.status_max}"`
}

type InviteSendInp struct {
	Hash  string `json:"hash" v:"required|length:32,32#{#iam.invite.hash_req}|{#iam.invite.hash_len}"`
	Email string `json:"email" v:"required|email#{#iam.invite.email_req}|{#iam.invite.email_fmt}"`
}

type InviteCheckInp struct {
	Hash string `json:"hash" v:"required|length:32,32#{#iam.invite.hash_req}|{#iam.invite.hash_len}"`
}

type InviteGrantInp struct {
	UserId      uint64 `json:"userId" v:"required#{#iam.invite.user_id_req}"`
	Count       int    `json:"count" v:"required|min:1#{#iam.invite.count_req}|{#iam.invite.count_min}"`
	IsTemporary bool   `json:"isTemporary"`
	ExpireDays  int    `json:"expireDays" v:"min:1#{#iam.invite.expire_days_min}"`
}
