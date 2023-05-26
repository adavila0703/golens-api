package utils

type CronJobScheduleType int

const (
	None CronJobScheduleType = iota
	EveryMinute
	EveryHour
	EveryDayAt12AM
	EveryMondayAt12AM
	EveryMonthOn1stAt12AM
)

func GetCronSchedule(cronSchedule CronJobScheduleType) string {
	switch cronSchedule {
	case EveryMinute:
		return "* * * * *"
	case EveryHour:
		return "0 * * * *"
	case EveryDayAt12AM:
		return "0 0 * * *"
	case EveryMondayAt12AM:
		return "0 0 * * 1"
	case EveryMonthOn1stAt12AM:
		return "0 0 1 * *"
	default:
		return ""
	}
}
