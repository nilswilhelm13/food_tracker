package main

import . "nilswilhelm.net/foodtracker/lib/structs"

type DailyGoalsAndIntake struct {
	goals  DailyGoals
	intake Intake
}

type Mapper struct {
	goals  []DailyGoals
	intake []Intake
}

func NewMapper(goals []DailyGoals, intake []Intake) Mapper {
	return Mapper{
		goals:  goals,
		intake: intake,
	}
}

func (m Mapper) DoMap() map[string]DailyGoalsAndIntake {
	var result = map[string]DailyGoalsAndIntake{}

	for _, goal := range m.goals {
		for _, intake := range m.intake {
			if goal.UserId == intake.UserId {
				result[goal.UserId] = DailyGoalsAndIntake{
					goals:  goal,
					intake: intake,
				}
			}
		}
	}
	return result
}
