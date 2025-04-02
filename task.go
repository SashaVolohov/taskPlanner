// This file describes the entity of the task.
// All fields of the entity are slices of integer values, except for the command itself.
// If the value is negative, this is equivalent to executing the command after this time interval.
package taskplanner

const (
	TaskMinute = iota
	TaskHour
	TaskDay
	TaskMonth
	TaskDayOfWeek
	TaskCommand
	CountOfNumberArguments = TaskCommand
)

const AnyTime = -1
const FirstTimeArgument = 0
const FirstSymbol = 0
const null = 0
const EachInteger = 1

type Task struct {
	minute    []int
	hour      []int
	day       []int
	month     []int
	dayOfWeek []int
	command   string
}

type TaskTimeParameters [CountOfNumberArguments][]int

func IsEachTimeParameterRelevant(eachTimeInteger int, currentTimeInteger int) bool {
	return IsEachTimeParameter(eachTimeInteger) && currentTimeInteger%eachTimeInteger == null
}

func IsEachTimeParameter(time int) bool {
	return time < null
}

func NewTask(minute []int, hour []int, day []int, month []int, dayOfWeek []int, command string) *Task {
	return &Task{
		minute:    minute,
		hour:      hour,
		day:       day,
		month:     month,
		dayOfWeek: dayOfWeek,
		command:   command,
	}
}

func (s *Task) GetTaskTimeParameters(currentParameters TaskTimeParameters) TaskTimeParameters {
	var timeParameters TaskTimeParameters = [CountOfNumberArguments][]int{s.minute, s.hour, s.day, s.month, s.dayOfWeek}

	for i := range timeParameters {
		for j := range timeParameters[i] {

			if timeParameters[i][j] == AnyTime || IsEachTimeParameterRelevant(timeParameters[i][j], currentParameters[i][FirstTimeArgument]) {
				timeParameters[i] = currentParameters[i]
			}

		}
	}

	return timeParameters
}

func (s *Task) GetCommandString() string {
	return s.command
}
