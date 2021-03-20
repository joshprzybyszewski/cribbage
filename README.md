# cribbage

This repo started as a challenge from a friend to calculate the distribution of all hands' points for every possible cribbage hand faster than his MATLAB implementation. Spoilers: golang is faster than MATLAB. Now it has become a playground to mess with unfamiliar technologies (like [mongodb](https://www.mongodb.com), [web assembly](https://webassembly.org/), React, hopefully [aws lambdas](https://aws.amazon.com/lambda/), and more), to show off what I have learned (like golang, MySQL+RDS, designing for scale), and to build up some technical ability between friends.

Cribbage has a ton of [seemingly made up rules](https://bicyclecards.com/how-to-play/cribbage/), and this project is my attempt to turn those into an interactive game for humans.

[![codecov](https://codecov.io/gh/joshprzybyszewski/cribbage/branch/master/graph/badge.svg)](https://codecov.io/gh/joshprzybyszewski/cribbage)

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
make serve
```

- Currently, it will default to a mysql DB. You need to have a mysql server stood up locally and have a database called `cribbage` existing on it.

4. Start playing cribbage.

- Run the React client using `make client`. Navigate to [localhost:3000](localhost:3000) and play!.
- You are able to interact with a barebones HTML client that the gin server has stood up at [localhost:8080/wasm](localhost:8080/wasm) which uses WebAssembly compiled from golang.
  - Using this option, you can create a user, "sign in" as that user, create a game with another user, and play through a game (although the UI is terrible:#). Please note, you need to refresh every time you make an action.
- If you're a sucker for pain, you can use our older "terminal interaction" (which may be broken:#). In a couple terminals, start a couple clients:

```bash
make goclient
```

    - From here, you should be directed through the game using [survey](https://github.com/AlecAivazis/survey).

Happy Playing!

## Legacy Binary

If you'd like to play the first version of our game, you can run the legacy player, which allows you to play dumb and calculated NPCs:

```bash
go run main.go -legacy
```

## Future Vision

On our TODO list:

- Implement websockets for real-time gameplay (we're sick of refreshing a browser after each action)
- Add proper auth using social logins and Oauth2
- Work a bit on infrastructure for a more automated release process; also, buy a real domain name

We will be using AWS free tier as hobbyists to get this deployed out into the cloud. Currently, we have persistent MySQL database in RDS. We're working on getting our app deployed so that you can play from anywhere. Someday, we'd like to have user auth provided by Oauth2 for legit sign-in, push notifications sent out using SNS, [AWS Lambdas](https://aws.amazon.com/lambda/) executing a game's actions, and potentially even settuing up a NoSQL [MongoDB](https://www.mongodb.com/) database in [AWS](https://docs.aws.amazon.com/quickstart/latest/mongodb/overview.html) just for fun.
