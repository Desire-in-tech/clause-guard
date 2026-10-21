package llm

var freelancerGuidance = map[string]string{
	"payment":               `Fair: payment within 14-30 days of invoice, clear amounts or milestones, a late fee or right to pause work, a deposit on larger projects. Risky: 60+ days, payment only after open-ended client "acceptance", pay-when-paid, unclear currency or fees.`,
	"intellectual_property": "Fair: client owns the final deliverables after full payment; freelancer keeps pre-existing tools and background IP and may show the work in a portfolio. Risky: assignment of ALL IP including pre-existing work, assignment before payment, no portfolio right, very broad \"work product\" definitions.",
	"termination":           "Fair: written notice (7-30 days), payment for work done up to termination, a kill fee if the client cancels early. Risky: instant termination without paying for work done, one-sided termination rights.",
	"liability":             "Fair: liability capped (for example at fees paid), no liability for indirect damages. Risky: unlimited liability, broad indemnification of the client.",
	"confidentiality":       "Fair: mutual, clearly defined scope, limited term (about 2-5 years), standard exclusions. Risky: perpetual, one-sided, or vague.",
	"non_compete":           "Fair: none, or narrow scope and short duration limited to direct competitors. Risky: broad exclusivity, long post-contract restrictions, barring work for other clients. Enforceability varies by jurisdiction.",
	"jurisdiction":          "Fair: governing law or venue that is reasonable for both sides, or neutral arbitration. Risky: exclusive foreign courts that make disputes impractical, one-sided arbitration costs.",
	"scope_of_work":         `Fair: clear deliverables, timelines, acceptance criteria, and paid change requests. Risky: vague "other duties as requested", unlimited scope.`,
	"revisions":             "Fair: a defined number of revision rounds, extra work billed. Risky: unlimited revisions without extra pay.",
	"other":                 "judge on overall balance and one-sidedness.",
}

var freelancerGuidanceOrder = []string{
	"payment",
	"intellectual_property",
	"termination",
	"liability",
	"confidentiality",
	"non_compete",
	"jurisdiction",
	"scope_of_work",
	"revisions",
	"other",
}
