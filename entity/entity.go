package entity

type AddTask struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type GetSearchTask struct {
	Date  *string `json:"date"`
	Title *string `json:"title"`
}
type DeleteTask struct {
}
type TokenJson struct {
	Token string `json:"token"`
}
type UserPass struct {
	Password string `json:"password"`
}
