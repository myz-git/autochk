import re
import requests

url = 'https://www.oracle.com/database/technologies/'
resp = requests.get(url, timeout=10)
match = re.search(r'Oracle Database (\d{2}ai)', resp.text)
if match:
    print('Latest Oracle DB version:', match.group(1))
