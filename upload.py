#!/usr/bin/env python3

import psycopg2
import toml
from dotenv import load_dotenv
from os import getenv, listdir
from os.path import isdir, join
from utils.db_interface import flyway_checksum_sum

# select sum(checksum) from flyway_schema_history;
OK_SCHEMA_VERSION = -7075107967

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
    assert db_schema_version == OK_SCHEMA_VERSION, \
        f"Database schema version mismatch: {db_schema_version}"

print("Database schema version OK")


AUTHOR = getenv('AUTHOR')

print("Iterating through all tasks to upload...")
for task_dir in listdir('upload'):
    print(f"Uploading task {task_dir}...")
    task_dirpath = join('upload', task_dir)
    if not isdir(task_dirpath):
        print(f"Error: {task_dir} is not a directory")
        continue
    try:
        cur = conn.cursor()

        problem_toml = toml.load(join(task_dir, 'problem.toml'))

        code, name = problem_toml['code'], problem_toml['name']
        time_ms = 1000*problem_toml['time']
        memory_kb = 1024*problem_toml['memory']
        type_id = problem_toml['type']
        authors = problem_toml['authors']

        conn.commit()

    except Exception as e:
        print(f"Error: {e}")
        conn.rollback()
    finally:
        cur.close()

conn.close()
