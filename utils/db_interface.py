from utils.cryptography import sha256_bytes


def create_user(
        cursor, username, email, hashed_password,
        first_name, last_name, is_admin=False):
    '''Creates a new user and returns their ID'''
    cursor.execute('''
        INSERT INTO users
        (username, email, hashed_password, first_name, last_name,is_admin)
        VALUES (%s, %s, %s, %s, %s, %s)
        ON CONFLICT DO NOTHING
        RETURNING id
    ''', (username, email, hashed_password, first_name, last_name, is_admin))
    return cursor.fetchone()[0]


def user_exists(cursor, username):
    '''Returns True if a user with the given username exists'''
    cursor.execute('''
        SELECT EXISTS(
            SELECT 1
            FROM users
            WHERE username = %s
        )
    ''', (username,))
    return cursor.fetchone()[0]


def get_user_id(cursor, username):
    '''Returns the ID of a user with the given username'''
    cursor.execute('''
        SELECT id FROM users
        WHERE username = %s
    ''', (username,))
    return cursor.fetchone()[0]


def get_user_by_username(cursor, username):
    if not user_exists(cursor, username):
        return None

    '''Returns a user with the given username'''
    cursor.execute('''
        SELECT * FROM users
        WHERE username = %s
    ''', (username,))
    return cursor.fetchone()


def create_task(cursor, created_by_id,
                relevant_version_id=None, published_version_id=None) -> int:
    '''Creates a new task and returns its ID'''
    cursor.execute('''
        INSERT INTO tasks
        (created_by_id, relevant_version_id, published_version_id)
        VALUES (%s, %s, %s)
        RETURNING id
    ''', (created_by_id, relevant_version_id, published_version_id))
    return cursor.fetchone()[0]


def update_task(cursor, task_id, created_by_id,
                relevant_version_id=None, published_version_id=None):
    cursor.execute('''
        UPDATE tasks
        SET created_by_id = %s, relevant_version_id = %s,
        published_version_id = %s
        WHERE id = %s
    ''', (created_by_id, relevant_version_id, published_version_id, task_id))


def create_version(cursor, task_id, short_code, full_name,
                   time_lim_ms, mem_lim_kb, testing_type_id, origin=None,
                   checker_id=None, interactor_id=None):
    '''Creates a new task version and returns its ID'''
    cursor.execute('''
        INSERT INTO task_versions
        (task_id, short_code, full_name, time_lim_ms, mem_lim_kb,
        testing_type_id, origin, checker_id, interactor_id,
        created_at)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, NOW())
        RETURNING id
    ''', (task_id, short_code, full_name, time_lim_ms, mem_lim_kb,
          testing_type_id, origin, checker_id, interactor_id))
    return cursor.fetchone()[0]


def create_md_statement(cursor,
                        story, input, output, notes, scoring,
                        task_version_id, lang_iso639_1):
    '''Create a new markdown_statement and returns its ID'''
    cursor.execute('''
        INSERT INTO markdown_statements
        (story, input, output, notes, scoring, task_version_id,
        lang_iso639_1)
        VALUES (%s, %s, %s, %s, %s, %s, %s)
        RETURNING id
    ''', (story, input, output, notes, scoring, task_version_id,
          lang_iso639_1))
    return cursor.fetchone()[0]


def create_version_author(cursor, task_version_id, author):
    '''Create a new task version author entry'''
    cursor.execute('''
        INSERT INTO version_authors
        (task_version_id, author)
        VALUES (%s, %s)
    ''', (task_version_id, author))


def textfile_exists(cursor, sha256):
    '''Returns True if a text file with the given sha256 exists'''
    cursor.execute('''
        SELECT EXISTS(
            SELECT 1
            FROM text_files
            WHERE sha256 = %s
        )
    ''', (sha256,))
    return cursor.fetchone()[0]


def get_textfile_id(cursor, sha256):
    '''Returns the ID of a text file with the given sha256'''
    cursor.execute('''
        SELECT id FROM text_files
        WHERE sha256 = %s
    ''', (sha256,))
    return cursor.fetchone()[0]


def create_textfile(cursor, sha256, content):
    '''Creates a new text file and returns its ID'''
    cursor.execute('''
        INSERT INTO text_files
        (sha256, content)
        VALUES (%s, %s)
        RETURNING id
    ''', (sha256, content))
    return cursor.fetchone()[0]


def flyway_checksum_sum(cursor):
    '''Returns the checksum sum of the flyway_schema_history table'''
    cursor.execute('''
        SELECT SUM(checksum) FROM flyway_schema_history
    ''')
    return cursor.fetchone()[0]


def create_checker(cursor, code):
    """Create a new checker and returns its ID"""
    cursor.execute('''
        INSERT INTO testlib_checkers
        (code)
        VALUES (%s)
        RETURNING id
    ''', (code,))
    return cursor.fetchone()[0]


def ensure_checker(cursor, code):
    """Create a new checker if it doesn't exist and returns its ID"""
    cursor.execute('''
        SELECT id FROM testlib_checkers
        WHERE code = %s
    ''', (code,))
    res = cursor.fetchone()
    if res is None:
        return create_checker(cursor, code)
    else:
        return res[0]


def assign_checker(cursor, task_version_id, checker_id):
    """Assign a checker to a task version"""
    cursor.execute('''
        UPDATE task_versions
        SET checker_id = %s
        WHERE id = %s
    ''', (checker_id, task_version_id))


def ensure_textfile(cursor, content):
    """Create a new textfile if it doesn't exist and returns its ID"""
    sha256 = sha256_bytes(content)
    decoded = content.decode('utf-8')
    cursor.execute('''
        SELECT id FROM text_files
        WHERE sha256 = %s
    ''', (sha256,))
    res = cursor.fetchone()
    if res is None:
        return create_textfile(cursor, sha256, decoded)
    else:
        return res[0]


def create_task_version_test(cursor, test_filename, task_version_id,
                             input_text_file_id, answer_text_file_id):
    cursor.execute('''
        INSERT INTO task_version_tests
        (test_filename, task_version_id,
        input_text_file_id, answer_text_file_id)
        VALUES (%s, %s, %s, %s)
        RETURNING id
    ''', (test_filename, task_version_id,
          input_text_file_id, answer_text_file_id))
    return cursor.fetchone()[0]
