import psycopg2
import toml
from dotenv import load_dotenv
from os import getenv, listdir
from os.path import isdir, join

load_dotenv()

DB_HOST, DB_PORT = getenv('DB_HOST'), getenv('DB_PORT')
DB_USER, DB_PASS = getenv('DB_USER'), getenv('DB_PASS')
DB_NAME = getenv('DB_NAME')

AUTHOR = getenv('AUTHOR')

conn = psycopg2.connect(f"""
host={DB_HOST} port={DB_PORT}
dbname={DB_NAME} user={DB_USER} password={DB_PASS}
""")

for task_dir in listdir('upload'):
    if not isdir(task_dir):
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
