from utils.db_interface import textfile_exists, \
    get_textfile_id, create_textfile
from cryptography import sha256_file


def ensure_textfile(cursor, filename):
    '''Ensure db has text_file entry and return its id'''
    file_hash = sha256_file(filename)

    if textfile_exists(cursor, file_hash):
        print('Textfile already exists. Fetching ID.')
        return get_textfile_id(cursor, file_hash)
    else:
        print(f'Creating new textfile with hash {file_hash}')
        with open(filename, 'r') as f:
            return create_textfile(cursor, file_hash, f.read())
