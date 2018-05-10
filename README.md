# Markdown Bullet Journal

![Markdown Bullet Journal Logo](https://github.com/dballard/markdown-bullet-journal/raw/master/Markdown-Bullet-Journal.png "Markdown Bullet Journal Logo")

Markdown Bullet Journal is a digital adaptation of [analog tech](http://bulletjournal.com/). For my personal productivity I found having a full markdown todo list file with daily migrations was the most optimal was to manage my time. I added in a utility to summarize my past work as the daily migrations made that hard to track.

These are a simple set of utilities that work for me. Nothing fancy

## Usage

### Windows

Download:

- [mdbj-migreate.exe](https://www.danballard.com/resources/mdbj/mdbj-summary.exe)
- [mdbj-summary.exe](https://www.danballard.com/resources/mdbj/mdbj-migrate.exe)

And place them in a directory. Run `mdbj-migrate` to generate a template to work from and each day after to 'migrate'. Run `mdbj-summary` to generate summary.txt to review work done.

### Linux & Mac

- Install Go

```
go install github.com/dballard/markdown-bullet-journal/tree/master/mdbj-migrate
go install github.com/dballard/markdown-bullet-journal/tree/master/mdbj-summary
```

Pick a directory you want to use and run `mdbj-migreate` to generate a template to work from. Run it on successive days to 'migrate'. Run `mdbj-summary` to print a summary of done work to the console.

### Recommendations

My mdbj directoy is in a cloud backed up location so I can also slightly awkwardly review it from my phone in a text editor.

## Documentation

### mdbj-migrate

When run in a directory, takes the last dated .md file, copies it to a new file with today's date, and dropes all lines marked completed (with a '[x]').

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

Nested lists with '-' and indentation

Todo and done with '[ ]' and '[x]'

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

Will get output as:

- 40 pushups
- 5 minutes of meditation

And then on migration the '4' and '1' will get reset to 0 and the tasks will not get dropped

