package models

type Chat struct {
	ChatType string  `json:"chat_type"`
	UserInfo UserIds `json:"user_ids"`
}

type UserIds struct {
	User1Id int `json:"user_1_id"`
	User2Id int `json:"user_2_id"`
}
