# mdex-dl
A simple python script to download manga from https://mangadex.org

API documentation: https://api.mangadex.org/docs.html

Script is tested working with the API version 5.2.28 at the time of writing this.

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

usage: mdex-dl.py [-h] [-s START] [-t TO] [-l LANG] [-g GROUPS] [-u UPLOADER] link

positional arguments:
link                  link to the manga or chapter

optional arguments:
-h, --help            show this help message and exit
-s START, --start START
                      chapter number to start downloading from
-t TO, --to TO        chapter number to stop downloading after
-l LANG, --lang LANG  comma separated language codes. default is en
-g GROUPS, --groups GROUPS
                      comma separated UUID's of scanlation groups
-u UPLOADER, --uploader UPLOADER
                      UUID of the uploader
```

For the most basic feature, just paste a manga or chapter link and you are good to go. For example:
```
$ python mdex-dl.py https://mangadex.org/title/a7c13d5c-3a2d-4dc4-bcd7-74c79d05f88b
```
You can specify a limit or range using the -s/--start or -t/--to options.
```
# Chapter 12 to end
$ python mdex-dl.py https://mangadex.org/title/a7c13d5c-3a2d-4dc4-bcd7-74c79d05f88b -s 12

# From beginning to chapter 20
$ python mdex-dl.py https://mangadex.org/title/a7c13d5c-3a2d-4dc4-bcd7-74c79d05f88b -t 20

# Chapter 12 to 20
$ python mdex-dl.py https://mangadex.org/title/a7c13d5c-3a2d-4dc4-bcd7-74c79d05f88b -s 12 -t 20
```
Pass a language code with -l or --lang option if you want to download chapters in some other language except English. Multiple language codes should be comma-separated. Note that if there are multiple releases of same chapter listed under a language, the downloader will download all of them.

You can also filter chapters by scanlation group or uploader. Pass their UUID with -g/--groups or -u/--uploader options. Multiple group-id's should be comma-separated.
```
$ python mdex-dl.py https://mangadex.org/title/a7c13d5c-3a2d-4dc4-bcd7-74c79d05f88b -l pt-br -g 7300735e-e7dc-4182-baa4-60d5568d4e63
```
## Here is a list of languages and their respective codes.
```
Arabic : ar
Bengali : bn
Bulgarian : bg
Burmese : my
Catalan : ca
Chinese (Simp) : zh
Chinese (Trad) : zh-hk
Czech : cs
Danish : da
Dutch : nl
English : en
Filipino : tl
Finnish : fi
French : fr
German : de
Greek : el
Hebrew : he
Hindi : hi
Hungarian : hu
Indonesian : id
Italian : it
Japanese : ja
Korean : ko
Lithuanian : lt
Malay : ms
Mongolian : mn
Norwegian : no
Persian : fa
Polish : pl
Portuguese (Br) : pt-br
Portuguese (Pt) : pt
Romanian : ro
Russian : ru
Serbo-Croatian : sh
Spanish (Es) : es
Spanish (LATAM) : es-la
Swedish : sv
Thai : th
Turkish : tr
Ukrainian : uk
Vietnamese : vi
```
## Work in progress
I want to keep this script simple, so not too many features will be added. There won't be any feature that requires http POST
request for now. I will keep refining the filtering options. Open a PR if you want to contribute.
