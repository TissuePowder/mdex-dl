# mdex-dl
A simple python script to download manga from https://mangadex.org

## Requirements
Python 3.5 or up. You may need to install some modules separately such as requests.
```
$ pip install requests
```
Or if your pip command invokes pip 2, then try with
```
$ pip3 install requests
```

## How to install
Nothing to install. Just download and run the mdex-dl.py file with python3.
```
$ python mdex-dl.py <manga_link>
```

## Usage
```
$ python mdex-dl.py -h

usage: mdex-dl.py [-h] [-s START] [-t TO] [-l LANG] link

positional arguments:
link                  link to the manga or chapter

optional arguments:
-h, --help            show this help message and exit
-s START, --start START
chapter number to start downloading from
-t TO, --to TO        chapter number to stop downloading after
-l LANG, --lang LANG  language code: en, ja, es-la etc. default is English or en
```

For the most basic feature, just paste a manga or chapter link and you are good to go. For example:
```
$ python mdex-dl.py https://mangadex.org/title/a7c13d5c-3a2d-4dc4-bcd7-74c79d05f88b
```
If you want to start downloading from a certain chapter, then pass the chapter number with -s or --start option. The following will start downloading from chapter 12, and keep downloading until all chapters are downloaded.
```
$ python mdex-dl.py https://mangadex.org/title/a7c13d5c-3a2d-4dc4-bcd7-74c79d05f88b -s 12
```
If you want to download upto a certain chapter, then pass the chapter number with -t or --to option. The following will start downloading from the beginning and stop downloading after chapter 20 is downloaded.
```
$ python mdex-dl.py https://mangadex.org/title/a7c13d5c-3a2d-4dc4-bcd7-74c79d05f88b -t 20
```
If you want to download in a range, then pass -s and -t both. The following will download from chapter 12 to 20.
```
$ python mdex-dl.py https://mangadex.org/title/a7c13d5c-3a2d-4dc4-bcd7-74c79d05f88b -s 12 -t 20
```
Pass a language code with the --lang option if you want to download chapters in some other language except English.

## WIP
I want to keep this script too simple, so no plans to add extra features or capabilites as of now. Open a PR if you want to add something.
 
