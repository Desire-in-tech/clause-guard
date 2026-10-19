package analysis

import (
	"context"
	"time"

	"clauseguard/backend/internal/contracts"
)

type Risk string

const (
	RiskGreen  Risk = "green"
	RiskYellow Risk = "yellow"
	RiskRed    Risk = "red"
)

type ClauseType string

const (
	TypePayment              ClauseType = "payment"
	TypeIntellectualProperty ClauseType = "intellectual_property"
	TypeTermination          ClauseType = "termination"
	TypeLiability            ClauseType = "liability"
	TypeConfidentiality      ClauseType = "confidentiality"
	TypeNonCompete           ClauseType = "non_compete"
	TypeJurisdiction         ClauseType = "jurisdiction"
	TypeScopeOfWork          ClauseType = "scope_of_work"
	TypeRevisions            ClauseType = "revisions"
	TypeOther                ClauseType = "other"
)

type Clause struct {
	ID               string     `json:"id" bson:"id"`
	Index            int        `json:"index" bson:"index"`
	Text             string     `json:"text" bson:"text"`
	Type             ClauseType `json:"type" bson:"type"`
	Risk             Risk       `json:"risk" bson:"risk"`
	Reason           string     `json:"reason" bson:"reason"`
	SuggestedRewrite string     `json:"suggested_rewrite" bson:"suggested_rewrite"`
}

type AnalysisResult struct {
	Summary     string   `json:"summary"`
	OverallRisk Risk     `json:"overall_risk"`
	Clauses     []Clause `json:"clauses"`
}

type Report struct {
	ID          string    `json:"id" bson:"_id"`
	Summary     string    `json:"summary" bson:"summary"`
	OverallRisk Risk      `json:"overall_risk" bson:"overall_risk"`
	Clauses     []Clause  `json:"clauses" bson:"clauses"`
	CreatedAt   time.Time `json:"-" bson:"created_at"`
}

type AnalyzeRequest struct {
	ContractText string `json:"contract_text"`
}

type Analyzer interface {
	Analyze(context.Context, []contracts.Clause) (AnalysisResult, error)
}
