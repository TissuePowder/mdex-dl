import requests
import sys
import time
import os
import argparse


def check_response_code(response, known_errors = []):

    if response.status_code == 200:
        return True
    else:
        if response.status_code in known_errors:
            for error in response.json()['errors']:
                print(error['detail'])
        else:
            print(response, response.text, end='\n')
        return False



# def login(req, api_url, username, password):

#     login_payload = {
#         "username": username,
#         "password": password
#     }
#     response = req.post(f"{api_url}/auth/login", json = login_payload)
#     if not check_response_code(response, [400, 401]):
#         sys.exit()



def get_manga_name(req, api_url, manga_id):

    response = req.get(f"{api_url}/manga/{manga_id}")
    if not check_response_code(response, [403, 404]):
        sys.exit()

    manga_info = response.json()
    manga_name = manga_info['data']['attributes']['title']['en']
    manga_name = manga_name.replace('/', '_')
    return manga_name



def download_chapter(req, api_url, manga_name, chapter):

    chapter_id = chapter['data']['id']
    chapter_hash = chapter['data']['attributes']['hash']
    scanlators = []

    for relationship in chapter['data']['relationships']:

        if relationship['type'] == "scanlation_group":

            group_id = relationship['id']
            response = req.get(f"{api_url}/group/{group_id}")

            if check_response_code(response, [403, 404]):
                group_info = response.json()
                group_name = group_info['data']['attributes']['name']
                group_name = group_name.replace('/', '_')
                scanlators.append(group_name)


    chapter_number = chapter['data']['attributes']['chapter']
    cnum = ""

    if not chapter_number:
        chapter_number = chapter['data']['attributes']['title']
        chapter_dir = chapter_number
    else:
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
        image = req.get(url)
        with open(output, 'wb') as f:
            f.write(image.content)
        print(output)
        pnum += 1



def download_manga(req, api_url, manga_id, start, to, lang):

    manga_name = get_manga_name(req, api_url, manga_id)
    if not os.path.exists(manga_name):
        os.mkdir(manga_name)

    offset = int(start)
    current = start
    counter = 0

    while current < to:

        params = {
            "manga" : manga_id,
            "translatedLanguage[]" : lang,
            "order[chapter]" : "asc",
            "limit" : 100,
            "offset" : offset
        }

        response = req.get(f"{api_url}/chapter", params = params)
        if not check_response_code(response, [400, 403]):
            sys.exit()

        chapter_info = response.json()

        if not chapter_info['results']:
            break

        for chapter in chapter_info['results']:

            chapter_number = chapter['data']['attributes']['chapter']

            if not chapter_number:
                chapter_number = chapter['data']['attributes']['title']
            else:
                current = float(chapter_number)
                if current < start:
                    continue
                if current >= to:
                    break

            download_chapter(req, api_url, manga_name, chapter)
            counter += 1

        offset += 100

    print(f"Total {counter} chapters downloaded.")



def main():

    parser = argparse.ArgumentParser()
    # parser.add_argument("-u", "--username", help="your mangadex username")
    # parser.add_argument("-p", "--password", help="your mangadex password")
    parser.add_argument(dest="link", help="link to the manga or chapter")
    parser.add_argument("-s", "--start", type=float, help="chapter number to start downloading from")
    parser.add_argument("-t", "--to", type=float, help="chapter number to stop downloading after")
    parser.add_argument("-l", "--lang", help="language code: en, ja, es-la etc. default is English or en")
    args = parser.parse_args()

    req = requests.session()
    api_url = 'https://api.mangadex.org'

    # username = args.username
    # password = args.password

    # if username and password:
    #     login(req, api_url, username, password)

    link = args.link.split('/')
    manga_id = ""
    chapter_id = ""
    lang = "en"

    if args.lang != None:
        lang = args.lang

    for i in range(0, len(link)):
        if link[i] == "title":
            manga_id = link[i+1]
            break
        elif link[i] == "chapter":
            chapter_id = link[i+1]
            break

    if manga_id == "" and chapter_id == "":
        print("Invalid link. Provide a valid manga or chapter link from mangadex.")
        sys.exit()

    if args.start != None:
        start = args.start
    else:
        start = float(0)

    if args.to != None:
        to = args.to + 1.0
    else:
        to = float(9999)

    if chapter_id:

        response = req.get(f"{api_url}/chapter/{chapter_id}")
        if not check_response_code(response, [403, 404]):
            sys.exit()
        chapter = response.json()
        manga_name = ""

        for relationship in chapter['data']['relationships']:
            if relationship['type'] == "manga":
                manga_id = relationship['id']
                manga_name = get_manga_name(req, api_url, manga_id)
                break

        if not manga_name:
            manga_name = chapter_id

        download_chapter(req, api_url, manga_name, chapter)

    else:
        download_manga(req, api_url, manga_id, start, to, lang)

    print("Download finished!")



if __name__ == "__main__":
    main()





