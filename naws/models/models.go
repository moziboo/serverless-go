package models

// PK : USER#KELLY
// SK : PROFILE
type UserProfile struct {
	PK    string `dynamodbav:"PK" json:"pk"`
	SK    string `dynamodbav:"SK" json:"sk"`
	Email string `dynamodbav:"Email" json:"email"`
	PfP   string `dynamodbav:"PfP" json:"pfp"`
	Phone string `dynamodbav:"Phone" json:"phone"`
	Name  string `dynamodbav:"Name" json:"name"`
}

// PK : USER#KELLY
// SK : POST#001
type Post struct {
	PK          string `dynamodbav:"PK" json:"pk"`
	SK          string `dynamodbav:"SK" json:"sk"`
	Title       string `dynamodbav:"Title" json:"title"`
	Description string `dynamodbav:"Description" json:"description"`
	Body        string `dynamodbav:"Body" json:"body"`
	Date        string `dynamodbav:"Date" json:"date"`
	Image       string `dynamodbav:"Image" json:"image"`
}
