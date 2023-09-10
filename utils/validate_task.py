from os.path import isdir, isfile, join
from os import listdir
import re


def is_lowercase_alphanum(s: str) -> bool:
    '''Returns True if s is a lowercase alphanumerical string,
    False otherwise'''
    return all(c.islower() or c.isdigit() for c in s)


def is_numerical_integer(value):
    return isinstance(value, int) or \
        (isinstance(value, str) and value.isdigit())


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

    assert 1 <= problem_toml['difficulty'] <= 5 and \
        is_numerical_integer(problem_toml['difficulty']), \
        f"Invalid difficulty: {problem_toml['difficulty']}"

    toml_fields = set(problem_toml.keys())
    for field in ['specification', 'code', 'name', 'time', 'memory',
                  'type', 'authors', 'tags', 'difficulty']:
        toml_fields.discard(field)

    assert len(toml_fields) == 0, \
        f"Invalid fields: {toml_fields}"


def validate_task_fs(task_dirpath, problem_toml):
    assert problem_toml['type'] == 'simple', \
        f"Invalid type: {problem_toml['type']}"

    # checker.cpp in evaluation dir
    assert isdir(join(task_dirpath, 'evaluation')), \
        "evaluation directory not found"
    assert isfile(join(task_dirpath, 'evaluation', 'checker.cpp')), \
        "checker.cpp not found in evaluation directory"

    # digit filenames for examples in examples dir
    if isdir(join(task_dirpath, 'examples')):
        example_pattern = re.compile(r'^\d+\.(in|ans)$')

        for filename in listdir(join(task_dirpath, 'examples')):
            assert not isdir(join(task_dirpath, 'examples', filename)), \
                f"Example {filename} is a fucking directory"
            assert example_pattern.match(filename), \
                f"Invalid example filename: {filename}"

    # check for at least one valid statement
    # for now just check if there is a statement.md
    # TODO: also check other languages, pdfs
    assert isfile(join(task_dirpath, 'statements',
                       'markdown', 'lv', 'story.md')), \
        "statement.md not found in statements/markdown/lv directory"

    # ensure tests exist
    assert isdir(join(task_dirpath, 'tests')), \
        "tests directory not found"
    assert len(listdir(join(task_dirpath, 'tests'))) > 0, \
        "no tests in tests directory"
    for test in listdir(join(task_dirpath, 'tests')):
        assert not isdir(join(task_dirpath, 'tests', test)), \
            f"Test {test} is a fucking directory"
        test_pattern = re.compile(r'^.+\.(in|ans)$')
        assert test_pattern.match(test), \
            f"Invalid test filename: {test}"

    # check for unnecessary directories or files
    dir_files = set(listdir(task_dirpath))
    for file in ['problem.toml', 'evaluation', 'examples', 'generation',
                 'scripts', 'solutions', 'statements', 'temporary', 'tests']:
        dir_files.discard(file)
