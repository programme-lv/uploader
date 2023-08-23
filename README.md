# Programme.lv task uploader

The python script contained in this repository can
be used by administrators of [programme.lv](https://programme.lv)
to upload tasks to the database.

## Usage

1) Create a task following the examples in [example-tasks](https://github.com/programme-lv/example-tasks) repository.
2) Place the task in the `upload` directory.
3) Run `python upload.py`.
4) Input your username when queried.

## Prerequisites

1) Install the required python packages with `pip install -r requirements.txt`.
2) Create `.env` file with database credentials. See `.env.example` for reference.

## Workflow

1) Create task if it doesn't exist in the database.
2) Create a new task version of the task.
	- tests
	- checker
	- statement
3) Assign the new task version as the relevant one.
