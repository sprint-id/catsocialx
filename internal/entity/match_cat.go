package entity

// {
// 	"matchCatId": "",
// 	"userCatId": "",
// 	"message": "" // not null, minLength: 5, maxLength: 120
// }

type MatchCat struct {
	ID        string `json:"id"`
	UserCatId string `json:"user_cat_id"`
	Message   string `json:"message"`
}