# Markdown Bullet Journal

![Markdown Bullet Journal Logo](https://github.com/dballard/markdown-bullet-journal/raw/master/src/Markdown-Bullet-Journal.png "Markdown Bullet Journal Logo")

Markdown Bullet Journal is a digital adaptation of analog tech. I found having a running todo list with daily migrations dropping done items worked best for my workflow. Add a simple summary app to show you all you have accomplished and I'm happy.

These are a simple set of utilities that work for me. Nothing fancy

## mdbj-migrate

When run in a directory, takes the last dated .md file, copies it to a new file with today's date, and dropes all lines marked completed (with a '[x]').

## mdbj-summary

Consumes all dated .md files in the directory and prints out all done tasks (lines with '[x]'). Properly collapses nested items into one line names like

- Complex task
    - [ ] Subpart A
        - [x] Task 1

into

"Complex task / Subpart A / Task 1"

## Markdown supported

The basics of headers with '#'

Nested lists with '-' and indentation

Todo and done with '[ ]' and '[x]'

Obviously you can use other markdown features such as **bold**, *italics* and [Links](https://guides.github.com/features/mastering-markdown/) but none of these trigger any special treatment with regards to Markdown Bullet Journal.

See the included demo file for a better idea.

### Extra Markdown Bullet Journal 'modules'

#### Daily Repetitive Tasks

These are tasks you might want to do a subset of on any given day, and possibly several times. You would like it tracked, but on migration you would like it 'reset to 0' not dropped. In my case I use it with a list of exercises I pick one to do a few times a day.

- [x] 4x10 - Pushups
- [ ] 0x10 - Crunches
- [ ] 0x10 - Lunges
- [x] 1x5 - minutes of meditation

Will get output as:

- 40 pushups
- 5 minutes of meditation

And then on migration the '4' and '1' will get reset to 0 and the tasks will not get dropped

