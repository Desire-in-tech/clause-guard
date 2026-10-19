package llm

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"clauseguard/backend/internal/analysis"
	"clauseguard/backend/internal/chat"
	"clauseguard/backend/internal/contracts"
)

type MockClient struct{}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) Analyze(_ context.Context, input []contracts.Clause) (analysis.AnalysisResult, error) {
	clauses := make([]analysis.Clause, 0, len(input))
	counts := map[analysis.Risk]int{
		analysis.RiskGreen:  0,
		analysis.RiskYellow: 0,
		analysis.RiskRed:    0,
	}
	overallRisk := analysis.RiskGreen

	for _, item := range input {
		clause := reviewClause(item)
		clauses = append(clauses, clause)
		counts[clause.Risk]++
		if riskWeight(clause.Risk) > riskWeight(overallRisk) {
			overallRisk = clause.Risk
		}
	}

	summary := fmt.Sprintf(
		"This mock review split the contract into %d clauses: %d marked high risk, %d needing a closer look, and %d lower risk. It is informational only, not legal advice.",
		len(clauses), counts[analysis.RiskRed], counts[analysis.RiskYellow], counts[analysis.RiskGreen],
	)
	return analysis.AnalysisResult{
		Summary:     summary,
		OverallRisk: overallRisk,
		Clauses:     clauses,
	}, nil
}

func reviewClause(item contracts.Clause) analysis.Clause {
	text := strings.ToLower(item.Text)
	result := analysis.Clause{
		ID:               item.ID,
		Index:            item.Index,
		Text:             item.Text,
		Type:             analysis.TypeOther,
		Risk:             analysis.RiskGreen,
		Reason:           "No common risk keyword was detected by the mock reviewer.",
		SuggestedRewrite: "Clarify the parties’ expectations in plain language and confirm the final wording with a qualified professional.",
	}

	switch {
	case hasAny(text, "non-compete", "non compete", "noncompetition", "compete"):
		result.Type = analysis.TypeNonCompete
		result.Risk = analysis.RiskRed
		result.Reason = "A restriction on competing work may limit future client opportunities, especially if its duration, scope, or geography is broad."
		result.SuggestedRewrite = "Limit any restriction to a clearly defined service, client relationship, and short time period, or remove it if it is not essential."
	case hasAny(text, "intellectual property", "work made for hire", "assign all rights", "all rights, title", "ownership of", "copyright"):
		result.Type = analysis.TypeIntellectualProperty
		result.Risk = analysis.RiskYellow
		if hasAny(text, "all rights", "perpetual", "irrevocable", "work made for hire", "assign all") {
			result.Risk = analysis.RiskRed
			result.Reason = "The wording may transfer broad ownership or reuse rights, potentially including work or materials beyond this project."
			result.SuggestedRewrite = "Transfer ownership only in the final deliverables after full payment, and expressly retain ownership of pre-existing tools, templates, and general know-how."
		} else {
			result.Reason = "The clause discusses intellectual property, but the scope of what transfers or remains yours should be explicit."
			result.SuggestedRewrite = "Define the specific deliverables being transferred and exclude pre-existing materials, tools, and general know-how."
		}
	case hasAny(text, "liability", "indemnif", "damages", "hold harmless"):
		result.Type = analysis.TypeLiability
		result.Risk = analysis.RiskYellow
		if hasAny(text, "unlimited", "without limit", "any and all", "regardless of") {
			result.Risk = analysis.RiskRed
			result.Reason = "The clause may create broad or uncapped responsibility for losses, including ones that may be outside your control."
			result.SuggestedRewrite = "Set a clear liability cap tied to fees paid under this contract and limit responsibility to direct losses caused by your own proven breach."
		} else {
			result.Reason = "The clause allocates responsibility for losses; check whether the scope and any cap are clear and balanced."
			result.SuggestedRewrite = "Define covered losses narrowly and include a reasonable liability cap and exceptions for indirect or consequential damages."
		}
	case hasAny(text, "terminate", "termination", "cancel", "cancellation"):
		result.Type = analysis.TypeTermination
		result.Risk = analysis.RiskYellow
		if hasAny(text, "immediately", "without notice", "at any time") {
			result.Risk = analysis.RiskRed
			result.Reason = "One party may be able to end the agreement quickly, which can leave work completed or time reserved unpaid."
			result.SuggestedRewrite = "Require written notice, payment for work completed and approved expenses, and a defined opportunity to remedy a breach before termination."
		} else {
			result.Reason = "The clause covers ending the agreement; check that notice, payment for work completed, and handoff expectations are included."
			result.SuggestedRewrite = "State the notice period, payment due for work completed, and how project materials will be handed over."
		}
	case hasAny(text, "payment", "invoice", "fee", "compensation", "payable", "net 30", "net 60"):
		result.Type = analysis.TypePayment
		result.Risk = analysis.RiskYellow
		if hasAny(text, "waive", "forfeit", "sole discretion", "when satisfied", "upon acceptance only") {
			result.Risk = analysis.RiskRed
			result.Reason = "Payment may depend on a broad condition or be waived, creating uncertainty about when you will be paid."
			result.SuggestedRewrite = "Set a payment due date after invoice receipt, define objective acceptance criteria, and preserve payment for work already completed."
		} else {
			result.Reason = "The clause discusses payment; a specific due date, invoice process, and late-payment handling should be stated."
			result.SuggestedRewrite = "Specify the fee, invoice schedule, payment due date, and a reasonable process for disputed invoice items."
		}
	case hasAny(text, "confidential", "confidentiality", "nondisclosure", "non-disclosure"):
		result.Type = analysis.TypeConfidentiality
		result.Risk = analysis.RiskYellow
		result.Reason = "Confidentiality duties can continue beyond the project; check that the covered information and duration are clearly limited."
		result.SuggestedRewrite = "Define confidential information, exclude information already public or independently developed, and set a reasonable end date."
	case hasAny(text, "jurisdiction", "governing law", "venue", "court"):
		result.Type = analysis.TypeJurisdiction
		result.Risk = analysis.RiskYellow
		result.Reason = "The named law or court location may be inconvenient or unfamiliar if a dispute arises."
		result.SuggestedRewrite = "Name a mutually workable governing law and dispute location, and identify a practical process for resolving disputes."
	case hasAny(text, "revision", "revisions", "rounds of changes", "unlimited changes"):
		result.Type = analysis.TypeRevisions
		result.Risk = analysis.RiskYellow
		if hasAny(text, "unlimited", "without limit", "as needed") {
			result.Risk = analysis.RiskRed
			result.Reason = "The wording may require an unlimited amount of revision work without additional payment."
			result.SuggestedRewrite = "Include a fixed number of revision rounds and define the rate or approval process for additional changes."
		} else {
			result.Reason = "Revision expectations should be tied to a clear number of rounds and a process for extra requests."
			result.SuggestedRewrite = "State the included revision rounds and the fee and approval process for additional work."
		}
	case hasAny(text, "scope", "deliverable", "services", "statement of work", "project"):
		result.Type = analysis.TypeScopeOfWork
		result.Risk = analysis.RiskYellow
		if hasAny(text, "as needed", "as requested", "all services", "from time to time") {
			result.Reason = "The scope may be open-ended, making it hard to tell what work is included in the agreed fee."
			result.SuggestedRewrite = "List concrete deliverables, milestones, exclusions, and a written approval process for out-of-scope requests."
		} else {
			result.Risk = analysis.RiskGreen
			result.Reason = "The clause appears to describe the work; confirm the deliverables, timing, and exclusions are specific."
			result.SuggestedRewrite = "Keep the deliverables, milestones, and exclusions explicit, and require written approval for scope changes."
		}
	}

	return result
}

func (m *MockClient) Chat(_ context.Context, request chat.Request) (chat.Response, error) {
	question := strings.TrimSpace(request.Question)
	if question == "" {
		return chat.Response{}, fmt.Errorf("question must not be empty")
	}
	if len(request.Clauses) == 0 {
		return chat.Response{
			Answer:    "There are no analyzed clauses to cite yet. Run an analysis first. This mock response is informational only, not legal advice.",
			Citations: []string{},
		}, nil
	}

	contextQuestion := question
	if len(strings.Fields(question)) <= 3 {
		for i := len(request.History) - 1; i >= 0; i-- {
			if request.History[i].Role == "user" && strings.TrimSpace(request.History[i].Content) != "" {
				contextQuestion += " " + request.History[i].Content
				break
			}
		}
	}

	keywords := tokens(contextQuestion)
	type match struct {
		clause chat.ClauseReference
		score  int
	}
	matches := make([]match, 0, len(request.Clauses))
	for _, clause := range request.Clauses {
		score := 0
		clauseText := tokens(clause.Text)
		for word := range keywords {
			if _, ok := clauseText[word]; ok {
				score++
			}
		}
		matches = append(matches, match{clause: clause, score: score})
	}
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].score > matches[j].score
	})

	selected := matches[0]
	if selected.score == 0 {
		return chat.Response{
			Answer: fmt.Sprintf(
				"I couldn't find a direct match for that question. Clause 1 may be worth checking: “%s”. This mock response is informational only, not legal advice.",
				shorten(selected.clause.Text, 180),
			),
			Citations: []string{selected.clause.ID},
		}, nil
	}

	answer := fmt.Sprintf(
		"Clause %d says: “%s”. This mock review flags it for a closer read because the wording may affect your obligations. Consider clarifying the terms in writing. This is informational only, not legal advice.",
		selectedIndex(selected.clause, request.Clauses),
		shorten(selected.clause.Text, 220),
	)
	return chat.Response{Answer: answer, Citations: []string{selected.clause.ID}}, nil
}

func selectedIndex(selected chat.ClauseReference, clauses []chat.ClauseReference) int {
	for index, clause := range clauses {
		if clause.ID == selected.ID {
			return index + 1
		}
	}
	return 1
}

func tokens(value string) map[string]struct{} {
	words := strings.FieldsFunc(strings.ToLower(value), func(r rune) bool {
		return r < 'a' || r > 'z'
	})
	result := make(map[string]struct{}, len(words))
	for _, word := range words {
		if len(word) > 2 {
			result[word] = struct{}{}
		}
	}
	return result
}

func hasAny(text string, phrases ...string) bool {
	for _, phrase := range phrases {
		if strings.Contains(text, phrase) {
			return true
		}
	}
	return false
}

func riskWeight(risk analysis.Risk) int {
	switch risk {
	case analysis.RiskRed:
		return 3
	case analysis.RiskYellow:
		return 2
	default:
		return 1
	}
}

func shorten(text string, max int) string {
	text = strings.Join(strings.Fields(text), " ")
	if len([]rune(text)) <= max {
		return text
	}
	runes := []rune(text)
	return string(runes[:max-1]) + "…"
}
