# qgrep

Search Stack Overflow and get real answers, not just links, without leaving your terminal.

`qgrep` is for the "I forgot the syntax for X" moment. Instead of opening a browser, clicking through a few tabs, and skimming for the one useful sentence, you run one command and get the top answers printed right there.

```
qgrep golang for loops
```

## What it does

1. Searches Stack Overflow for your query (sorted by votes).
2. Takes the top 3 matching questions.
3. Fetches the highest-voted answer for each one.
4. Strips the HTML down to readable text, with code blocks wrapped in `[code]` / `[/code]` markers so they stand out from the prose.

Each result shows the question title, score, a link back to the full thread, and the answer itself:

```
Title: <question title>
Score: <votes>
Link: <stackoverflow url>
Answer: <explanation text>
[code]
<code from the answer>
[/code]
<more explanation>
=====================================================================
```

## Requirements

- Go 1.20 or newer (check with `go version`)
- Internet access (it talks to the public Stack Exchange API, no API key needed)

## Install

Linux / macOS:

```
git clone https://github.com/Peejidle/qgrep.git
cd qgrep
go mod tidy
mkdir -p ~/.local/bin
go build -o ~/.local/bin/qgrep .
```

`go mod tidy` downloads the one external dependency, [goquery](https://github.com/PuerkitoBio/goquery).

Make sure `~/.local/bin` is on your `PATH`. Check with:

```
echo $PATH
```

If it isn't listed, add this line to your shell config (`~/.bashrc` or `~/.zshrc`), then restart your terminal or `source` the file:

```
export PATH="$HOME/.local/bin:$PATH"
```

Confirm it worked:

```
which qgrep
```

Windows: I haven't set up install steps for Windows. You can still run it from the project folder with `go run . <your query>`.

## Usage

```
qgrep <your search query>
```

The query is just whatever you type after the command, so no quotes are needed:

```
qgrep golang read file line by line
qgrep golang switch case
qgrep golang defer
```

Specific queries work better than broad ones. `golang for loops` will pull in anything loosely related to loops, while `golang range over map` narrows it down.

## Notes and limitations

- **API rate limit:** the Stack Exchange API allows 300 requests per day per IP without a key. Each `qgrep` run makes up to 4 requests (1 search + 1 answer lookup for each of the 3 results), so that's roughly 75 searches a day.
- **Only the top answer** is shown for each question. The full thread is one click away via the link.
- **Formatting is basic.** Only paragraphs and code blocks are kept. Lists, blockquotes, and inline formatting are dropped from the answer text.
- Stack Overflow only, for now.

## Roadmap

- [ ] Refactor `main` by pulling the per-question work into a helper function
- [ ] Flags: choose how many results to show, and skip ahead to the next set
- [ ] Tests with `go test`
- [ ] GitHub Actions CI
- [ ] Better terminal formatting for code blocks
- [ ] Preserve lists and other formatting from answers
- [ ] Maybe someday: other search modes (news, reference lookups) behind the same interface

## How it's built

Single-file Go program using only the standard library plus goquery:

- `net/http` and `encoding/json` for the Stack Exchange API (`search/advanced` to find questions, `questions/{id}/answers` to fetch answer bodies)
- `goquery` to parse the HTML in each answer and walk its elements in order, so prose and code stay interleaved the way the original author wrote them

## Credits

Answer content comes from Stack Overflow and is licensed under [CC BY-SA](https://creativecommons.org/licenses/by-sa/4.0/). Each result includes a link back to the original post.
