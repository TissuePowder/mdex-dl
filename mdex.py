import requests
import os
import sys
import time
import enum
import argparse


parser = argparse.ArgumentParser()
parser.add_argument("-u", "--username", help="your mangadex username")
parser.add_argument("-p", "--password", help="your mangadex password")
parser.add_argument("-m", "--manga-id", required=True, help="ID of the manga")
parser.add_argument("-s", "--start", type=float, help="chapter number to start downloading from")
parser.add_argument("-t", "--to", type=float, help="chapter number to stop downloading after")
parser.add_argument("-k", "--keep-name", action="store_true", help="keep original filename")
parser.add_argument("-v", "--verbose", action="store_true", help="turns on verbosity")
args = parser.parse_args()

req = requests.session()
api_url = 'https://api.mangadex.org'

response = req.get(api_url)
if response.status_code != 200:
    print(response)
    print(response.text)
    sys.exit()

username = args.username
password = args.password

if username and password:
    login_payload = {
        "username": username,
        "password": password
    }
    response = req.post(f"{api_url}/auth/login", json = login_payload)
    if response.status_code != 200:
        if response.status_code == 400 or response.status_code == 401:
            for error in response.json()['errors']:
                print(error['detail'])
        else:
            print(response)
            print(response.text)
        sys.exit()

else:
    print("Downloader will be working without logging in.")

manga_id = args.manga_id

response = req.get(f"{api_url}/manga/{manga_id}")
if response.status_code != 200:
    if response.status_code == 403 or response.status_code == 404:
        for error in response.json()['errors']:
            print(error['detail'])
    else:
        print(response)
        print(response.text)
    sys.exit()

manga_info = response.json()

manga_name = manga_info['data']['attributes']['title']['en']
manga_name = manga_name.replace('/', '_')

if not os.path.exists(manga_name):
    os.mkdir(manga_name)


if args.start != None:
    start = float(args.start)
else:
    start = float(0)

if args.to != None:
    to = float(args.to + 1.0)
else:
    to = float(9999999)


offset = max(0, int(start)-20)
current = start

while current < to:

    params = {
        "manga" : manga_id,
        "translatedLanguage[]" : "en",
        "order[chapter]" : "asc",
        "limit" : 100,
        "offset" : offset
    }

    response = req.get(f"{api_url}/chapter", params = params)

    if response.status_code != 200:
        if response.status_code == 400 or response.status_code == 403:
            for error in response.json()['errors']:
                print(error['detail'])
        else:
            print(response)
            print(response.text)
        sys.exit()

    chapter_info = response.json()

    if not chapter_info['results']:
        break

    #print(chapter_info)

    for chapter in chapter_info['results']:

        chapter_id = chapter['data']['id']
        chapter_hash = chapter['data']['attributes']['hash']
        scanlators = []

        for relationship in chapter['relationships']:

            if relationship['type'] == "scanlation_group":
                group_id = relationship['id']
                response = req.get(f"{api_url}/group/{group_id}")
                if response.status_code != 200:
                    if response.status_code == 403 or response.status_code == 404:
                        for error in response.json()['errors']:
                            print(error['detail'])
                    else:
                        print(response)
                        print(response.text)
                else:
                    group_info = response.json()
                    group_name = group_info['data']['attributes']['name']
                    group_name = group_name.replace('/', '_')
                    scanlators.append(group_name)

        chapter_number = chapter['data']['attributes']['chapter']
        current = float(chapter_number)
        if current < start:
            continue
        if current >= to:
            break

        if "." in chapter_number:
            x = chapter_number.split(".")
            cnum = x[0].zfill(3) + "." + x[1]
        else:
            cnum = chapter_number.zfill(3)

        chapter_dir = 'Chapter ' + cnum
        if scanlators:
            chapter_dir += " " + "[" + " + ".join(scanlators) + "]"


        status_code = 404
        tries = 0
        while status_code != 200:
            if tries > 5:
                print(f"Problem with MangaDex@Home network, couldn't fetch chapter {chapter_number}.")
                sys.exit()

            response = req.get(f"{api_url}/at-home/server/{chapter_id}")
            status_code = response.status_code

            if status_code == 200:
                break
            else:
                tries += 1
                time.sleep(3)


        base_url = response.json()['baseUrl']
        filenames = chapter['data']['attributes']['data']

        if not os.path.exists(f"{manga_name}/{chapter_dir}"):
            os.makedirs(f"{manga_name}/{chapter_dir}")

        pnum = 1
        for filename in filenames:
            url = f"{base_url}/data/{chapter_hash}/{filename}"
            ext = os.path.splitext(filename)[1]
            output = f"{manga_name}/{chapter_dir}/c{cnum}_p{pnum:04}{ext}"
            print(output)
            image = req.get(url)
            with open(output, 'wb') as f:
                f.write(image.content)
            pnum += 1

    offset += 100

print("Download finished!")