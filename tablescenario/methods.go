package tablescenario

import "github.com/kaiquegarcia/gopest/scenario"

func (table *tableScenario) GivenCases(cases ...tableCase) *tableScenario {
	table.cases = cases
	return table
}

func (table *tableScenario) When(action scenario.Action) *tableScenario {
	table.action = action
	return table
}

func (table *tableScenario) Run() {
	if len(table.cases) == 0 {
		table.test.Fatal("there's no cases to run")
		table.test.FailNow()
	}

	for _, c := range table.cases {
		scenario.New(table.test, "case " + c.title).
			Given(c.args...).
			When(table.action).
			Expect(c.expectations...).
			Run()
	}
}