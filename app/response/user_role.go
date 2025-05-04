package response

type List_User_Role struct {
	User_ID string `bun:"user_id" json:"user_id"`
	Role_ID string `bun:"role_id" json:"role_id"`
}

type UserRoleByUserID struct {
	User_ID    string `bun:"user_id" json:"user_id"`
	First_Name string `bun:"first_name" json:"first_name"`
	Last_Name  string `bun:"last_name" json:"last_name"`
	Role_ID    string `bun:"role_id" json:"role_id"`
	Role_Name  string `bun:"role_name" json:"role_name"`
}
