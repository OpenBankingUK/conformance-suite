package model

type ResourceIDs struct {
	AccountIDs   []ResourceAccountID   `json:"account_ids" "validation:min=1"`
	StatementIDs []ResourceStatementID `json:"statement_ids" "validation:min=1"`
}

type ResourceAccountID struct {
	AccountID string `json:"account_id"`
}

type ResourceStatementID struct {
	StatementID string `json:"statement_id"`
}
