package main

import (
	"bytes"
	"strings"
	"testing"
)

func extractAnswers(output string) []string {
	var answers []string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Belongs: ") {
			answers = append(answers, strings.TrimPrefix(line, "Belongs: "))
		}
	}
	return answers
}

func runTest(input, parserType string) (string, int) {
	scanner := strings.NewReader(input)
	var buf bytes.Buffer
	exitCode := run(scanner, &buf, parserType, "", "")
	return buf.String(), exitCode
}

func TestArithmeticExpressions(t *testing.T) {
	input := `ET
+n()
4
E -> E+T
E -> T
T -> n
T -> (E)
E
(n)
((n))
(n)+(n)+(n)
n
n)
()
exit`

	expectedAnswers := []string{"YES", "YES", "YES", "YES", "NO", "NO"}

	for _, parserType := range []string{"lr0", "lr1"} {
		t.Run(parserType, func(t *testing.T) {
			output, code := runTest(input, parserType)
			if code != 0 {
				t.Errorf("Expected no error, got error code %d", code)
			}
			answers := extractAnswers(output)
			if len(answers) != len(expectedAnswers) {
				t.Errorf("Expected %d answers, got %d", len(expectedAnswers), len(answers))
			}
			for i, ans := range expectedAnswers {
				if i >= len(answers) {
					break
				}
				if answers[i] != ans {
					t.Errorf("Word №%d: expected output %q, got %q", i+1, ans, answers[i])
				}
			}
		})
	}
}

func TestExtendedParentheses(t *testing.T) {
	input := `S
()*
2
S -> (S)S
S -> *
S
(*)*
((*)*)*
*
(*)*)
exit`

	expectedAnswers := []string{"YES", "YES", "YES", "NO"}

	for _, parserType := range []string{"lr0", "lr1"} {
		t.Run(parserType, func(t *testing.T) {
			output, code := runTest(input, parserType)
			if code != 0 {
				t.Errorf("Expected no error, got error code %d", code)
			}
			answers := extractAnswers(output)
			if len(answers) != len(expectedAnswers) {
				t.Errorf("Expected %d answers, got %d", len(expectedAnswers), len(answers))
			}
			for i, ans := range expectedAnswers {
				if i >= len(answers) {
					break
				}
				if answers[i] != ans {
					t.Errorf("Word №%d: expected output %q, got %q", i+1, ans, answers[i])
				}
			}
		})
	}
}

func TestParentheses(t *testing.T) {
	input := `S
()
2
S ->
S -> S(S)
S
()
(())
()()()()((((()()()(()))(())(()))(()))())
())
()()()()((((()()((()))(())(()))(()))())
exit`

	expectedAnswers := []string{"YES", "YES", "YES", "NO", "NO"}

	for _, parserType := range []string{"lr0", "lr1"} {
		t.Run(parserType, func(t *testing.T) {
			output, code := runTest(input, parserType)
			if code != 0 {
				t.Errorf("Expected no error, got error code %d", code)
			}
			answers := extractAnswers(output)
			if len(answers) != len(expectedAnswers) {
				t.Errorf("Expected %d answers, got %d", len(expectedAnswers), len(answers))
			}
			for i, ans := range expectedAnswers {
				if i >= len(answers) {
					break
				}
				if answers[i] != ans {
					t.Errorf("Word №%d: expected output %q, got %q", i+1, ans, answers[i])
				}
			}
		})
	}
}

func TestShiftReduceConflict(t *testing.T) {
	input := `SA
a
3
S ->
S -> A
A -> a
S
a
exit`

	expectedAnswers := []string{"YES"}

	for _, parserType := range []string{"lr0", "lr1"} {
		t.Run(parserType, func(t *testing.T) {
			output, code := runTest(input, parserType)
			if parserType == "lr0" {
				if code == 0 {
					t.Error("Expected error, got no error")
				}
				if !strings.Contains(output, "The grammar is not LR(0)") {
					t.Errorf("Expected grammar error message, got %q", output)
				}
			} else {
				if code != 0 {
					t.Errorf("Expected no error, got error code %d", code)
				}
				answers := extractAnswers(output)
				if len(answers) != len(expectedAnswers) {
					t.Errorf("Expected %d answers, got %d", len(expectedAnswers), len(answers))
				}
				for i, ans := range expectedAnswers {
					if i >= len(answers) {
						break
					}
					if answers[i] != ans {
						t.Errorf("Word №%d: expected output %q, got %q", i+1, ans, answers[i])
					}
				}
			}
		})
	}
}

func TestParenthesesShiftReduceConflict(t *testing.T) {
	input := `S
()
2
S ->
S -> (S)S
S
()
exit`

	expectedAnswers := []string{"YES"}

	for _, parserType := range []string{"lr0", "lr1"} {
		t.Run(parserType, func(t *testing.T) {
			output, code := runTest(input, parserType)
			if parserType == "lr0" {
				if code == 0 {
					t.Error("Expected error, got no error")
				}
				if !strings.Contains(output, "The grammar is not LR(0)") {
					t.Errorf("Expected grammar error message, got %q", output)
				}
			} else {
				if code != 0 {
					t.Errorf("Expected no error, got error code %d", code)
				}
				answers := extractAnswers(output)
				if len(answers) != len(expectedAnswers) {
					t.Errorf("Expected %d answers, got %d", len(expectedAnswers), len(answers))
				}
				for i, ans := range expectedAnswers {
					if i >= len(answers) {
						break
					}
					if answers[i] != ans {
						t.Errorf("Word №%d: expected output %q, got %q", i+1, ans, answers[i])
					}
				}
			}
		})
	}
}

func TestReduceReduceConflict(t *testing.T) {
	input := `SAB
x
4
S -> A
S -> B
A -> x
B -> x
S
x
exit`

	for _, parserType := range []string{"lr0", "lr1"} {
		t.Run(parserType, func(t *testing.T) {
			output, code := runTest(input, parserType)
			if code == 0 {
				t.Error("Expected error, got no error")
			}
			if parserType == "lr0" {
				if !strings.Contains(output, "The grammar is not LR(0)") {
					t.Errorf("Expected grammar error message, got %q", output)
				}
			} else {
				if !strings.Contains(output, "The grammar is not LR(1)") {
					t.Errorf("Expected grammar error message, got %q", output)
				}
			}
		})
	}
}

func TestEmptyLanguage(t *testing.T) {
	input := `

1
S ->
S

a
exit`

	expectedAnswers := []string{"YES", "NO"}

	for _, parserType := range []string{"lr0", "lr1"} {
		t.Run(parserType, func(t *testing.T) {
			output, code := runTest(input, parserType)
			if code != 0 {
				t.Errorf("Expected no error, got error code %d", code)
			}
			answers := extractAnswers(output)
			if len(answers) != len(expectedAnswers) {
				t.Errorf("Expected %d answers, got %d", len(expectedAnswers), len(answers))
			}
			for i, ans := range expectedAnswers {
				if i >= len(answers) {
					break
				}
				if answers[i] != ans {
					t.Errorf("Word №%d: expected output %q, got %q", i+1, ans, answers[i])
				}
			}
		})
	}
}

func TestMultiCWithD(t *testing.T) {
	input := `SC
cd
3
S -> CC
C -> cC
C -> d
S
dd
ccdcccccccdccccd
cccdcccd
ccccccdcccccd
exit`

	expectedAnswers := []string{"YES", "NO", "YES", "YES"}

	for _, parserType := range []string{"lr0", "lr1"} {
		t.Run(parserType, func(t *testing.T) {
			output, code := runTest(input, parserType)
			if code != 0 {
				t.Errorf("Expected no error, got error code %d", code)
			}
			answers := extractAnswers(output)
			if len(answers) != len(expectedAnswers) {
				t.Errorf("Expected %d answers, got %d", len(expectedAnswers), len(answers))
			}
			for i, ans := range expectedAnswers {
				if i >= len(answers) {
					break
				}
				if answers[i] != ans {
					t.Errorf("Word №%d: expected output %q, got %q", i+1, ans, answers[i])
				}
			}
		})
	}
}
