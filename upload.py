#!/usr/bin/env python3

import psycopg2
import toml
from dotenv import load_dotenv
from os import getenv, listdir
from os.path import splitext
from os.path import isdir, join
from utils.db_interface import flyway_checksum_sum, get_user_by_username, \
    create_task, create_version, update_task, ensure_checker, \
    assign_checker, ensure_textfile, create_task_version_test
from utils.validate_task import validate_toml, validate_task_fs

# select sum(checksum) from flyway_schema_history;
OK_DB_SCHEMA_VERSION = -7075107967

OK_TASK_SPEC_VERSION = "1.0"

load_dotenv()

conn = psycopg2.connect(f"""
host={getenv('DB_HOST')} port={getenv('DB_PORT')}
dbname={getenv('DB_NAME')}
user={getenv('DB_USER')} password={getenv('DB_PASS')}
""")

print(f"Connected to {getenv('DB_NAME')} at " +
      f"{getenv('DB_HOST')}:{getenv('DB_PORT')}")

with conn.cursor() as cur:
    db_schema_version = flyway_checksum_sum(cur)
assert db_schema_version == OK_DB_SCHEMA_VERSION, \
    f"Database schema version mismatch: {db_schema_version}"

print("Database schema version OK")

OWNER = getenv('OWNER')
with conn.cursor() as cur:
    owner = get_user_by_username(cur, OWNER)

assert owner is not None, f"OWNER {owner} not found in database"
print(f"OWNER {[owner[i] for i in [1,2,4,5]]} OK")


print("Iterating through all tasks to upload...")
for task_dir in listdir('upload'):
    print(f"Uploading task {task_dir}...")
    task_dirpath = join('upload', task_dir)
    if not isdir(task_dirpath):
        print(f"Error: {task_dir} is not a directory")
        continue
    try:
        cur = conn.cursor()

        problem_toml = toml.load(join(task_dirpath, 'problem.toml'))
        assert problem_toml['specification'] == OK_TASK_SPEC_VERSION, \
            f"Specification version mismatch: {problem_toml['specification']}"
        print("Specification version OK")

        validate_toml(problem_toml)
        print("Validated problem.toml OK")

        validate_task_fs(task_dirpath, problem_toml)
        print("Validated task filesystem OK")

        task_id = create_task(cur, owner[0])
        print(f"Created task {task_id}")

        version_id = create_version(
            cur, task_id,
            problem_toml["code"], problem_toml["name"],
            problem_toml["time"]*1000, problem_toml["memory"]*1024,
            problem_toml["type"]
        )
        print(f"Created version {version_id}")

        update_task(cur, task_id, owner[0], version_id)
        print(f"Assigned relevant version {version_id} to task {task_id}")

        checker_path = join(task_dirpath, 'evaluation', 'checker.cpp')
        with open(checker_path, 'r') as checker_file:
            checker_code = checker_file.read()
        checker_id = ensure_checker(cur, checker_code)
        print(f"Ensured checker {checker_id} exists")

        assign_checker(cur, version_id, checker_id)
        print(f"Assigned checker {checker_id} to version {version_id}")

        # upload testcases
        tests = set()
        tests_path = join(task_dirpath, 'tests')
        for test in listdir(tests_path):
            tests.add(splitext(test)[0])

        for test in tests:
            print(f"Uploading test \"{test}\"...")
            input_path = join(tests_path, f"{test}.in")
            answer_path = join(tests_path, f"{test}.ans")
            input_file = open(input_path, 'rb')
            answer_file = open(answer_path, 'rb')

            input = input_file.read()
            input_text_file_id = ensure_textfile(cur, input)
            print(f"Ensured input textfile {input_text_file_id} exists")

            answer = answer_file.read()
            answer_text_file_id = ensure_textfile(cur, answer)
            print(f"Ensured answer textfile {answer_text_file_id} exists")

            create_task_version_test(
                cur, test, version_id,
                input_text_file_id, answer_text_file_id)
            print(f"Created test {test} for version {version_id}")

            input_file.close()
            answer_file.close()

        conn.commit()

    except Exception as e:
        print(f"Error: {e}")
        conn.rollback()
    finally:
        cur.close()

conn.close()
