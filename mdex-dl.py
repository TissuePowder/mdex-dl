import requests
import sys
import time
import os
import argparse
import zipfile

group_list = {}

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



def get_manga_name(req, api_url, manga_id):

    response = req.get(f"{api_url}/manga/{manga_id}")
    if not check_response_code(response, [403, 404]):
        sys.exit()

    manga_info = response.json()
    manga_name = manga_info['data']['attributes']['title']['en']
    manga_name = manga_name.replace('/', '_')
    return manga_name



def download_chapter(req, api_url, manga_name, chapter):

    global group_list
    chapter_id = chapter['id']
    chapter_hash = chapter['attributes']['hash']
    scanlators = []

    for relationship in chapter['relationships']:

        if relationship['type'] == "scanlation_group":

            group_id = relationship['id']

            if group_id in group_list:
                scanlators.append(group_list[group_id])
            else:
                response = req.get(f"{api_url}/group/{group_id}")

                if check_response_code(response, [403, 404]):
                    group_info = response.json()
                    group_name = group_info['data']['attributes']['name']
                    group_name = group_name.replace('/', '_')
                    scanlators.append(group_name)
                    group_list[group_id] = group_name


    chapter_number = chapter['attributes']['chapter']
    cnum = ""

    if not chapter_number:
        chapter_number = chapter['attributes']['title']
        chapter_dir = f"{manga_name} - {chapter_number}"
    else:
        if "." in chapter_number:
            x = chapter_number.split(".")
            cnum = x[0].zfill(4) + "." + x[1]
        else:
            cnum = chapter_number.zfill(4)
        chapter_dir = f"{manga_name} - c{cnum}"


    if scanlators:
        scanlators = " + ".join(scanlators)
        chapter_dir += f" [{scanlators}]"


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
    filenames = chapter['attributes']['data']

    if not os.path.exists(f"{manga_name}"):
        os.makedirs(f"{manga_name}")

    pnum = 1
    for filename in filenames:
        url = f"{base_url}/data/{chapter_hash}/{filename}"
        ext = os.path.splitext(filename)[1]
        output_zip = f"{manga_name}/{chapter_dir}.cbz"
        output_file = f"p{pnum:04}{ext}"
        if cnum:
            output_file = f"c{cnum}_{output_file}"
        image = req.get(url)
        z = zipfile.ZipFile(output_zip, 'a', compression=zipfile.ZIP_DEFLATED)
        z.writestr(output_file, image.content)
        z.close()
        print(output_file)
        pnum += 1



def download_manga(req, api_url, manga_id, start, to, lang, groups, uploader):

    manga_name = get_manga_name(req, api_url, manga_id)
    if not os.path.exists(manga_name):
        os.mkdir(manga_name)

    offset = 0
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
        if groups:
            params['groups[]'] = groups
        if uploader:
            params['uploader'] = uploader

        response = req.get(f"{api_url}/chapter", params = params)
        if not check_response_code(response, [400, 403]):
            sys.exit()

        chapter_info = response.json()

        if not chapter_info['data']:
            break

        for chapter in chapter_info['data']:

            chapter_number = chapter['attributes']['chapter']

            if not chapter_number:
                chapter_number = chapter['attributes']['title']
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
    parser.add_argument(dest="link", help="link to the manga or chapter")
    parser.add_argument("-s", "--start", type=float, help="chapter number to start downloading from")
    parser.add_argument("-t", "--to", type=float, help="chapter number to stop downloading after")
    parser.add_argument("-l", "--lang", help="comma separated language codes. default is en")
    parser.add_argument("-g", "--groups", help="comma separated UUID's of scanlation groups")
    parser.add_argument("-u", "--uploader", help="UUID of the uploader")
    args = parser.parse_args()

    req = requests.session()
    api_url = 'https://api.mangadex.org'

    link = args.link.split('/')
    manga_id = ""
    chapter_id = ""
    lang = "en"
    groups = ""
    uploader = ""

    if args.lang != None:
        lang = args.lang.split(',')

    if args.groups != None:
        groups = args.groups.split(',')

    if args.uploader != None:
        uploader = args.uploader

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
        download_manga(req, api_url, manga_id, start, to, lang, groups, uploader)

    print("Download finished!")



if __name__ == "__main__":
    main()





