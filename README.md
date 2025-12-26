# PinView

`pinview` is a simple terminal pager inspired by `less`, with a focus on **fixed header/footer** and **keyboard-driven navigation**.

## Features

- Scroll vertically and horizontally
- Pin header and footer lines
- Vim-like key bindings
- Works with files or standard input (pipes)
- Minimal, predictable behavior


## Usage

From a file:

```bash
pinview filename.txt
```


From standard input:

```bash
column -s, -t filename.csv | pinview
```

## Key Bindings


Movement
```
j / k        scroll down / up
h / l        scroll left / right
```

Paging
```
d            page down
u            page up
```

Jump
```
g            go to top
G            go to bottom
```

Pinning
```
T / t        pin / unpin header
B / b        pin / unpin footer
```

Help & Quit
```
?            show help
q            quit
```

#### Help Screen

Press `?` to show a help screen.

The help screen displays a short key summary and exits on any key press.

Press any key to return to pager

This is intentional: during help mode, all input is treated as “return”.
