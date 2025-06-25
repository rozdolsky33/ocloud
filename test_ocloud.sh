#!/bin/bash

# Test script for ocloud CLI
# This script tests various combinations of commands and flags for the ocloud CLI

# Source environment variables from .env file
# Function to print section headers
print_header() {
    echo "============================================================"
    echo "  $1"
    echo "============================================================"
    echo ""
}

# Function to run a command and print the command before executing
run_command() {
    echo "$ $@"
    "$@"
    echo ""
}

# Test version command and flag
print_header "Testing version command and flag"
run_command ./bin/ocloud version
run_command ./bin/ocloud -v
run_command ./bin/ocloud --version

# Test root command with help
print_header "Testing root command with help"
run_command ./bin/ocloud --help

# Test root command with global flags
print_header "Testing root command with global flags"
run_command ./bin/ocloud --compartment $OCI_COMPARTMENT
run_command ./bin/ocloud -c $OCI_COMPARTMENT

# Test compute command
print_header "Testing compute command"
run_command ./bin/ocloud compute --help
run_command ./bin/ocloud comp --help

# Test compute instance command
print_header "Testing compute instance command"
run_command ./bin/ocloud compute instance --help
run_command ./bin/ocloud comp inst --help

# Test compute instance list command
print_header "Testing compute instance list command"
run_command ./bin/ocloud compute instance list
run_command ./bin/ocloud compute instance list
run_command ./bin/ocloud compute instance list --limit 10 --page 1 --json
run_command ./bin/ocloud compute instance list -m 10 -p 1 -j
run_command ./bin/ocloud comp inst l

# Test compute instance find command
print_header "Testing compute instance find command"
run_command ./bin/ocloud compute instance find "test"
run_command ./bin/ocloud compute instance find "test"
run_command ./bin/ocloud compute instance find "test" --image-details --json
run_command ./bin/ocloud compute instance find "test" -i -j
run_command ./bin/ocloud comp inst f "test"

# Test with debug flag
print_header "Testing with debug flag"
run_command ./bin/ocloud -d compute instance list
run_command ./bin/ocloud --debug compute instance list

# Test with color flag
print_header "Testing with color flag"
run_command ./bin/ocloud --color compute instance list

# Test with disable concurrency flag
print_header "Testing with disable concurrency flag"
run_command ./bin/ocloud -x compute instance list
run_command ./bin/ocloud --disable-concurrency compute instance list

# Test identity command
print_header "Testing identity command"
run_command ./bin/ocloud identity --help
run_command ./bin/ocloud ident --help
run_command ./bin/ocloud idt --help

# Test identity compartment command
print_header "Testing identity compartment command"
run_command ./bin/ocloud identity compartment --help
run_command ./bin/ocloud identity compart --help
run_command ./bin/ocloud ident compart --help

# Test identity compartment list command
print_header "Testing identity compartment list command"
run_command ./bin/ocloud identity compartment list
run_command ./bin/ocloud identity compartment list
run_command ./bin/ocloud identity compartment list --limit 10 --page 1 --json
run_command ./bin/ocloud identity compartment list -m 10 -p 1 -j
run_command ./bin/ocloud ident compart l

# Test identity compartment find command
print_header "Testing identity compartment find command"
run_command ./bin/ocloud identity compartment find "test"
run_command ./bin/ocloud identity compartment find "test"
run_command ./bin/ocloud identity compartment find "test" --json
run_command ./bin/ocloud identity compartment find "test" -j
run_command ./bin/ocloud ident compart f "test"

# Test database command
print_header "Testing database command"
run_command ./bin/ocloud database --help
run_command ./bin/ocloud db --help

# Test database autonomousdb command
print_header "Testing database autonomousdb command"
run_command ./bin/ocloud database autonomousdb --help
run_command ./bin/ocloud database adb --help
run_command ./bin/ocloud db adb --help

# Test database autonomousdb list command
print_header "Testing database autonomousdb list command"
run_command ./bin/ocloud database autonomousdb list
run_command ./bin/ocloud database autonomousdb list
run_command ./bin/ocloud database autonomousdb list --limit 10 --page 1 --json
run_command ./bin/ocloud database autonomousdb list -m 10 -p 1 -j
run_command ./bin/ocloud db adb l

# Test database autonomousdb find command
print_header "Testing database autonomousdb find command"
run_command ./bin/ocloud database autonomousdb find "test"
run_command ./bin/ocloud database autonomousdb find "test"
run_command ./bin/ocloud database autonomousdb find "test" --json
run_command ./bin/ocloud database autonomousdb find "test" -j
run_command ./bin/ocloud db adb f "test"

print_header "All tests completed"
