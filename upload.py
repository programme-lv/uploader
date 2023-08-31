#!/usr/bin/env python3

import psycopg2
import toml
from dotenv import load_dotenv
from os import getenv, listdir
from os.path import isdir, join
from utils.db_interface import flyway_checksum_sum, get_user_by_username
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

        conn.commit()

    except Exception as e:
        print(f"Error: {e}")
        conn.rollback()
    finally:
        cur.close()

conn.close()
