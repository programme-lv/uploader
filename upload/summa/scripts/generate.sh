#!/usr/bin/env bash

ROOT_DIR=$(readlink -f "$(dirname "$0")/..")
TEMP_DIR="$ROOT_DIR/temporary"
GEN_DIR="$ROOT_DIR/generation"
SOL_DIR="$ROOT_DIR/solutions"

if [ ! -f "$GEN_DIR/generator.cpp" ]; then
    echo "generator.cpp not found"
    exit 1
fi

if [ ! -f "$GEN_DIR/generate.jl" ]; then
    echo "generate.jl not found"
    exit 1
fi

if [ ! -f "$SOL_DIR/correct.cpp" ]; then
    echo "correct.cpp not found"
    exit 1
fi

if [ ! -d "$TEMP_DIR" ]; then
    echo "creating temporary directory"
    mkdir -p $TEMP_DIR
fi

echo "compiling generator"
g++ -o $TEMP_DIR/generator $GEN_DIR/generator.cpp

echo "compiling solution"
g++ -o $TEMP_DIR/solution $SOL_DIR/correct.cpp

GEN_SCRIPT=$(readlink -f "$GEN_DIR/generate.jl")
GEN_EXE=$(readlink -f "$TEMP_DIR/generator")
SOL_EXE=$(readlink -f "$TEMP_DIR/solution")
TEST_DIR="$ROOT_DIR/tests"

export GEN_EXE
export SOL_EXE

if [ -d "$TEST_DIR" ]; then
    echo "removing old test directory"
    rm -r $TEST_DIR
fi

echo "creating new test directory"
mkdir -p $TEST_DIR

echo "generating tests"
cd $TEST_DIR
julia $GEN_SCRIPT

echo "finished"

