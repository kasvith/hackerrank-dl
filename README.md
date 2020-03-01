# HackerRank Submission Downloader

HackerRank Submission Downloader allows you to download submissions for a contest
in HackerRank.

The tool organizes the submissions according to the problem name and their programming language and then by account name.

For example,

```bash
out/2020-03-01-06-01
├── a-walk-to-remember-1
│   ├── c
│   │   ├── Cygnus_UOK.c
│   │   ├── KingCoders_SEUSL.c
│   │   └── StackTrace_NSBM_.c
│   ├── cpp
│   │   ├── BlackHawks_USJ.cpp
│   │   ├── Code_Clan_RUSL.cpp
│   │   ├── Cyborgs_USJ.cpp
│   │   ├── FRIDAY_USJ.cpp
│   │   ├── JPAC_UOP.cpp
│   │   ├── Paradox_UOJ.cpp
│   │   └── XCODERS_UOP.cpp
│   ├── csharp
│   │   └── CodeMart_UOR.cs
```

## Development

This tool was written in Go

1. Install Golang 1.14 or newer
2. Clone the repo in your local machine
3. Run `go build` to build binary for your platform

### To cross compile

Run `chmod +x ./build/build.sh && ./build/build.sh`

## How to use

Download binary for your platform from [release](https://github.com/kasvith/hackerrank-dl/releases) page (or build yourself :hammer:).

The usage of the tool is simple. You have to provide a `config.yaml` in your current working directory and tool will pick this up. `config.yaml` contains all the necessary info to download submissions from a hackerrank contest.

> Note: You should be a moderator or an administrator of the contest in order to use this tool for download submissions

sample config file `config.yaml` is shown below

```yaml
# slug name of the contest
contest: our-contest
# browser cookies for HackerRank as semicolon seperated string
cookies: >-
  abc=123; xyz=123;
# output directory
output: results
```

Run `./hackerrank-dl` with `config.yaml` in the working directory to start the magic.

## How to get contest slug name

1. Go to HackerRank
2. Go to your contest
3. Slug name is in the pattern of `https://www.hackerrank.com/aces-coders-v8`
4. In here `aces-coders-v8` is the **slug name**
5. Copy it and use it in `config.yaml`

## How to get HackerRank browser cookies

HackerRank does not provide a public API, thus we have to use browser cookies for authentication

> Following instructions are based on Chrome

1. Open Chrome
2. Login to HackerRank
3. Make sure you are an admin/moderator of your contest
4. Install [EditThisCookie](https://chrome.google.com/webstore/detail/editthiscookie/fngmhnnpilhplaeedifhccceomclgfbg/related?utm_source=chrome-ntp-icon) browser extension for chrome
5. Open extension
6. Click **Settings** Icon
7. Set **Choose preferred export format for cookies** to **Semicolon seperated name=value pairs**
8. Open up `https://hackerrank.com`
9. Open `EditThisCookie` extension
10. Click `Export` button (arrow out)
11. Paste the value for cookies in `config.yaml`
