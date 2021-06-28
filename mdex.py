import requests
import os
import sys
import argparse


parser = argparse.ArgumentParser()
parser.add_argument("-u", "--username", help="your mangadex username")
parser.add_argument("-p", "--password", help="your mangadex password")
parser.add_argument("-m", "--manga-id", required=True, help="ID of the manga")
parser.add_argument("-s", "--start", help="chapter number to start downloading from")
parser.add_argument("-t", "--to", help="chapter number to stop downloading after")
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


#print(manga_info)

manga_name = manga_info['data']['attributes']['title']['en']

#manga_vol_count = manga_info['data']['attributes']['lastVolume']
#volume_pad = len(manga_vol_count)

#manga_chapter_count = manga_info['data']['attributes']['lastChapter']
#chapter_pad = len(manga_chapter_count)


if not os.path.exists(manga_name):
    os.mkdir(manga_name)


loop = 100#int(manga_chapter_count)

while loop > 0:

    offset = 0

    params = {
        "manga" : manga_id,
        "translatedLanguage[]" : "en",
        "limit" : 100,
        "offset" : offset
    }

    chapter_info = req.get(f"{api_url}/chapter", params = params).json()

    print(chapter_info)

    for ch in chapter_info['results']:

        ch_id = ch['data']['id']
        ch_hash = ch['data']['attributes']['hash']
        volume_dir = 'volume ' + (ch['data']['attributes']['volume'])#.zfill(volume_pad)
        chapter_dir = 'chapter ' + (ch['data']['attributes']['chapter'])#.zfill(chapter_pad)
        filenames = ch['data']['attributes']['data']
        base_url = req.get(f"{api_url}/at-home/server/{ch_id}").json()['baseUrl']

        print(volume_dir)

        if not os.path.exists(f"{manga_name}/{volume_dir}/{chapter_dir}"):
            os.makedirs(f"{manga_name}/{volume_dir}/{chapter_dir}")

        num = 1
        for filename in filenames:
            url = f"{base_url}/data/{ch_hash}/{filename}"
            ext = os.path.splitext(filename)[1]
            output = f"{manga_name}/{volume_dir}/{chapter_dir}/{num:03}{ext}"
            image = req.get(url)
            with open(output, 'wb') as f:
                f.write(image.content)
            num += 1

            print(output)

        loop -= 100
        offset += 100

