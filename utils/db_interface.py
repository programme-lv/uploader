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
    '''Returns a user with the given username'''
    cursor.execute('''
        SELECT * FROM users
        WHERE username = %s
    ''', (username,))
    return cursor.fetchone()


def create_task(cursor, created_by_id,
                relevant_version_id=None, published_version_id=None):
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
                   checker_text_id=None, interactor_text_id=None):
    '''Creates a new task version and returns its ID'''
    cursor.execute('''
        INSERT INTO task_versions
        (task_id, short_code, full_name, time_lim_ms, mem_lim_kb,
        testing_type_id, origin, checker_text_id, interactor_text_id)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
        RETURNING id
    ''', (task_id, short_code, full_name, time_lim_ms, mem_lim_kb,
          testing_type_id, origin, checker_text_id, interactor_text_id))
    return cursor.fetchone()[0]


def create_md_statement(cursor,
                        story, input, output, notes, scoring,
                        task_version_id):
    '''Create a new markdown_statement and returns its ID'''
    cursor.execute('''
        INSERT INTO markdown_statements
        (story, input, output, notes, scoring, task_version_id)
        VALUES (%s, %s, %s, %s, %s, %s)
        RETURNING id
    ''', (story, input, output, notes, scoring, task_version_id))
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


def create_version_test(cursor, task_version_id,
                        input_text_file_id, answer_text_file_id,
                        test_filename=""):
    '''Link task version with test and returns its ID'''
    cursor.execute('''
        INSERT INTO task_version_tests
        (task_version_id, input_text_file_id, answer_text_file_id,
        test_filename)
        VALUES (%s, %s, %s, %s)
        RETURNING id
    ''', (task_version_id, input_text_file_id, answer_text_file_id,
          test_filename))
    return cursor.fetchone()[0]
