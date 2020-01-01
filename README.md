# cribbage
This repo started as a challenge from a friend to calculate the distribution of all hands' points for every possible cribbage hand faster than his MATLAB implementation. Spoilers: golang is faster than MATLAB. Now it's a place to teach some friends how to code and a playground to mess with unfamiliar technologies in golang (like [mongodb](https://www.mongodb.com), [web assembly](https://webassembly.org/), and hopefully [aws lambdas](https://aws.amazon.com/lambda/) and more).

Cribbage has a ton of [seemingly made up rules](https://bicyclecards.com/how-to-play/cribbage/), and this project is my attempt to turn those into an interactive game for humans.

[![codecov](https://codecov.io/gh/joshprzybyszewski/cribbage/branch/master/graph/badge.svg)](https://codecov.io/gh/joshprzybyszewski/cribbage)

[![Build Status](https://travis-ci.org/joshprzybyszewski/cribbage.png)](https://travis-ci.org/joshprzybyszewski/cribbage)

## How to install

1. Go get this project

```bash
go get -u github.com/joshprzybyszewski/cribbage
```

2. Install everything you need

```bash
make install
```

3. In one terminal, start the server

```bash
go run main.go
```

  - to start with an in-memory database, use `-db=memory`

4. Start playing cribbage. 
  - In a couple other terminals, start a couple clients:

```bash
go run main.go -client
```

  - From here, you should be directed through the game using [survey](https://github.com/AlecAivazis/survey).
  
Happy Playing!

## Legacy Binary
If you'd like to play the first version of our game, you can run the legacy player, which allows you to play dumb and calculated NPCs:
```bash
go run main.go -legacy
```

## Future Vision
We are taking steps to get this into a state where the server can be deployed out into [AWS Lambdas](https://aws.amazon.com/lambda/) with a NoSQL-backed, [MongoDB](https://www.mongodb.com/) database in [AWS](https://docs.aws.amazon.com/quickstart/latest/mongodb/overview.html).
