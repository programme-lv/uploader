import os
from db_interface import create_version_test, textfile_exists, \
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


def import_tests(cursor, task_dir):
    '''Imports all tests from the given task directory'''
    test_dir = f'{task_dir}/tests'

    testnames = set()
    for filename in os.listdir(test_dir):
        testnames.add(filename.split('.')[0])

    for testname in testnames:
        in_filename = f'{test_dir}/{testname}.in'
        print(f'Processing {in_filename}')
        in_id = ensure_textfile(cursor, in_filename)

        ans_filename = f'{test_dir}/{testname}.ans'
        print(f'Processing {ans_filename}')
        ans_id = ensure_textfile(cursor, ans_filename)

        create_version_test(cursor, in_id, ans_id, testname)
        print(f'Test {testname} added.')
