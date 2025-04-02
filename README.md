# Task Planner
A very simple task scheduler. Tasks are read from a file specified by the configuration.

## Config
| Type | Description |
| ---- | ----------- |
| taskFile | Path to the file containing the list of tasks. |
| anyTimeSymbol | Symbol that will be defined by the scheduler as "always execute" |
| tasksSeparationSymbol | Symbol separating tasks. By default, tasks are separated by lines. |
| multiTimeSeparationSymbol | A character separating a list of values ​​that will cause the task to be executed. The default is a comma. |
| eachSymbol | The symbol at the beginning of the value that will be defined as "fulfill every minute/hour/etc." |

## Standard task file syntax

```
[minute] [hour] [day] [month] [dayOfWeek] [command]
```

| Type | Description |
| ---- | ----------- |
| minute | The minute the task must be completed. |
| hour | The hour the task must be completed. |
| day | The day the task must be completed. |
| month | The month the task must be completed. |
| dayOfWeek | The dayOfWeek the task must be completed. |
| command | A command line instruction that will be executed if all conditions are met. |

Example:

```
0 10 12 11 * node index --dev
```

'*' symbol means that this argument can take any value.

You can also set up a task to run after a certain period of time, for example:

```
/5 10 12 11 * node index --dev
```

You can also add multiple valid values:

```
0,10,25 10 12 11 * node index --dev
```

...Or combine everything together if needed:

```
3,24,26,/5 10 12 11 * node index --dev
```

## Build
```
go build cmd/main.go
```
## Build and run
```
go run cmd/main.go
```