package report

type ReportSignatureChain struct {
	Type    string `json:"type"`
	Creator string `json:"creator"`
	Domain  string `json:"domain"`
	Nounce  string `json:"nounce"`
	Value   string `json:"value"`
}
