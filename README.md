# rwthcorona-bot

[![telegram](https://img.shields.io/badge/Telegram-%40RWTHcorona__bot-blue?style=social&logo=telegram)](https://t.me/RWTHcorona_bot)

### Running

##### Installation

Install [`golang`](https://golang.org/).

##### Preparing

1. Copy `.env.sample` to `.env` and edit it. (`cp .env.sample .env`)
2. Start with `./scripts/start.sh`.

### Dataset

Currently the data from the NRW Dashboard is used, only for _Aachen (Städteregion)_. [Source](https://www.lzg.nrw.de/covid19/covid19_mags.html).

### ToDo

- [ ] Deployment
    - [ ] Docker with GitHub Actions
    - [ ] Deploy to a raspberry
- [ ] Visualization
    - [ ] Line Chart
    - [ ] Heat Map
- [ ] Commands
    - [x] `/aktuell`
    - [x] `/altersgruppe`
    - [x] `/info`
    - [ ] `/verlauf`
- [ ] Scheduling based on chat (_command to receive automatic updates_)

### License

_just do, what you'd like to do_
