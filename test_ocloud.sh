#!/bin/bash

# Test script for ocloud CLI
# This script tests various combinations of commands and flags for the ocloud CLI

# Define color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Array to store errors
errors=()

# Function to print section headers
print_header() {
    echo "============================================================"
    echo "  $1"
    echo "============================================================"
    echo ""
}

# Function to run a command and print the command before executing
# Captures exit code and displays errors in red
run_command() {
    echo "$" "$@"
    "$@"
    exit_code=$?

    if [ $exit_code -ne 0 ]; then
        # Print error message in red
        echo -e "${RED}Command failed with exit code $exit_code${NC}"

        # Store the error for summary - concatenate all args into a single string
        cmd_str="Command failed: $(printf "%s " "$@")"
        errors+=("${cmd_str%?}")  # Remove trailing space
    fi

    echo ""
}
# Test version command and flag
print_header "Testing ocloud settings"
run_command ./bin/ocloud

# Test version command and flag
print_header "Testing version command and flag"
run_command ./bin/ocloud version
run_command ./bin/ocloud -v
run_command ./bin/ocloud --version

# Test root command with help
print_header "Testing root command with help"
run_command ./bin/ocloud --help

# Test config info map-file command
print_header "Testing config info map-file command"
run_command ./bin/ocloud config info map-file
run_command ./bin/ocloud config info map-file --json
run_command ./bin/ocloud config info map-file --realm OC1
run_command ./bin/ocloud config info map-file --realm OC1 --json

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
run_command ./bin/ocloud compute instance get
run_command ./bin/ocloud compute instance get --limit 10 --page 1 --json
run_command ./bin/ocloud compute instance get -m 10 -p 1 -j
run_command ./bin/ocloud comp inst get

# Test compute instance search command
print_header "Testing compute instance search command"
run_command ./bin/ocloud compute instance search "roster"
run_command ./bin/ocloud compute instance search "roster" --all --json
run_command ./bin/ocloud compute instance search "roster" -A -j
run_command ./bin/ocloud comp inst s "roster"

# Test compute image command
print_header "Testing compute image command"
run_command ./bin/ocloud compute image --help
run_command ./bin/ocloud comp img --help

# Test compute image get command
print_header "Testing compute image get command"
run_command ./bin/ocloud compute image get
run_command ./bin/ocloud compute image get --limit 10 --page 1 --json
run_command ./bin/ocloud compute image get -m 10 -p 1 -j
run_command ./bin/ocloud comp img get

# Test compute image find command
print_header "Testing compute image find command"
run_command ./bin/ocloud compute image search "Oracle-Linux"
run_command ./bin/ocloud compute image search "Oracle-Linux" --json
run_command ./bin/ocloud compute image search "Oracle-Linux" -j
run_command ./bin/ocloud comp img s "Oracle-Linux"

# Test compute oke command
print_header "Testing compute oke command"
run_command ./bin/ocloud compute oke --help
run_command ./bin/ocloud comp oke --help

# Test compute oke get command
print_header "Testing compute oke list command"
run_command ./bin/ocloud compute oke get
run_command ./bin/ocloud compute oke get --limit 10 --page 1 --json
run_command ./bin/ocloud compute oke get -m 10 -p 1 -j
run_command ./bin/ocloud comp oke get

# Test compute oke find command
print_header "Testing compute oke find command"
run_command ./bin/ocloud compute oke search "orion"
run_command ./bin/ocloud compute oke search "orion" --json
run_command ./bin/ocloud compute oke search "orion" -j
run_command ./bin/ocloud comp oke s "orion"

# Test with debug flag
print_header "Testing with debug flag"
run_command ./bin/ocloud -d compute instance get
run_command ./bin/ocloud --debug compute instance get

# Test with color flag
print_header "Testing with color flag"
run_command ./bin/ocloud --color compute instance get

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
print_header "Testing identity compartment get command"
run_command ./bin/ocloud identity compartment get
run_command ./bin/ocloud identity compartment get --limit 10 --page 1 --json
run_command ./bin/ocloud identity compartment get -m 10 -p 1 -j


# Test identity compartment find command
print_header "Testing identity compartment find command"
run_command ./bin/ocloud identity compartment find "sand"
run_command ./bin/ocloud identity compartment find "sand" --json
run_command ./bin/ocloud identity compartment find "sand" -j
run_command ./bin/ocloud ident compart f "sand"

# Test identity policy command
print_header "Testing identity policy command"
run_command ./bin/ocloud identity policy --help
run_command ./bin/ocloud identity pol --help
run_command ./bin/ocloud ident pol --help

# Test identity policy get command
print_header "Testing identity policy get command"
run_command ./bin/ocloud identity policy get
run_command ./bin/ocloud identity policy get
run_command ./bin/ocloud identity policy get --limit 10 --page 1 --json
run_command ./bin/ocloud identity policy get -m 10 -p 1 -j

# Test identity policy find command
print_header "Testing identity policy find command"
run_command ./bin/ocloud identity policy find "monitor"
run_command ./bin/ocloud identity policy find "monitor" --json
run_command ./bin/ocloud identity policy find "monitor" -j
run_command ./bin/ocloud ident pol f "monitor"

# Test network command
print_header "Testing network command"
run_command ./bin/ocloud network --help
run_command ./bin/ocloud net --help

# Test network subnet command
print_header "Testing network subnet command"
run_command ./bin/ocloud network subnet --help
run_command ./bin/ocloud network sub --help
run_command ./bin/ocloud net sub --help

# Test network subnet list command
print_header "Testing network subnet list command"
run_command ./bin/ocloud network subnet list
run_command ./bin/ocloud network subnet list --limit 10 --page 1 --json
run_command ./bin/ocloud network subnet list -m 10 -p 1 -j
run_command ./bin/ocloud net sub l

# Test network subnet find command
print_header "Testing network subnet find command"
run_command ./bin/ocloud network subnet find "pub"
run_command ./bin/ocloud network subnet find "pub" --json
run_command ./bin/ocloud network subnet find "pub" -j
run_command ./bin/ocloud net sub f "pub"

# Test network vcn command
print_header "Testing network vcn command"
run_command ./bin/ocloud network vcn --help
run_command ./bin/ocloud net vcn --help

# Test network vcn get command (no interactive list)
print_header "Testing network vcn get command"
run_command ./bin/ocloud network vcn get
run_command ./bin/ocloud network vcn get --limit 10 --page 1 --json
run_command ./bin/ocloud network vcn get -m 10 -p 1 -j
# with network-related flags
run_command ./bin/ocloud network vcn get --gateway --subnet --nsg --route-table --security-list
run_command ./bin/ocloud network vcn get --all
# with short aliases for flags
run_command ./bin/ocloud network vcn get -G -S -N -R -L -j
run_command ./bin/ocloud network vcn get -A -j

# Test network vcn find command
print_header "Testing network vcn find command"
run_command ./bin/ocloud network vcn find "prod"
run_command ./bin/ocloud network vcn find "prod" --json
run_command ./bin/ocloud network vcn find "prod" --all
run_command ./bin/ocloud network vcn find "prod" -A -j

# Test network loadbalancer get command
print_header "Testing network loadbalancer get command"
run_command ./bin/ocloud network loadbalancer get
run_command ./bin/ocloud network loadbalancer get --limit 10 --page 1 --json
run_command ./bin/ocloud network loadbalancer get -m 10 -p 1 -j
run_command ./bin/ocloud network loadbalancer get --all
run_command ./bin/ocloud net lb get
run_command ./bin/ocloud net lb get -A -j

# Test network loadbalancer find command
print_header "Testing network loadbalancer find command"
run_command ./bin/ocloud network loadbalancer find "prod"
run_command ./bin/ocloud network loadbalancer find "prod" --json
run_command ./bin/ocloud network loadbalancer find "prod" --all
run_command ./bin/ocloud net lb f "prod"
run_command ./bin/ocloud net lb f "prod" -A -j

# Test storage object-storage get command
print_header "Testing storage object-storage get command"
run_command ./bin/ocloud storage object-storage get
run_command ./bin/ocloud storage object-storage get --limit 10 --page 1 --json
run_command ./bin/ocloud storage object-storage get -m 10 -p 1 -j
run_command ./bin/ocloud storage object-storage get
run_command ./bin/ocloud storage os get
run_command ./bin/ocloud storage os get -j

# Test database command
print_header "Testing database command"
run_command ./bin/ocloud database --help
run_command ./bin/ocloud db --help

# Test database autonomousdb command
print_header "Testing database autonomousdb command"
run_command ./bin/ocloud database autonomous --help
run_command ./bin/ocloud database adb --help
run_command ./bin/ocloud db adb --help

# Test database autonomousdb list command
print_header "Testing database autonomousdb list command"
run_command ./bin/ocloud database autonomous get
run_command ./bin/ocloud database autonomous get --limit 10 --page 1 --json
run_command ./bin/ocloud database autonomous get -m 10 -p 1 -j

# Test database autonomousdb find command
print_header "Testing database autonomousdb find command"
run_command ./bin/ocloud database autonomous find "test"
run_command ./bin/ocloud database autonomous find "test" --json
run_command ./bin/ocloud database autonomous find "test" -j
run_command ./bin/ocloud db adb f "test"

# Test version command and flag
print_header "Testing ocloud settings"
run_command ./bin/ocloud

print_header "All tests completed"

# Display error summary if there were any errors
if [ ${#errors[@]} -gt 0 ]; then
    echo -e "${RED}ERROR SUMMARY:${NC}"
    echo -e "${RED}=============${NC}"
    for error in "${errors[@]}"; do
        echo -e "${RED}$error${NC}"
    done
    echo ""
    echo -e "${RED}Total errors: ${#errors[@]}${NC}"
    exit 1
else
    echo -e "${GREEN}All commands completed successfully!${NC}"
fi
