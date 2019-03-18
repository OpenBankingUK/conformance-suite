package model

type ResourceIDs struct {
	AccountIDs   []ResourceAccountID   `json:"account_ids"`
	StatementIDs []ResourceStatementID `json:"statement_ids"`
}

type ResourceAccountID struct {
	AccountID string `json:"account_id"`
}

type ResourceStatementID struct {
	StatementID string `json:"statement_id"`
}
