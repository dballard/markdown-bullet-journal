# Markdown Bullet Journal

![Markdown Bullet Journal Logo](https://github.com/dballard/markdown-bullet-journal/raw/master/Markdown-Bullet-Journal.png "Markdown Bullet Journal Logo")

Markdown Bullet Journal is a digital adaptation of [analog tech](http://bulletjournal.com/). When using analog pen + paper, bullet journal migrations are expensive in space and time, however in digital form they are cheap. I found that for my personal productivity having a full markdown todo list file with daily migrations was the most optimal way to manage my time and a digital bullet journal enabled this. I added in a utility to summarize my past work as the daily migrations intentionally removed it. I have also extended this tool to my tastes adding in custom support for repetitive daily tasks and pomodoros.

These are a simple set of utilities that work for me. Nothing fancy.

## Usage

### Windows

Download:

- [mdbj-migreate.exe](https://www.danballard.com/resources/mdbj/mdbj-summary.exe)
- [mdbj-summary.exe](https://www.danballard.com/resources/mdbj/mdbj-migrate.exe)

And place them in a directory you want to use. Run `mdbj-migrate` to generate a template to work from in that directory and each day after to 'migrate' to create the new day's file in that directory. Run `mdbj-summary` to generate summary.txt in the same directory to review work done.

### Linux & Mac

- Install Go

```
go install github.com/dballard/markdown-bullet-journal/tree/master/mdbj-migrate
go install github.com/dballard/markdown-bullet-journal/tree/master/mdbj-summary
```

Pick a directory you want to use and run `mdbj-migreate` to generate a template to work from in that directory. Run `mdbj-migrate` on succesive days and it will find the last dated file in the directory and 'migrate' it. Run `mdbj-summary` in the directory to print a summary of all done work to the console.

### Recommendations

My mdbj directoy is in a cloud backed up location so I can also slightly awkwardly review it from my phone in a text editor.

## Documentation

### mdbj-migrate

When run in a directory, takes the last dated .md file, copies it to a new file with today's date, and drops all lines marked completed (with a '[x]').

### mdbj-summary

Consumes all dated .md files in the directory and prints out all done tasks (lines with '[x]'). Properly collapses nested items into one line names like

```
- Complex task
    - [ ] Subpart A
        - [x] Task 1
```

into

"Complex task / Subpart A / Task 1"

### Markdown supported

The basics of headers with '#'

Nested lists with '-' for notes and indentation

Todo and done with '[ ]' for open todo item, '[x]' for done todo item, and '[-]' for dropped todo item

Obviously you can use other markdown features such as **bold**, *italics* and [Links](https://guides.github.com/features/mastering-markdown/) but none of these trigger any special treatment with regards to Markdown Bullet Journal.

See the included demo file for a better idea.

### Extra Markdown Bullet Journal 'modules'

#### Daily Repetitive Tasks

These are tasks you might want to do a subset of on any given day, and possibly several times. You would like it tracked, but on migration you would like it 'reset to 0' not dropped. In my case I use it with a list of exercises I pick one to do a few times a day.

```
- [x] 4x10 - Pushups
- [ ] 0x10 - Crunches
- [ ] 0x10 - Lunges
- [x] 1x5 - minutes of meditation
```

under `summary` will show as

```
- 40 pushups
- 5 minutes of meditation
```

And then on `migrate` the '4' and '1' will get reset to 0 and the tasks will not get dropped

```
- [ ] 0x10 - Pushups
- [ ] 0x10 - Crunches
- [ ] 0x10 - Lunges
- [ ] 0x5 - minutes of meditation
```

#### Pomodoro ####

If you want to track pomodoro sessions, simply add '.'s inside the square brackets of todo items. They will not be considered done until an 'x' is included and thus will migrate to clean items the next day. They will however count towards pomodoro summaries.

```
- [..] Big Task
    - [x] Part A
    - [x] Part B
    - [x] Part C
    - [ ] Part D
- [ ] Other Task
    - [..x] Thing 1
    - [ ] Thing 2
```

will `migrate` to

```
- [ ] Big Task
    - [ ] Part D
- [ ] Other Task
    - [ ] Thing 2
```

and `summary` will be

```
    Big Task / Part A
    Big Task / Part B
    Big Task / Part C
    Other Task / Thing 1
4 / 8 - 4 Pomodoros
```

