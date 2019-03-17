package model

type ResourceIDs struct {
	AccountIDs   []ResourceAccountID   `json:"account_ids" validate:"min=1"`
	StatementIDs []ResourceStatementID `json:"statement_ids" validate:"min=1"`
}

type ResourceAccountID struct {
	AccountID string `json:"account_id" validate:"not_empty"`
}

type ResourceStatementID struct {
	StatementID string `json:"statement_id" validate:"not_empty"`
}
