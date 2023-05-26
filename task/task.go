package task

import (
	"fmt"
	"golens-api/utils"
)

func GetUpdateTaskFunc(cronSchedule utils.CronJobScheduleType) func() {
	switch cronSchedule {
	case utils.EveryMinute:
		return UpdateCoverageTask_EveryMinute
	case utils.EveryDayAt12AM:
		return UpdateCoverageTask_EveryDay
	case utils.EveryMondayAt12AM:
		return UpdateCoverageTask_EveryWeek
	case utils.EveryMonthOn1stAt12AM:
		return UpdateCoverageTask_EveryMonth
	default:
		return func() {}
	}
}

func UpdateCoverageTask_EveryMinute() {
	fmt.Println("test")
}

func UpdateCoverageTask_EveryDay() {

}

func UpdateCoverageTask_EveryWeek() {

}

func UpdateCoverageTask_EveryMonth() {

}

func updateCoverage() {

}

func jobCleanup() {

}
