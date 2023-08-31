
def valid_task_fs(task_dir, specification) -> (bool, str):
    '''Returns  (True, None) if task_dir is a valid task,
                (False, error) otherwise'''


def is_lowercase_alphanum(s: str) -> bool:
    '''Returns True if s is a lowercase alphanumerical string,
    False otherwise'''
    return all(c.islower() or c.isdigit() for c in s)


def validate_toml(problem_toml: dict):
    '''Returns True if toml is a valid problem.toml,
    False otherwise'''
    assert is_lowercase_alphanum(problem_toml['code']), \
        f"Invalid code: {problem_toml['code']}"
    assert len(problem_toml['code']) <= 32, \
        f"Code too long: {problem_toml['code']}"

    assert len(problem_toml['name']) <= 128, \
        f"Name too long: {problem_toml['name']}"

    assert 0 < problem_toml['time'] <= 10, \
        f"Invalid time: {problem_toml['time']}"

    assert 0 < problem_toml['memory'] <= 2048, \
        f"Invalid memory: {problem_toml['memory']}"

    assert problem_toml['type'] == 'simple', \
        f"Invalid type: {problem_toml['type']}"

    toml_fields = set(problem_toml.keys())
    for field in ['specification', 'code', 'name', 'time', 'memory',
                  'type', 'authors', 'tags']:
        toml_fields.discard(field)

    assert len(toml_fields) == 0, \
        f"Invalid fields: {toml_fields}"
