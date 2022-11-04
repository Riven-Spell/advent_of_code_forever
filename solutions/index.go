package solutions

import "github.com/Riven-Spell/advent_of_code_forever/inputs"

// InputGenerator generates input with a given complexity (number of elements to compute).
// It is intended for benchmarking, testing, and generating giga inputs.
type InputGenerator func(complexity uint64) (input string, solution *inputs.Solution)

// Solution is an interface that exposes a simple structure:
// The runner should call Prepare() on it to prepare the input
// Then call the part 1 and part 2 functions if wanted.
type Solution interface {
	Prepare(input string)
	Part1() any // runtime.DeepEquals is used for comparisons
	Part2() any // runtime.DeepEquals is used for comparisons
}

type Year struct {
	Days    [25]*Day
	Year    uint
	LastDay uint
}

type Day struct {
	Solution          Solution
	Generator         InputGenerator // not mandatory for solution but mandatory for benchmarking
	DefaultComplexity uint64
}

type SolutionIndex struct {
	years    map[uint]*Year
	lastYear uint
}

var Index = &SolutionIndex{
	years: map[uint]*Year{
		2015: {
			LastDay: 0, // This is invalid but create will fix it
		},
	},
	lastYear: 2015,
}

func (i *SolutionIndex) GetCurrentDay() (day, year uint) {
	return i.years[i.lastYear].LastDay, i.lastYear
}

func (i *SolutionIndex) GetCurrentDayForYear(year uint) uint {
	targetYear, ok := i.years[year]
	if ok {
		return targetYear.LastDay
	} else {
		return 0
	}
}

func (i *SolutionIndex) Get(day, year uint) *Day {
	targetYear, ok := i.years[year]
	if !ok || day < 1 || day > 25 {
		return nil
	}

	return targetYear.Days[day-1]
}

func (i *SolutionIndex) Insert(day, year uint, solution *Day) {
	if year > i.lastYear {
		i.lastYear = year
	}

	targetYear, ok := i.years[year]
	if !ok {
		targetYear = &Year{
			Days:    [25]*Day{},
			Year:    year,
			LastDay: day,
		}
		i.years[year] = targetYear
		ok = true
	}

	if day > targetYear.LastDay {
		targetYear.LastDay = day
	}

	targetYear.Days[day-1] = solution
}
