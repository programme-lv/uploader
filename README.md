# Programme.lv task uploader

The python script contained in this repository can
be used by administrators of [programme.lv](https://programme.lv)
to upload tasks to the database.

## Usage

1) Create a task following the examples in [example-tasks](https://github.com/programme-lv/example-tasks) repository.
2) Place the task in the `upload` directory as a directory or a `.zip`.
3) Create `.env` file. See `.env.example` for reference.
4) Run `python upload.py`.

Each task belongs to a specific user. The user is determined by the
username provided in the `.env` file. The user must exist in the database.

## Prerequisites

1) Install the required python packages with `pip install -r requirements.txt`.

The `requirements.txt` isn't maintained thoroughly. If you encounter any
missing packages, please add them to the file.

## Workflow

The script will do the following:
1) assert task specification version;
2) validate provided `problem.toml` & file structure;
3) create task if it doesn't exist in the database;
4) create a new task version of the task;
5) create tests for the task version;
	- tests
	- checker
	- statement
5) assign the new task version as the relevant one;
