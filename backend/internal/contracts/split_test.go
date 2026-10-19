package contracts

import "testing"

func TestSplitSeparatesClausesAtBlankLines(t *testing.T) {
	got := Split("  First line\r\ncontinues here  \r\n\r\n \t \r\nSecond clause\r\n")
	if len(got) != 2 {
		t.Fatalf("Split() returned %d clauses, want 2", len(got))
	}
	if got[0].ID != "clause-1" || got[0].Index != 1 || got[0].Text != "First line\ncontinues here" {
		t.Errorf("first clause = %#v", got[0])
	}
	if got[1].ID != "clause-2" || got[1].Index != 2 || got[1].Text != "Second clause" {
		t.Errorf("second clause = %#v", got[1])
	}
}
