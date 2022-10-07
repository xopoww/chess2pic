# chess2pic
A tool for visualizing chess positions and games

![](./assets/chess2pic.gif)

## Features

 - parse FEN positions and turn them into PNG images
 - parse PGN games and turn them into GIF animations

*note: current parsers have limited capabilities, refer to docs for more info*


## Build

```
make build
# binaries will be in build/ directory
```


## Usage

Visualize position from FEN notation:
```
chess2pic -notation fen -data "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
```

Or from file:
```
chess2pic -notation fen -in position.fen
```

Create GIFs from PGN games in a similar way:
```
chess2pic -notation pgn -in game.pgn
```

Use `chess2pic -help` for full info on command line arguments.