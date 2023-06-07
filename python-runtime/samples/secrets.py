import os
import yaml

data = os.getenv('DATA')
dct = yaml.safe_load(data)

USER= dct['user']
PASSWORD = os.getenv('SECURE_DATA')

print(USER,PASSWORD)
